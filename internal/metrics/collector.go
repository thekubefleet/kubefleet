package metrics

import (
	"context"
	"fmt"
	"time"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
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
	clientset *kubernetes.Clientset
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

	return &Collector{clientset: clientset}, nil
}

// CollectPodMetrics collects metrics for all pods in a namespace
func (c *Collector) CollectPodMetrics(ctx context.Context, namespace string) ([]ResourceMetric, error) {
	pods, err := c.clientset.CoreV1().Pods(namespace).List(ctx, metav1.ListOptions{})
	if err != nil {
		return nil, fmt.Errorf("failed to list pods in namespace %s: %w", namespace, err)
	}

	var metrics []ResourceMetric
	for _, pod := range pods.Items {
		// For now, we'll use placeholder metrics
		// In a real implementation, you'd query the metrics API or Prometheus
		metric := ResourceMetric{
			Namespace: namespace,
			Name:      pod.Name,
			Kind:      "Pod",
			CPU:       0.0, // TODO: Get actual CPU usage
			Memory:    0.0, // TODO: Get actual memory usage
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

	var metrics []ResourceMetric
	for _, deployment := range deployments.Items {
		// For now, we'll use placeholder metrics
		// In a real implementation, you'd aggregate metrics from the deployment's pods
		metric := ResourceMetric{
			Namespace: namespace,
			Name:      deployment.Name,
			Kind:      "Deployment",
			CPU:       0.0, // TODO: Get actual CPU usage
			Memory:    0.0, // TODO: Get actual memory usage
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
