# Infrastructure Standards (Terraform, Docker, Kubernetes, CI/CD)

Infrastructure is code and follows the same rules as code: small focused units,
one task per module, everything reviewed, nothing hand-edited in a console.
If it isn't in Terraform or a manifest, it doesn't exist.

## Terraform

### Layout — modules are the "functions" of infra

```
infra/terraform/
├── modules/                  reusable, environment-agnostic
│   ├── network/              VPC, subnets, endpoints
│   ├── eks/                  cluster, node groups, IRSA
│   ├── kafka/                MSK (or Redpanda on EKS via helm_release)
│   ├── dynamo/               tables, TTL, capacity mode
│   ├── storage/              S3 buckets, lifecycle rules
│   ├── verdict-api/          Lambda + API Gateway
│   └── agent-verification/   Step Functions + IAM
└── envs/
    ├── local/                LocalStack targets (dev/prod parity, factor X)
    ├── staging/
    └── prod/                 each env: backend.tf, main.tf, terraform.tfvars
```

Rules:

- Environments differ ONLY in `terraform.tfvars` and backend config. An env's
  `main.tf` is nothing but module calls — the orchestrator pattern again.
- Every module declares typed, validated, described `variable` blocks and
  explicit `output` blocks. No module reaches into another's internals;
  wiring happens via outputs → inputs in the env layer.
- Remote state in S3 with state locking; state is never committed, never
  edited by hand. `terraform fmt` + `terraform validate` + tflint in CI.
- `plan` runs on every PR and is posted as a comment; `apply` only from CI on
  merge to main, never from a laptop.

```hcl
# modules/dynamo/variables.tf
variable "table_name" {
	type        = string
	description = "DynamoDB table name, prefixed by environment"

	validation {
		condition     = can(regex("^[a-z][a-z0-9-]+$", var.table_name))
		error_message = "table_name must be lowercase kebab-case."
	}
}

variable "ttl_attribute" {
	type        = string
	description = "Attribute holding epoch-seconds expiry; empty disables TTL"
	default     = ""
}
```

```hcl
# envs/staging/main.tf — module calls only, like a composition root
module "network" {
	source      = "../../modules/network"
	environment = var.environment
	cidr_block  = var.vpc_cidr
}

module "verdicts_table" {
	source        = "../../modules/dynamo"
	table_name    = "${var.environment}-datacat-verdicts"
	hash_key      = "session_id"
	ttl_attribute = "expires_at"
}
```

### IAM — least privilege, per service

Every service gets its own IAM role via IRSA (IAM Roles for Service Accounts).
Policies name exact resources — no `"Resource": "*"`, no shared "app role".
The enforcement service can write its decisions table and nothing else.

## Docker

Multi-stage, distroless, non-root, pinned bases. One image per service; the
same image runs locally and in EKS (factors V and X).

```dockerfile
# services/enforcement/Dockerfile
FROM golang:1.26 AS build
WORKDIR /src
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 go build -trimpath -ldflags="-s -w" -o /out/enforcement ./cmd/server

FROM gcr.io/distroless/static-debian12:nonroot
COPY --from=build /out/enforcement /enforcement
USER nonroot
ENTRYPOINT ["/enforcement"]
```

Rules: `.dockerignore` always; no secrets in build args or layers; images
tagged with the git SHA (never a mutable `latest` in deployments); Trivy scan
gates the push.

## Kubernetes

Every workload ships the full checklist — a Deployment missing any of these is
incomplete:

