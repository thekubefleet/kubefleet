# Contributing to KubeFleet

Thank you for your interest in contributing to KubeFleet! This document provides guidelines and information for contributors.

## Table of Contents

- [Code of Conduct](#code-of-conduct)
- [Getting Started](#getting-started)
- [Development Setup](#development-setup)
- [Making Changes](#making-changes)
- [Testing](#testing)
- [Submitting Changes](#submitting-changes)
- [Issue Guidelines](#issue-guidelines)
- [Pull Request Guidelines](#pull-request-guidelines)

## Code of Conduct

This project and everyone participating in it is governed by our Code of Conduct. By participating, you are expected to uphold this code.

## Getting Started

1. **Fork the repository** on GitHub
2. **Clone your fork** locally
3. **Set up the development environment** (see below)
4. **Create a feature branch** for your changes

## Development Setup

### Prerequisites

- Go 1.24 or later
- Node.js 24 or later
- Docker (optional, for container builds)
- Kubernetes cluster (for testing)
- protoc (Protocol Buffers compiler)

### Local Development

1. **Clone and setup:**

   ```bash
   git clone https://github.com/thekubefleet/kubefleet.git
   cd kubefleet
   go mod tidy
   ```

2. **Generate protobuf code:**

   ```bash
   protoc --go_out=. --go_opt=paths=source_relative --go-grpc_out=. --go-grpc_opt=paths=source_relative proto/agent.proto
   ```

3. **Setup React dashboard:**

   ```bash
   cd dashboard
   npm install
   cd ..
   ```

4. **Start development servers:**

   ```bash
   # Terminal 1: Dashboard server
   go run ./cmd/server
   
   # Terminal 2: React dev server
   cd dashboard && npm start
   
   # Terminal 3: Agent (optional, for testing)
   go run ./cmd/agent
   ```

## Making Changes

### Project Structure

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
├── proto/              # Protobuf definitions
├── deploy/             # Kubernetes manifests
└── docs/               # Documentation
```

### Code Style Guidelines

#### Go Code

- Follow [Effective Go](https://golang.org/doc/effective_go.html)
- Use `gofmt` for formatting
- Add comments for exported functions and types
- Keep functions small and focused
- Use meaningful variable names

#### React/TypeScript Code

- Follow TypeScript best practices
- Use functional components with hooks
- Follow Material-UI design patterns
- Add proper TypeScript types
- Use meaningful component and variable names

#### General Guidelines

- Write clear commit messages
- Add tests for new functionality
- Update documentation for API changes
- Follow existing code patterns

## Testing

### Go Tests

```bash
# Run all tests
go test ./...

# Run tests with coverage
go test -cover ./...

# Run specific package tests
go test ./internal/k8s
```

### React Tests

```bash
cd dashboard
npm test
```

### Integration Tests

```bash
# Build and test containers
docker build -t kubefleet-agent:test .
docker build -f Dockerfile.dashboard -t kubefleet-dashboard:test .

# Deploy to test cluster
kubectl apply -f deploy/dashboard-deployment.yaml
kubectl apply -f deploy/agent-deployment.yaml
```

## Submitting Changes

### Before Submitting

1. **Ensure tests pass** locally
2. **Update documentation** if needed
3. **Check code formatting** (`gofmt`, `npm run format`)
4. **Verify the build** works (`go build`, `npm run build`)

### Commit Message Format

Use conventional commit format:

```
type(scope): description

[optional body]

[optional footer]
```

Examples:

- `feat(agent): add namespace filtering`
- `fix(dashboard): resolve chart rendering issue`
- `docs(readme): update installation instructions`

### Types

- `feat`: New feature
- `fix`: Bug fix
- `docs`: Documentation changes
- `style`: Code style changes
- `refactor`: Code refactoring
- `test`: Adding tests
- `chore`: Maintenance tasks

## Issue Guidelines

### Before Creating an Issue

1. **Search existing issues** to avoid duplicates
2. **Check the documentation** for solutions
3. **Try the latest version** of the project

### Issue Template

When creating an issue, please include:

- **Description**: Clear description of the problem
- **Steps to reproduce**: Detailed steps to reproduce the issue
- **Expected behavior**: What you expected to happen
- **Actual behavior**: What actually happened
- **Environment**: OS, Go version, Node version, Kubernetes version
- **Screenshots**: If applicable

## Pull Request Guidelines

### Before Creating a PR

1. **Ensure your branch is up to date** with main
2. **Write clear commit messages**
3. **Add tests** for new functionality
4. **Update documentation** if needed
5. **Test locally** before submitting

### PR Template

When creating a PR, please include:

- **Description**: What does this PR do?
- **Type of change**: Bug fix, feature, documentation, etc.
- **Testing**: How was this tested?
- **Breaking changes**: Any breaking changes?
- **Screenshots**: If UI changes

### PR Review Process

1. **Automated checks** must pass
2. **Code review** by maintainers
3. **Address feedback** and make requested changes
4. **Maintainer approval** required for merge

## Getting Help

- **GitHub Issues**: For bugs and feature requests
- **GitHub Discussions**: For questions and general discussion
- **Documentation**: Check the README and docs folder

## Recognition

Contributors will be recognized in:

- GitHub contributors list
- Release notes
- Project documentation

Thank you for contributing to KubeFleet! 🚀

