# kubefleet

A Kubernetes monitoring solution with an agent that collects cluster data and a React dashboard for visualization.

## Features

### Agent
- **Namespace Discovery**: Automatically discovers all namespaces in the cluster
- **Resource Monitoring**: Lists pods, deployments, and services in each namespace
- **Metrics Collection**: Gathers performance metrics (CPU, memory) from resources
- **gRPC Communication**: Sends data to the dashboard server via gRPC

### Dashboard
- **Real-time Monitoring**: Live updates every 30 seconds
- **Cluster Overview**: High-level statistics and resource counts
- **Namespace Explorer**: Detailed view of all namespaces and their resources
- **Performance Charts**: Interactive charts for CPU and memory metrics
- **Modern UI**: Built with React, TypeScript, and Material-UI

## Project Structure

```
kubefleet/
├── cmd/
│   ├── agent/          # Agent entrypoint
│   └── server/         # Dashboard server entrypoint
├── internal/
│   ├── k8s/            # Kubernetes API logic
│   ├── metrics/        # Metrics collection
│   ├── grpcclient/     # gRPC client logic
│   └── server/         # Dashboard server logic
├── dashboard/          # React frontend
│   ├── src/
│   │   ├── components/ # React components
│   │   └── App.tsx     # Main app
│   └── package.json
├── proto/              # Protobuf definitions
├── deploy/             # Kubernetes manifests
└── Dockerfile*         # Container builds
```

## Quick Start

### Prerequisites

1. **Go** (>=1.18)
2. **Node.js** (>=16)
3. **protoc** (for protobuf/gRPC)
4. **Docker** (for building containers)
5. **Kubernetes cluster** (for deployment)

### Local Development

1. **Clone and setup:**
   ```sh
   git clone <your-repo>
   cd kubefleet
   go mod tidy
   ```

2. **Generate protobuf code:**
   ```sh
   protoc --go_out=. --go_opt=paths=source_relative --go-grpc_out=. --go-grpc_opt=paths=source_relative proto/agent.proto
   ```

3. **Start the dashboard server:**
   ```sh
   go run ./cmd/server
   ```

4. **Start the React development server:**
   ```sh
   cd dashboard
   npm start
   ```

5. **Run the agent (in another terminal):**
   ```sh
   go run ./cmd/agent
   ```

6. **Access the dashboard:** http://localhost:3000

### Building Containers

```sh
# Build agent
docker build -t kubefleet-agent:latest .

# Build dashboard
docker build -f Dockerfile.dashboard -t kubefleet-dashboard:latest .
```

### Deploying to Kubernetes

1. **Deploy the dashboard:**
   ```sh
   kubectl apply -f deploy/dashboard-deployment.yaml
   ```

2. **Deploy the agent:**
   ```sh
   kubectl apply -f deploy/agent-deployment.yaml
   ```

3. **Check the deployment:**
   ```sh
   kubectl get pods -l app=kubefleet-dashboard
   kubectl get pods -l app=kubefleet-agent
   ```

4. **Access the dashboard:**
   ```sh
   kubectl port-forward svc/kubefleet-dashboard 3000:3000
   ```
   Then visit http://localhost:3000

## Configuration

### Environment Variables

**Agent:**
- `KUBEFLEET_SERVER_ADDR`: gRPC server address (default: localhost:50051)

**Dashboard Server:**
- `HTTP_PORT`: HTTP server port (default: 3000)
- `GRPC_PORT`: gRPC server port (default: 50051)

### RBAC Permissions

The agent requires the following permissions:
- Read access to namespaces, pods, services, and deployments
- Read access to metrics API (if available)

## API Endpoints

### Dashboard Server

- `GET /api/data` - Get all historical data
- `GET /api/data/latest` - Get the latest data point
- `GET /api/health` - Health check endpoint

### gRPC Service

The agent sends data using the `AgentReporter` service:

```protobuf
service AgentReporter {
  rpc ReportData(AgentData) returns (ReportResponse);
}
```

## Development

### Adding New Metrics

1. Update the protobuf definition in `proto/agent.proto`
2. Regenerate the Go code
3. Update the metrics collector in `internal/metrics/`
4. Update the dashboard components to display the new data

### Adding New Dashboard Components

1. Create a new component in `dashboard/src/components/`
2. Add it to the main App component
3. Implement data fetching from the API endpoints

## Troubleshooting

### Common Issues

1. **Agent can't connect to dashboard:**
   - Check the `KUBEFLEET_SERVER_ADDR` environment variable
   - Verify the dashboard service is running

2. **Dashboard shows no data:**
   - Check if the agent is running and sending data
   - Verify the gRPC connection is working
   - Check the dashboard server logs

3. **React app not loading:**
   - Ensure the React development server is running on port 3001
   - Check browser console for errors

## Next Steps

- [ ] Implement actual metrics collection from Kubernetes Metrics API
- [ ] Add Prometheus integration for detailed metrics
- [ ] Add configuration for collection intervals
- [ ] Implement retry logic for gRPC failures
- [ ] Add authentication and authorization
- [ ] Create Helm chart for easier deployment
- [ ] Add alerting capabilities
- [ ] Implement resource usage trends and predictions
