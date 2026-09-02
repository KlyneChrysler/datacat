# datacat — project instructions

Real-time traffic classification platform (human vs AI-agent vs abusive
automation). See `README.md` for architecture.

## Before writing ANY code

Read `docs/standards/` — it is the binding contract:

- `docs/standards/README.md` — universal rules (budgets, hexagonal layers, DI)
- `docs/standards/twelve-factor.md` — 12factor.net compliance, all 12 factors
- `docs/standards/go.md` — Go services
- `docs/standards/java.md` — Spring Boot + Flink
- `docs/standards/react.md` — dashboard
- `docs/standards/infra.md` — Terraform, Docker, K8s, CI/CD

Code that violates the standards is wrong even if it works.

## Hard rules (summary — full versions in standards)

- Hexagonal layers everywhere: `domain → ports → app → adapters`, dependencies
  point inward only. Domain imports zero infrastructure.
- Orchestrators delegate (one call per line); leaf functions ≤ 15 lines target,
  30 hard cap. One task per function — no "and" in names.
- Constructor injection only. No globals, no singletons. Concrete types meet
  only in composition roots (`main.go`, Spring context, `main.jsx`, env
  `main.tf`).
- Interfaces: 1–3 methods, defined at the consumer.
- Immutable domain objects; transformations return new values.
- Every error wrapped with context or handled; never swallowed.
- All config from environment variables, validated at startup, crash on missing.
- Structured JSON logs to stdout only. Never log secrets or payloads.
- Files 200–400 lines (800 cap). One exported concept per file.
- TDD on domain/app layers; 80% coverage floor; hand-written fakes in Go,
  Mockito allowed in Java, MSW at fetch boundary in React.

## Repo conventions

- Monorepo: Go services in `services/*` (own `go.mod` each), shared Go code in
  `pkg/*`, Java in `services/policy-service` + `services/classifier-job`,
  React in `web/dashboard`, Terraform in `infra/terraform`, manifests in
  `deploy/k8s`.
- Commits: `<type>: <description>`, lowercase imperative, ≤ 72 chars
  (`feat`, `fix`, `refactor`, `docs`, `test`, `chore`, `perf`, `ci`).
- Never guess dependency versions — resolve with the toolchain at add time
  (`go get`, Gradle catalog, `npm install`) and commit lockfiles.
- Local-first development: kind + Docker Compose. AWS is ephemeral
  (`terraform destroy` between sessions); budget alarm exists before any
  other AWS resource.

## Toolchain (verified on this machine, 2026-09)

Go 1.26.3 · Java 24 (compile with `--release 21` for Flink compat — verify at
scaffold) · Node 25 · Docker 29 · Terraform 1.15 · kubectl 1.34.
Not yet installed: kind, helm — install at scaffold phase.
