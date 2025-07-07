package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/thekubefleet/kubefleet/internal/grpcclient"
	"github.com/thekubefleet/kubefleet/internal/k8s"
	"github.com/thekubefleet/kubefleet/internal/metrics"
	agentpb "github.com/thekubefleet/kubefleet/proto"
)

func main() {
	fmt.Println("KubeFleet Agent starting...")

	// Get server address from environment or use default
	serverAddr := os.Getenv("KUBEFLEET_SERVER_ADDR")
	if serverAddr == "" {
		serverAddr = "localhost:50051" // Default gRPC server address
	}

	// Initialize Kubernetes client
	k8sClient, err := k8s.NewClient()
	if err != nil {
		log.Fatalf("Failed to create Kubernetes client: %v", err)
	}

	// Initialize metrics collector
	metricsCollector, err := metrics.NewCollector()
	if err != nil {
		log.Fatalf("Failed to create metrics collector: %v", err)
	}

	// Initialize gRPC client
	grpcClient, err := grpcclient.NewClient(serverAddr)
	if err != nil {
		log.Fatalf("Failed to create gRPC client: %v", err)
	}
	defer grpcClient.Close()

	// Main loop
	ticker := time.NewTicker(30 * time.Second) // Report every 30 seconds
	defer ticker.Stop()

	ctx := context.Background()

	for {
		select {
		case <-ticker.C:
			if err := collectAndReport(ctx, k8sClient, metricsCollector, grpcClient); err != nil {
				log.Printf("Error collecting and reporting data: %v", err)
			}
		}
	}
}

func collectAndReport(ctx context.Context, k8sClient *k8s.Client, metricsCollector *metrics.Collector, grpcClient *grpcclient.Client) error {
	// Get all namespaces
	namespaces, err := k8sClient.GetNamespaces(ctx)
	if err != nil {
		return fmt.Errorf("failed to get namespaces: %w", err)
	}

	// Collect resource information for each namespace
	var resourceInfos []*agentpb.ResourceInfo
	for _, namespace := range namespaces {
		pods, err := k8sClient.GetPodsInNamespace(ctx, namespace)
		if err != nil {
			log.Printf("Failed to get pods in namespace %s: %v", namespace, err)
			continue
		}

		deployments, err := k8sClient.GetDeploymentsInNamespace(ctx, namespace)
		if err != nil {
			log.Printf("Failed to get deployments in namespace %s: %v", namespace, err)
			continue
		}

		resourceInfo := grpcclient.ConvertResourceInfo(namespace, pods, deployments)
		resourceInfos = append(resourceInfos, resourceInfo)
	}

	// Collect metrics
	metricsData, err := metricsCollector.CollectAllMetrics(ctx, namespaces)
	if err != nil {
		return fmt.Errorf("failed to collect metrics: %w", err)
	}

	// Convert metrics to protobuf format
	protoMetrics := grpcclient.ConvertResourceMetrics(metricsData)

	// Create agent data
	agentData := &agentpb.AgentData{
		Resources: resourceInfos,
		Metrics:   protoMetrics,
		Timestamp: time.Now().Unix(),
	}

	// Send data via gRPC
	if err := grpcClient.SendAgentData(ctx, agentData); err != nil {
		return fmt.Errorf("failed to send agent data: %w", err)
	}

	fmt.Printf("Successfully reported data for %d namespaces with %d metrics\n", len(namespaces), len(protoMetrics))
	return nil
}
