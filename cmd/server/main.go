package main

import (
	"context"
	"log"
	"net"
	"net/http"
	"os"
	"sync"

	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"

	"github.com/thekubefleet/kubefleet/internal/server"
	agentpb "github.com/thekubefleet/kubefleet/proto"
)

type grpcServer struct {
	agentpb.UnimplementedAgentReporterServer
	dataStore *server.DataStore
	mu        sync.RWMutex
}

func (s *grpcServer) ReportData(ctx context.Context, data *agentpb.AgentData) (*agentpb.ReportResponse, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	// Store the received data
	s.dataStore.StoreAgentData(data)

	log.Printf("Received data from agent: %d resources, %d metrics", len(data.Resources), len(data.Metrics))

	return &agentpb.ReportResponse{
		Success: true,
		Message: "Data received successfully",
	}, nil
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

	// Initialize data store
	dataStore := server.NewDataStore()

	// Create gRPC server
	grpcSrv := grpc.NewServer()
	agentpb.RegisterAgentReporterServer(grpcSrv, &grpcServer{
		dataStore: dataStore,
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