```yaml
# deploy/k8s/enforcement/deployment.yaml (Helm-templated; literals shown)
apiVersion: apps/v1
kind: Deployment
metadata:
  name: enforcement
spec:
  replicas: 2
  selector:
    matchLabels: { app: enforcement }
  template:
    metadata:
      labels: { app: enforcement }
    spec:
      serviceAccountName: enforcement        # IRSA binding
      securityContext:
        runAsNonRoot: true
        seccompProfile: { type: RuntimeDefault }
      containers:
        - name: enforcement
          image: <ecr>/datacat/enforcement:<git-sha>
          ports:
            - containerPort: 8080
          envFrom:
            - configMapRef: { name: enforcement-config }   # factor III
            - secretRef: { name: enforcement-secrets }
          resources:
            requests: { cpu: 100m, memory: 128Mi }
            limits: { memory: 256Mi }
          readinessProbe:
            httpGet: { path: /readyz, port: 8080 }
            periodSeconds: 5
          livenessProbe:
            httpGet: { path: /healthz, port: 8080 }
            periodSeconds: 10
          securityContext:
            allowPrivilegeEscalation: false
            readOnlyRootFilesystem: true
            capabilities: { drop: ["ALL"] }
      terminationGracePeriodSeconds: 30       # matches app shutdown timeout
```

Rules:

- `/healthz` = process is alive; `/readyz` = dependencies reachable. Never the
  same check: a Kafka outage must flip readiness (stop traffic), not liveness
  (restart loops fix nothing).
- ConfigMaps for non-secret env, Secrets for secrets; both created by
  Terraform/Helm, never `kubectl create` by hand.
- Resource requests AND memory limits always. CPU limits omitted deliberately
  (throttling hurts latency-sensitive services); revisit per service.
- Probes + PodDisruptionBudget + `terminationGracePeriodSeconds` aligned with
  the app's graceful-shutdown timeout (factor IX).
- One Helm chart per service under `deploy/k8s/`; env differences live in
  `values-<env>.yaml` only.

## CI/CD (GitHub Actions)

One reusable workflow per language, called by thin per-service workflows.
Pipeline stages, in order, each gating the next:

```
lint → unit tests → build image → Trivy scan → push to ECR (sha tag)
     → deploy to staging (helm upgrade) → smoke test → [manual gate] → prod
```

```yaml
# .github/workflows/service-go.yaml (reusable, abbreviated)
on:
  workflow_call:
    inputs:
      service: { required: true, type: string }

jobs:
  verify:
    runs-on: ubuntu-latest
    defaults: { run: { working-directory: services/${{ inputs.service }} } }
    steps:
      - uses: actions/checkout@v4
      - uses: actions/setup-go@v5
        with: { go-version-file: services/${{ inputs.service }}/go.mod }
      - run: golangci-lint run
      - run: go test -race -cover ./...
  package:
    needs: verify
    steps:
      - run: docker build -t $ECR/datacat/${{ inputs.service }}:${{ github.sha }} .
      - run: trivy image --exit-code 1 --severity HIGH,CRITICAL $IMAGE
      - run: docker push $IMAGE   # OIDC auth to AWS — no long-lived keys
```

DevSecOps gates, all blocking:

- **Secrets**: gitleaks scan on every PR.
- **Dependencies**: `govulncheck` (Go), OWASP dependency-check (Java),
  `npm audit` (React) — fail on known-exploitable.
- **Images**: Trivy, fail on HIGH/CRITICAL.
- **Auth to AWS**: GitHub OIDC federation only. Long-lived access keys are
  forbidden in CI and on laptops alike.
- **Terraform**: `plan` on PR, `apply` on merge, tflint + checkov.

## Observability (CloudWatch)

- All services log JSON to stdout → Fluent Bit DaemonSet → CloudWatch Logs
  (factor XI). Log groups per service, retention set by Terraform.
- Metrics: each service exposes Prometheus-style `/metrics`; key business
  metrics (verdicts/sec by class, enforcement actions, classifier lag) also
  published to CloudWatch for alarming.
- Alarms as code (Terraform): consumer lag, verdict latency p99, error rates,
  Flink checkpoint failures. An alarm without a runbook link is incomplete.

## Cost discipline

- Local first: kind + Docker Compose (Redpanda, DynamoDB Local, LocalStack,
  Flink session cluster) for daily development — $0.
- AWS envs are ephemeral: `terraform apply` to exercise, `terraform destroy`
  between sessions. Budget alarm at a fixed monthly cap is part of the
  Terraform baseline, created before any other AWS resource.
