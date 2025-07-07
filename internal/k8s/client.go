package k8s

import (
	"context"
	"fmt"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
)

type Client struct {
	clientset *kubernetes.Clientset
}

// NewClient creates a new Kubernetes client
func NewClient() (*Client, error) {
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

	return &Client{clientset: clientset}, nil
}

// GetNamespaces returns all namespaces in the cluster
func (c *Client) GetNamespaces(ctx context.Context) ([]string, error) {
	namespaces, err := c.clientset.CoreV1().Namespaces().List(ctx, metav1.ListOptions{})
	if err != nil {
		return nil, fmt.Errorf("failed to list namespaces: %w", err)
	}

	var names []string
	for _, ns := range namespaces.Items {
		names = append(names, ns.Name)
	}
	return names, nil
}

// GetPodsInNamespace returns all pods in a specific namespace
func (c *Client) GetPodsInNamespace(ctx context.Context, namespace string) ([]string, error) {
	pods, err := c.clientset.CoreV1().Pods(namespace).List(ctx, metav1.ListOptions{})
	if err != nil {
		return nil, fmt.Errorf("failed to list pods in namespace %s: %w", namespace, err)
	}

	var podNames []string
	for _, pod := range pods.Items {
		podNames = append(podNames, pod.Name)
	}
	return podNames, nil
}

// GetDeploymentsInNamespace returns all deployments in a specific namespace
func (c *Client) GetDeploymentsInNamespace(ctx context.Context, namespace string) ([]string, error) {
	deployments, err := c.clientset.AppsV1().Deployments(namespace).List(ctx, metav1.ListOptions{})
	if err != nil {
		return nil, fmt.Errorf("failed to list deployments in namespace %s: %w", namespace, err)
	}

	var deploymentNames []string
	for _, deployment := range deployments.Items {
		deploymentNames = append(deploymentNames, deployment.Name)
	}
	return deploymentNames, nil
}

// GetServicesInNamespace returns all services in a specific namespace
func (c *Client) GetServicesInNamespace(ctx context.Context, namespace string) ([]string, error) {
	services, err := c.clientset.CoreV1().Services(namespace).List(ctx, metav1.ListOptions{})
	if err != nil {
		return nil, fmt.Errorf("failed to list services in namespace %s: %w", namespace, err)
	}

	var serviceNames []string
	for _, service := range services.Items {
		serviceNames = append(serviceNames, service.Name)
	}
	return serviceNames, nil
}
