# Twelve-Factor Compliance — https://12factor.net

Every datacat service implements all 12 factors. This table is binding; the
per-language docs show the code that enforces each factor.

## I. Codebase — one codebase tracked in revision control, many deploys

One monorepo (`datacat/`), one git history. Each service is a separately
deployable app *within* the codebase; shared code lives in `pkg/` and Java
`common` modules — shared code is a library, never a copy-paste. Deploys
(local kind, staging EKS, prod EKS) differ only in config, never in code.

## II. Dependencies — explicitly declare and isolate

- Go: `go.mod`/`go.sum` per service; builds with `CGO_ENABLED=0`; never relies
  on system-wide tools at runtime.
- Java: Gradle version catalogs; no dependency without an explicit version.
- React: `package.json` with lockfile committed.
- Nothing shells out to implicit system binaries. Container images are the
  isolation boundary: distroless, containing exactly the binary and nothing else.

## III. Config — store config in the environment

- All config that varies between deploys (broker addresses, table names, ports,
  credentials, log level) comes from environment variables. Zero config files
  baked into images, zero constants like `if env == "prod"`.
- Every service validates its full config **at startup** and crashes loudly if
  anything is missing (see `go.md` § Config, `java.md` § Configuration).
- In Kubernetes, env comes from ConfigMaps (non-secret) and Secrets (secret),
  both created by Terraform/Helm — never hand-edited.
- The React build receives only its API base URL; everything else is served by
  the backend.

## IV. Backing services — treat as attached resources

Kafka, DynamoDB, S3, Postgres (policy-service) are addressed **only** by URL/
name from config and accessed **only** through ports (interfaces). Swapping
local Redpanda for MSK, or local DynamoDB for the real one, is a config change
plus at most one adapter — app and domain layers never know the difference.

## V. Build, release, run — strictly separate stages

- **Build**: CI produces an immutable image tagged with the git SHA.
- **Release**: image + environment config combine at deploy time (Helm values
  per environment). Every release is identifiable (SHA + config version).
- **Run**: containers execute the released artifact. No builds at runtime, no
  `kubectl edit`, no SSH-and-patch. Rollback = redeploy previous release.

## VI. Processes — stateless, share-nothing

Services keep no state between requests. Session verdicts live in DynamoDB,
stream state lives in Flink's checkpointed state (backed by S3), events live in
Kafka. Any pod can be killed at any moment and a replacement picks up cleanly.
No sticky sessions, no local file caches that matter.

## VII. Port binding — export services via port binding

Every service is self-contained and binds its own port from `PORT` env (Go
`net/http`, Spring embedded server). No service assumes an external web server.
Kubernetes Services route to those ports; health endpoints (`/healthz`,
`/readyz`) bind on the same listener.

## VIII. Concurrency — scale out via the process model

Scaling is horizontal: more pods (HPA on CPU/RPS), more Kafka partitions, more
Flink task slots. Internal concurrency (goroutines, Flink operators) is for
throughput within one process, never a substitute for scaling out. Nothing
daemonizes; the process is the unit.

## IX. Disposability — fast startup, graceful shutdown

- Startup: milliseconds for Go, and Spring is kept lean; a pod is `ready` only
  after its dependencies are reachable (readiness probe).
- Shutdown: every service traps SIGTERM, stops accepting work, drains in-flight
  requests/messages, commits Kafka offsets, then exits within the pod's
  `terminationGracePeriodSeconds`. See `go.md` § Graceful shutdown.
- Kafka consumers are idempotent, so a crash mid-batch is safe (at-least-once
  delivery + idempotent handling = effectively-once).

## X. Dev/prod parity — keep development, staging, production similar

The local stack (Docker Compose: Redpanda, DynamoDB Local, LocalStack, Flink)
runs the same images and the same env-var contract as EKS. Same Terraform
modules provision staging and prod with different variables. The gap between
"works locally" and "works on EKS" is config, not code.

## XI. Logs — treat logs as event streams

Every process writes structured JSON lines to **stdout only**. No log files, no
log rotation in-app, no in-process log shipping. The platform (container
runtime → CloudWatch / Fluent Bit) owns routing and retention. Log schema:
`ts`, `level`, `service`, `msg`, plus contextual keys (`session_id`,
`trace_id`). Never log secrets or full payloads.

## XII. Admin processes — run as one-off processes

Migrations (Flyway for policy-service), backfills, and topic creation run as
one-off Kubernetes Jobs using the *same image* and *same config* as the
service — never from a laptop against prod, never as code paths hidden inside
the server startup.
