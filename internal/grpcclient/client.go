package grpcclient

import (
	"context"
	"fmt"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	"github.com/thekubefleet/kubefleet/internal/metrics"
	agentpb "github.com/thekubefleet/kubefleet/proto"
)

type Client struct {
	conn   *grpc.ClientConn
	client agentpb.AgentReporterClient
}

// NewClient creates a new gRPC client
func NewClient(serverAddr string) (*Client, error) {
	conn, err := grpc.Dial(serverAddr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return nil, fmt.Errorf("failed to connect to gRPC server: %w", err)
	}

	client := agentpb.NewAgentReporterClient(conn)

	return &Client{
		conn:   conn,
		client: client,
	}, nil
}

// Close closes the gRPC connection
func (c *Client) Close() error {
	return c.conn.Close()
}

// SendAgentData sends agent data to the UI server
func (c *Client) SendAgentData(ctx context.Context, data *agentpb.AgentData) error {
	ctx, cancel := context.WithTimeout(ctx, 30*time.Second)
	defer cancel()

	response, err := c.client.ReportData(ctx, data)
	if err != nil {
		return fmt.Errorf("failed to send agent data: %w", err)
	}

	if !response.Success {
		return fmt.Errorf("server returned error: %s", response.Message)
	}

	return nil
}

// ConvertResourceInfo converts internal resource info to protobuf format
func ConvertResourceInfo(namespace string, pods, deployments []string) *agentpb.ResourceInfo {
	return &agentpb.ResourceInfo{
		Namespace:   namespace,
		Pods:        pods,
		Deployments: deployments,
	}
}

// ConvertResourceMetrics converts internal metrics to protobuf format
func ConvertResourceMetrics(metricsData []metrics.ResourceMetric) []*agentpb.ResourceMetrics {
	var protoMetrics []*agentpb.ResourceMetrics

	for _, metric := range metricsData {
		protoMetric := &agentpb.ResourceMetrics{
			Namespace: metric.Namespace,
			Name:      metric.Name,
			Kind:      metric.Kind,
			Cpu:       metric.CPU,
			Memory:    metric.Memory,
		}
		protoMetrics = append(protoMetrics, protoMetric)
	}

	return protoMetrics
}
