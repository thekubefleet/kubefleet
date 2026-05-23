# KubeFleet Helm Chart

This chart deploys the full KubeFleet application:

- `kubefleet-dashboard` Deployment + Service (+ optional Ingress)
- `kubefleet-agent` Deployment
- Agent ServiceAccount + ClusterRole + ClusterRoleBinding

## Prerequisite: metrics-server

KubeFleet reads CPU and memory usage from Kubernetes Metrics API. Install metrics-server first:

```bash
helm repo add metrics-server https://kubernetes-sigs.github.io/metrics-server/
helm repo update metrics-server
helm upgrade -i metrics-server metrics-server/metrics-server -n kube-system --set args={--kubelet-insecure-tls}
```

## Install

```bash
helm upgrade -i kubefleet ./charts/kubefleet -n kubefleet --create-namespace
```

## Useful values

- `agent.image.repository`, `agent.image.tag`
- `dashboard.image.repository`, `dashboard.image.tag`
- `dashboard.ingress.enabled`
- `dashboard.ingress.host`
- `namespaceOverride`

## Example custom install

```bash
helm upgrade -i kubefleet ./charts/kubefleet \
  -n kubefleet --create-namespace \
  --set dashboard.ingress.enabled=true \
  --set dashboard.ingress.host=kubefleet.example.com \
  --set agent.image.repository=ghcr.io/thekubefleet/kubefleet-agent \
  --set dashboard.image.repository=ghcr.io/thekubefleet/kubefleet-dashboard
```
