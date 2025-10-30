package main

import (
	"context"
	"log"
	"net"
	"net/http"
	"os"
	"strings"
	"sync"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"

	"github.com/thekubefleet/kubefleet/internal/k8s"
	"github.com/thekubefleet/kubefleet/internal/server"
	agentpb "github.com/thekubefleet/kubefleet/proto"
)

type grpcServer struct {
	agentpb.UnimplementedAgentReporterServer
	dataStore *server.DataStore
	k8sClient *k8s.Client
	mu        sync.RWMutex
}

func (s *grpcServer) ReportData(ctx context.Context, data *agentpb.AgentData) (*agentpb.ReportResponse, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	// Store the received data
	s.dataStore.StoreAgentData(data)

	log.Printf("Received data from agent: %d resources, %d metrics, %d logs", len(data.Resources), len(data.Metrics), len(data.Logs))

	return &agentpb.ReportResponse{
		Success: true,
		Message: "Data received successfully",
	}, nil
}

func (s *grpcServer) StreamPodLogs(req *agentpb.LogRequest, stream agentpb.AgentReporter_StreamPodLogsServer) error {
	ctx := stream.Context()

	// Get containers for the pod
	containers, err := s.k8sClient.GetPodContainers(ctx, req.Namespace, req.PodName)
	if err != nil {
		return err
	}

	// If no specific container requested, get logs from all containers
	if req.ContainerName == "" {
		for _, containerName := range containers {
			if err := s.streamContainerLogs(ctx, req, containerName, stream); err != nil {
				log.Printf("Error streaming logs for container %s: %v", containerName, err)
			}
		}
	} else {
		// Stream logs for specific container
		if err := s.streamContainerLogs(ctx, req, req.ContainerName, stream); err != nil {
			return err
		}
	}

	// Send completion signal
	return stream.Send(&agentpb.LogStream{
		Logs:       []*agentpb.PodLog{},
		IsComplete: true,
	})
}

func (s *grpcServer) streamContainerLogs(ctx context.Context, req *agentpb.LogRequest, containerName string, stream agentpb.AgentReporter_StreamPodLogsServer) error {
	// Get initial logs
	logLines, err := s.k8sClient.GetPodLogs(ctx, req.Namespace, req.PodName, containerName, int64(req.TailLines), false)
	if err != nil {
		return err
	}

	// Convert and send initial logs
	podLogs := convertLogsToProto(req.Namespace, req.PodName, containerName, logLines)
	if len(podLogs) > 0 {
		if err := stream.Send(&agentpb.LogStream{
			Logs:       podLogs,
			IsComplete: false,
		}); err != nil {
			return err
		}
	}

	// If follow is requested, continue streaming
	if req.Follow {
		lastLogTime := time.Now()
		ticker := time.NewTicker(5 * time.Second) // Check for new logs every 5 seconds
		defer ticker.Stop()

		for {
			select {
			case <-ctx.Done():
				return ctx.Err()
			case <-ticker.C:
				// Get logs since last check
				newLogLines, err := s.k8sClient.GetPodLogsSince(ctx, req.Namespace, req.PodName, containerName, lastLogTime)
				if err != nil {
					log.Printf("Error getting new logs: %v", err)
					continue
				}

				if len(newLogLines) > 0 {
					newPodLogs := convertLogsToProto(req.Namespace, req.PodName, containerName, newLogLines)
					if err := stream.Send(&agentpb.LogStream{
						Logs:       newPodLogs,
						IsComplete: false,
					}); err != nil {
						return err
					}
					lastLogTime = time.Now()
				}
			}
		}
	}

	return nil
}

func convertLogsToProto(namespace, podName, containerName string, logLines []string) []*agentpb.PodLog {
	var podLogs []*agentpb.PodLog
	now := time.Now().Unix()

	for _, line := range logLines {
		// Split log lines (they might contain multiple lines)
		lines := strings.Split(strings.TrimSpace(line), "\n")
		for _, logLine := range lines {
			if logLine == "" {
				continue
			}

			protoLog := &agentpb.PodLog{
				Namespace:     namespace,
				PodName:       podName,
				ContainerName: containerName,
				LogLine:       logLine,
				Timestamp:     now,
				Level:         parseLogLevel(logLine),
			}
			podLogs = append(podLogs, protoLog)
		}
	}

	return podLogs
}

func parseLogLevel(logLine string) string {
	line := strings.ToUpper(strings.TrimSpace(logLine))

	if strings.Contains(line, "ERROR") || strings.Contains(line, "FATAL") {
		return "ERROR"
	}
	if strings.Contains(line, "WARN") {
		return "WARN"
	}
	if strings.Contains(line, "DEBUG") {
		return "DEBUG"
	}
	if strings.Contains(line, "INFO") {
		return "INFO"
	}

	return "INFO" // Default level
}

func main() {
	// Get port from environment or use default
	grpcPort := os.Getenv("GRPC_PORT")
	if grpcPort == "" {
		grpcPort = "50051"
	}

	httpPort := os.Getenv("HTTP_PORT")
	if httpPort == "" {
		httpPort = "3000"
	}

	// Initialize Kubernetes client
	k8sClient, err := k8s.NewClient()
	if err != nil {
		log.Fatalf("Failed to create Kubernetes client: %v", err)
	}

	// Initialize data store
	dataStore := server.NewDataStore()

	// Create gRPC server
	grpcSrv := grpc.NewServer()
	agentpb.RegisterAgentReporterServer(grpcSrv, &grpcServer{
		dataStore: dataStore,
		k8sClient: k8sClient,
	})

	// Enable reflection for debugging
	reflection.Register(grpcSrv)

	// Start gRPC server
	go func() {
		lis, err := net.Listen("tcp", ":"+grpcPort)
		if err != nil {
			log.Fatalf("Failed to listen for gRPC: %v", err)
		}
		log.Printf("gRPC server listening on port %s", grpcPort)
		if err := grpcSrv.Serve(lis); err != nil {
			log.Fatalf("Failed to serve gRPC: %v", err)
		}
	}()

	// Create HTTP server for the dashboard
	httpServer := server.NewHTTPServer(dataStore)

	// Start HTTP server
	log.Printf("HTTP server listening on port %s", httpPort)
	if err := http.ListenAndServe(":"+httpPort, httpServer); err != nil {
		log.Fatalf("Failed to serve HTTP: %v", err)
	}
}
