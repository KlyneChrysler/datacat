# datacat

Real-time traffic classification platform: tells **humans**, **verified AI agents**, and **abusive automation** apart on live web traffic, then enforces policy on the verdict.

## The problem

Automated traffic is now a huge share of the web, and "bots bad, humans good" has collapsed:

- AI agents transact on behalf of real customers. Blocking them loses revenue; ignoring scrapers burns money.
- User-Agent strings are spoofed, CAPTCHAs are solved by the agents themselves, IP reputation fails against residential proxies.
- Sites cannot price, police, or serve what they cannot classify.

datacat answers, per session, within seconds: **human / verified agent / unverified automation / abusive** — using behavioral and protocol signals that are hard to fake — and turns verdicts into enforcement (allow, rate-limit, challenge, block) and analytics.

## Architecture

```
Client traffic
      |
edge-proxy (Go, EKS) ──── emits request events + sensor beacons
      |                    (timing, header order, TLS fingerprint)
      v
   Kafka ──── events partitioned by session ID
      |
classifier-job (Flink/Java) ──── keyed state per session, windows over
      |                          timing variance, navigation entropy,
      |                          fingerprint consistency → verdict stream
      +──> DynamoDB ──── session verdicts, fingerprint reputation
      +──> S3 ─────────── raw event archive (Parquet)
      +──> Kafka verdicts topic ──> enforcement (Go): block/limit/challenge
      |
Step Functions ──── agent-verification workflow (register → prove
      |             ownership → issue credential → periodic re-check)
policy-service (Spring Boot) ──── rules CRUD, JPA
dashboard (React) ──── live human/agent split, session explorer
Lambda + API Gateway ──── public verdict API
traffic-sim (Go) ──── human / polite-agent / scraper generators
```

Everything is provisioned by Terraform, deployed to EKS via GitHub Actions, observed through CloudWatch.

## Repository layout

```
datacat/
├── services/
│   ├── edge-proxy/          Go   — request interception, event emission
│   ├── enforcement/         Go   — consumes verdicts, applies actions
│   ├── traffic-sim/         Go   — traffic generators (test harness)
│   ├── policy-service/      Java — Spring Boot rules CRUD
│   └── classifier-job/      Java — Flink streaming classification
├── pkg/                     shared Go libraries (events, kafkax, httpx, obsx)
├── web/dashboard/           React SPA
├── infra/terraform/         modules: network, eks, kafka, dynamo, lambda, sfn
├── deploy/k8s/              manifests / Helm charts per service
├── .github/workflows/       CI/CD pipelines
└── docs/
    └── standards/           THE code standard — read before writing anything
```

## Status

**The full loop works locally.** Traffic through the proxy becomes Kafka
events; the Flink job classifies sessions on behavior (timing regularity,
request rate, path entropy, user agent) and identity (Ed25519 trusted-agent
signatures); enforcement turns verdicts into decisions; the gate applies
them — block (403), challenge (proof-of-work interstitial with signed
clearances), rate-limit (per-session token buckets), allow. The dashboard
charts classifications live. Terraform provisions the (local) infrastructure;
Helm charts deploy to kind; CI gates every push.

Try it: `docker compose up -d`, submit the classifier jar to the Flink
cluster (see `deploy/k8s/README.md` and compose comments), then
`docker compose --profile demo up traffic-sim` and watch the scraper get
caught.

Still ahead: the AWS phase (EKS/MSK/DynamoDB via the Terraform staging env),
TLS fingerprinting, and true wire-order header capture.

`docs/standards/` is the contract for all code in this repo — read
[docs/standards/README.md](docs/standards/README.md) before contributing.
