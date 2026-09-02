# Kubernetes deployment (kind, locally)

One Helm chart per service; environment differences live in values files
only. Every Deployment ships the full checklist from
docs/standards/infra.md: distinct readiness (`/readyz`) and liveness
(`/healthz`) probes, envFrom ConfigMap (+ optional Secret), resource requests
and memory limits, runAsNonRoot + readOnlyRootFilesystem + dropped
capabilities, a per-service ServiceAccount (the IRSA seam for AWS), and
`terminationGracePeriodSeconds` aligned with the app's shutdown timeout.

Charts:

- `platform/` — local backing services: single-node Redpanda, the whoami demo
  upstream, and the `topic-init` Job (twelve-factor XII admin process).
- `edge-proxy/`, `enforcement/` — the datacat services.

The Flink classifier is not deployed here yet; it runs on the Compose Flink
cluster until the Flink-operator phase. In AWS, `platform/` is replaced by
MSK + real services via Terraform.

## Run it

```bash
# 1. Cluster
kind create cluster --name datacat

# 2. Images (repo root context) loaded into kind - no registry needed
docker build -f services/edge-proxy/Dockerfile  -t datacat/edge-proxy:dev .
docker build -f services/enforcement/Dockerfile -t datacat/enforcement:dev .
kind load docker-image datacat/edge-proxy:dev datacat/enforcement:dev --name datacat

# 3. Install
helm install platform    deploy/k8s/platform    -n datacat --create-namespace
helm install edge-proxy  deploy/k8s/edge-proxy  -n datacat
helm install enforcement deploy/k8s/enforcement -n datacat
kubectl -n datacat rollout status deploy/edge-proxy deploy/enforcement

# 4. Reach it
kubectl -n datacat port-forward svc/edge-proxy 8080:8080 &
kubectl -n datacat port-forward svc/enforcement 8081:8081 &
curl -b dc_session=me localhost:8080/          # proxied via whoami
curl localhost:8081/v1/decisions/me            # decision lookup

# 5. Tear down
kind delete cluster --name datacat
```

Rebuild-and-redeploy loop after a code change:

```bash
docker build -f services/edge-proxy/Dockerfile -t datacat/edge-proxy:dev .
kind load docker-image datacat/edge-proxy:dev --name datacat
kubectl -n datacat rollout restart deploy/edge-proxy
```
