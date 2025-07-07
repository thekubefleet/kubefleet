# KubeFleet 🚀

[![Go Version](https://img.shields.io/badge/Go-1.24+-blue.svg)](https://golang.org)
[![Node Version](https://img.shields.io/badge/Node-18+-green.svg)](https://nodejs.org)
[![License](https://img.shields.io/badge/License-MIT-yellow.svg)](LICENSE)
[![PRs Welcome](https://img.shields.io/badge/PRs-welcome-brightgreen.svg)](CONTRIBUTING.md)
[![Contributor Covenant](https://img.shields.io/badge/Contributor%20Covenant-2.0-4baaaa.svg)](CODE_OF_CONDUCT.md)

A modern Kubernetes monitoring solution with an intelligent agent that collects cluster data and a beautiful React dashboard for real-time visualization.

## ✨ Features

### 🕵️ Agent
- **Smart Discovery**: Automatically discovers all namespaces and resources
- **Resource Monitoring**: Tracks pods, deployments, and services in real-time
- **Performance Metrics**: Collects CPU and memory usage data
- **gRPC Communication**: Fast, efficient data transmission to dashboard
- **Kubernetes Native**: Runs as a pod with proper RBAC permissions

### 📊 Dashboard
- **Real-time Updates**: Live data refresh every 30 seconds
- **Cluster Overview**: High-level statistics and resource counts
- **Namespace Explorer**: Interactive drill-down into namespaces and resources
- **Performance Charts**: Beautiful, interactive charts for metrics visualization
- **Modern UI**: Dark theme with Material-UI components
- **Responsive Design**: Works perfectly on desktop and mobile

## 🏗️ Architecture

```
┌─────────────────┐    gRPC     ┌─────────────────┐    HTTP     ┌─────────────────┐
│   Kubernetes    │ ──────────► │   KubeFleet     │ ──────────► │   React         │
│   Agent         │             │   Dashboard     │             │   Frontend      │
│                 │             │   Server        │             │                 │
└─────────────────┘             └─────────────────┘             └─────────────────┘
```

## 🚀 Quick Start

### Prerequisites

- **Go** (>=1.18)
- **Node.js** (>=16)
- **protoc** (Protocol Buffers compiler)
- **Docker** (for container builds)
- **Kubernetes cluster** (for deployment)

### Local Development

1. **Clone and setup:**
   ```bash
   git clone https://github.com/your-username/kubefleet.git
   cd kubefleet
   go mod tidy
   ```

2. **Generate protobuf code:**
   ```bash
   protoc --go_out=. --go_opt=paths=source_relative --go-grpc_out=. --go-grpc_opt=paths=source_relative proto/agent.proto
   ```

3. **Start the dashboard server:**
   ```bash
   go run ./cmd/server
   ```

4. **Start the React development server:**
   ```bash
   cd dashboard
   npm start
   ```

5. **Run the agent (in another terminal):**
   ```bash
   go run ./cmd/agent
   ```

6. **Access the dashboard:** http://localhost:3000

### Docker Deployment

```bash
# Build containers
docker build -t kubefleet-agent:latest .
docker build -f Dockerfile.dashboard -t kubefleet-dashboard:latest .

# Deploy to Kubernetes
kubectl apply -f deploy/dashboard-deployment.yaml
kubectl apply -f deploy/agent-deployment.yaml

# Access the dashboard
kubectl port-forward svc/kubefleet-dashboard 3000:3000
```

## 📁 Project Structure

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
├── .github/            # GitHub templates and workflows
└── docs/               # Documentation
```

## ⚙️ Configuration

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

## 🔌 API Reference

### Dashboard Server Endpoints

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

## 🤝 Contributing

We welcome contributions! Please see our [Contributing Guide](CONTRIBUTING.md) for details.

### Development Setup

1. Fork the repository
2. Create a feature branch
3. Make your changes
4. Add tests
5. Submit a pull request

### Code Style

- **Go**: Follow [Effective Go](https://golang.org/doc/effective_go.html)
- **React**: Use TypeScript, functional components, and Material-UI
- **Commits**: Use conventional commit format

## 📋 Roadmap

- [ ] **Enhanced Metrics**: Kubernetes Metrics API integration
- [ ] **Prometheus Support**: Direct Prometheus metrics collection
- [ ] **Alerting**: Built-in alerting capabilities
- [ ] **Authentication**: User authentication and authorization
- [ ] **Multi-cluster**: Support for multiple Kubernetes clusters
- [ ] **Helm Chart**: Official Helm chart for easy deployment
- [ ] **Performance**: Resource usage trends and predictions
- [ ] **Custom Dashboards**: User-configurable dashboard layouts

## 🐛 Troubleshooting

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

## 📄 License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.

## 🆘 Support

- **Issues**: [GitHub Issues](https://github.com/your-username/kubefleet/issues)
- **Discussions**: [GitHub Discussions](https://github.com/your-username/kubefleet/discussions)
- **Security**: [Security Policy](SECURITY.md)

## 🙏 Acknowledgments

- [Kubernetes client-go](https://github.com/kubernetes/client-go) for Kubernetes API integration
- [Material-UI](https://mui.com/) for the beautiful React components
- [Recharts](https://recharts.org/) for data visualization
- [gRPC](https://grpc.io/) for efficient communication

---

**Made with ❤️ by the KubeFleet community**
