package metrics

import (
	"context"
	"fmt"
	"time"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"

	versioned "k8s.io/metrics/pkg/client/clientset/versioned"
)

type ResourceMetric struct {
	Namespace string
	Name      string
	Kind      string
	CPU       float64
	Memory    float64
	Timestamp time.Time
}

type Collector struct {
	clientset     *kubernetes.Clientset
	metricsClient *versioned.Clientset
}

// NewCollector creates a new metrics collector
func NewCollector() (*Collector, error) {
	var config *rest.Config
	var err error

	// Try in-cluster config first (when running inside Kubernetes)
	config, err = rest.InClusterConfig()
	if err != nil {
		// Fall back to kubeconfig file (for local development)
		config, err = clientcmd.BuildConfigFromFlags("", clientcmd.RecommendedHomeFile)
		if err != nil {
			return nil, fmt.Errorf("failed to get kubeconfig: %w", err)
		}
	}

	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		return nil, fmt.Errorf("failed to create kubernetes client: %w", err)
	}

	metricsClient, err := versioned.NewForConfig(config)
	if err != nil {
		return nil, fmt.Errorf("failed to create metrics client: %w", err)
	}

	return &Collector{clientset: clientset, metricsClient: metricsClient}, nil
}

// CollectPodMetrics collects metrics for all pods in a namespace
func (c *Collector) CollectPodMetrics(ctx context.Context, namespace string) ([]ResourceMetric, error) {
	pods, err := c.clientset.CoreV1().Pods(namespace).List(ctx, metav1.ListOptions{})
	if err != nil {
		return nil, fmt.Errorf("failed to list pods in namespace %s: %w", namespace, err)
	}

	podMetricsList, err := c.metricsClient.MetricsV1beta1().PodMetricses(namespace).List(ctx, metav1.ListOptions{})
	if err != nil {
		return nil, fmt.Errorf("failed to get pod metrics in namespace %s: %w", namespace, err)
	}

	// Map pod name to metrics
	metricsMap := make(map[string]struct{ cpu, mem float64 })
	for _, podMetric := range podMetricsList.Items {
		var totalCPU, totalMem int64
		for _, c := range podMetric.Containers {
			cpu := c.Usage.Cpu().MilliValue() // millicores
			mem := c.Usage.Memory().Value()   // bytes
			totalCPU += cpu
			totalMem += mem
		}
		metricsMap[podMetric.Name] = struct{ cpu, mem float64 }{
			cpu: float64(totalCPU) / 1000.0,            // convert to cores
			mem: float64(totalMem) / (1024.0 * 1024.0), // convert to MiB
		}
	}

	var metrics []ResourceMetric
	for _, pod := range pods.Items {
		m := metricsMap[pod.Name]
		metric := ResourceMetric{
			Namespace: namespace,
			Name:      pod.Name,
			Kind:      "Pod",
			CPU:       m.cpu,
			Memory:    m.mem,
			Timestamp: time.Now(),
		}
		metrics = append(metrics, metric)
	}

	return metrics, nil
}

// CollectDeploymentMetrics collects metrics for all deployments in a namespace
func (c *Collector) CollectDeploymentMetrics(ctx context.Context, namespace string) ([]ResourceMetric, error) {
	deployments, err := c.clientset.AppsV1().Deployments(namespace).List(ctx, metav1.ListOptions{})
	if err != nil {
		return nil, fmt.Errorf("failed to list deployments in namespace %s: %w", namespace, err)
	}

	podMetricsList, err := c.metricsClient.MetricsV1beta1().PodMetricses(namespace).List(ctx, metav1.ListOptions{})
	if err != nil {
		return nil, fmt.Errorf("failed to get pod metrics in namespace %s: %w", namespace, err)
	}

	// Map pod name to metrics
	metricsMap := make(map[string]struct{ cpu, mem float64 })
	for _, podMetric := range podMetricsList.Items {
		var totalCPU, totalMem int64
		for _, c := range podMetric.Containers {
			cpu := c.Usage.Cpu().MilliValue()
			mem := c.Usage.Memory().Value()
			totalCPU += cpu
			totalMem += mem
		}
		metricsMap[podMetric.Name] = struct{ cpu, mem float64 }{
			cpu: float64(totalCPU) / 1000.0,
			mem: float64(totalMem) / (1024.0 * 1024.0),
		}
	}

	var metrics []ResourceMetric
	for _, deployment := range deployments.Items {
		selector := deployment.Spec.Selector.MatchLabels
		labelSelector := metav1.FormatLabelSelector(&metav1.LabelSelector{MatchLabels: selector})
		pods, err := c.clientset.CoreV1().Pods(namespace).List(ctx, metav1.ListOptions{LabelSelector: labelSelector})
		if err != nil {
			return nil, fmt.Errorf("failed to list pods for deployment %s: %w", deployment.Name, err)
		}
		var totalCPU, totalMem float64
		for _, pod := range pods.Items {
			m := metricsMap[pod.Name]
			totalCPU += m.cpu
			totalMem += m.mem
		}
		metric := ResourceMetric{
			Namespace: namespace,
			Name:      deployment.Name,
			Kind:      "Deployment",
			CPU:       totalCPU,
			Memory:    totalMem,
			Timestamp: time.Now(),
		}
		metrics = append(metrics, metric)
	}

	return metrics, nil
}

// CollectAllMetrics collects metrics for all resources in all namespaces
func (c *Collector) CollectAllMetrics(ctx context.Context, namespaces []string) ([]ResourceMetric, error) {
	var allMetrics []ResourceMetric

	for _, namespace := range namespaces {
		// Collect pod metrics
		podMetrics, err := c.CollectPodMetrics(ctx, namespace)
		if err != nil {
			return nil, fmt.Errorf("failed to collect pod metrics for namespace %s: %w", namespace, err)
		}
		allMetrics = append(allMetrics, podMetrics...)

		// Collect deployment metrics
		deploymentMetrics, err := c.CollectDeploymentMetrics(ctx, namespace)
		if err != nil {
			return nil, fmt.Errorf("failed to collect deployment metrics for namespace %s: %w", namespace, err)
		}
		allMetrics = append(allMetrics, deploymentMetrics...)
	}

	return allMetrics, nil
}
