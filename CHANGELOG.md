# Changelog

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.0.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [Unreleased]

### Added
- **Real-time Log Monitoring**: Live pod log streaming with automatic updates
- **Multi-container Log Support**: View logs from all containers in a pod
- **Log Level Detection**: Automatic parsing of ERROR, WARN, INFO, DEBUG levels
- **Log Filtering**: Filter logs by level and search functionality
- **Log Download**: Export logs as text files
- **Enhanced UI**: Log viewer modal with Material-UI components
- **gRPC Log Streaming**: New StreamPodLogs endpoint for real-time log delivery
- **HTTP Log Endpoints**: REST API endpoints for log retrieval
- **Proxy Configuration**: React development server proxy for API calls

### Changed
- Updated protobuf definitions to include log-related messages
- Enhanced Kubernetes client with log collection methods
- Improved dashboard UI with log viewing capabilities
- Updated TypeScript configuration for better compatibility

## [1.1.0] - 2024-01-XX

### Added
- **Log Collection**: Agent now collects pod logs every 30 seconds
- **Log Viewer Component**: Full-featured React component for log display
- **Log API Endpoints**: HTTP endpoints for retrieving logs by namespace/pod/container
- **Log Streaming**: gRPC streaming for real-time log updates
- **Log Level Parsing**: Automatic detection of log levels from log content
- **Container Selection**: Choose specific containers or view all
- **Search and Filter**: Find specific text and filter by log level
- **Download Functionality**: Export filtered logs as text files
- **Auto-refresh**: Toggle real-time log updates
- **Color-coded Logs**: Visual distinction for different log levels

### Technical Details
- **Backend**: Enhanced k8s client with GetPodLogs, GetPodContainers methods
- **Frontend**: New PodLogs.tsx component with filtering and search
- **API**: New /api/logs endpoints for log retrieval
- **gRPC**: StreamPodLogs method for real-time streaming
- **UI**: Modal dialog for log viewing with responsive design

## [1.0.0] - 2024-01-XX

### Added
- **Agent Features:**
  - Kubernetes namespace discovery
  - Pod, deployment, and service monitoring
  - Performance metrics collection (CPU, memory)
  - gRPC client for data transmission
  - Automatic retry logic
  - RBAC integration

- **Dashboard Features:**
  - Real-time cluster overview
  - Interactive namespace explorer
  - Performance metrics charts
  - Dark theme Material-UI design
  - Responsive layout
  - Auto-refresh every 30 seconds

- **Infrastructure:**
  - Multi-stage Docker builds
  - Kubernetes deployment manifests
  - Service and ingress configuration
  - Health checks and monitoring

- **Development:**
  - TypeScript React frontend
  - Go backend with gRPC server
  - Protocol Buffers for data serialization
  - Comprehensive test suite
  - CI/CD pipeline with GitHub Actions

### Technical Details
- **Agent:** Go 1.24, client-go, gRPC
- **Dashboard:** React 18, TypeScript, Material-UI, Recharts
- **Communication:** gRPC with Protocol Buffers
- **Deployment:** Docker containers, Kubernetes manifests
- **Monitoring:** Built-in health checks and metrics

---

## Version History

- **1.1.0** - Log monitoring and real-time log streaming
- **1.0.0** - Initial release with core monitoring functionality
- **Unreleased** - Development version with latest features

## Contributing

Please read [CONTRIBUTING.md](CONTRIBUTING.md) for details on our code of conduct and the process for submitting pull requests. 