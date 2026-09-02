# datacat

## What is datacat

datacat is a real time traffic classification platform. It watches live web
traffic and tells, per session, whether the visitor is a human, a verified AI
agent, unverified automation, or an abusive bot. It then acts on that verdict
at the edge: allow, rate limit, challenge with a proof of work page, or block.

It is built as an event driven system of small services:

```
Client traffic
      |
edge-proxy (Go) ......... observes requests, verifies agent signatures,
      |                   enforces decisions (403, 429, challenge page)
      v
   Kafka (Redpanda) ..... request events, keyed by session
      |
classifier-job (Flink) .. sliding windows, behavior signals plus identity,
      |                   emits a verdict when a classification changes
      v
enforcement (Go) ........ applies policy, stores decisions, publishes them
      |                   back to Kafka and serves a small HTTP API
      v
dashboard (React) ....... live chart of classifications
traffic-sim (Go) ........ demo personas: human, polite agent, scraper
```

Everything runs locally for free with Docker Compose. Helm charts deploy the
services to a kind cluster, and Terraform modules provision the (local)
infrastructure. The engineering rules live in `docs/standards/` and are
binding for every line of code.

## The problem it solves

A large share of web traffic is now automated, and the old line between bots
and humans has collapsed:

- AI agents shop and browse on behalf of real customers. Blocking them loses
  revenue, while a scraper hitting the same endpoint only burns money.
- The old signals are dead. User agent strings are spoofed, CAPTCHAs are
  solved by the same automation they try to stop, and IP reputation fails
  against residential proxies.
- When abuse happens, the response is usually manual: someone greps logs days
  later.

datacat answers the question a site actually needs answered, within seconds
and per session: who is this, and what should happen to their next request.
Humans pass untouched, agents that prove their identity with a signature get
a rate limited lane, suspicious sessions must solve a browser challenge, and
abusive sessions get blocked, all automatically.

## How to set it up

Prerequisites: Docker Desktop, Go 1.26+, Node 20+, and a JDK for Gradle
(the build auto provisions the right Java toolchains). No cloud account is
needed, everything runs on your machine.

1. Start the infrastructure and services:

```bash
docker compose up -d
docker compose up topic-init
```

This starts Redpanda (Kafka), a demo upstream site, the Flink cluster, the
edge proxy on port 8080, and enforcement on port 8081, then creates the
three topics. Note: topics vanish on `docker compose down`, so rerun
`topic-init` after every fresh start.

2. Build and submit the classifier job:

```bash
cd services/classifier-job && ./gradlew shadowJar && cd ../..
docker cp services/classifier-job/build/libs/classifier-job-0.0.1-SNAPSHOT.jar datacat-flink-jobmanager-1:/tmp/classifier-job.jar
docker exec datacat-flink-jobmanager-1 flink run -d /tmp/classifier-job.jar
```

3. Start the dashboard:

```bash
cd web/dashboard
npm install
echo "VITE_API_BASE_URL=http://localhost:8081" > .env
npm run dev
```

4. Send demo traffic with the three personas:

```bash
docker compose --profile demo up traffic-sim
```

5. Watch it work:

- Dashboard: http://localhost:5173 (all four classes fill in within a minute)
- Flink UI: http://localhost:8086
- Decisions API: `curl http://localhost:8081/v1/decisions/sim-scraper`
- Live verdicts: `docker exec datacat-redpanda-1 rpk topic consume datacat.verdicts`
- Feel the gate yourself: `curl -b dc_session=sim-scraper localhost:8080/x`
  returns 403 once the scraper is blocked, while your own browser session
  passes normally at http://localhost:8080

6. Stop everything:

```bash
docker compose down
```

## Credentials needed

No external accounts or paid services are required for the local setup.
The values below are environment variables, all with working local defaults
already wired into `docker-compose.yml`:

| Variable | Service | Required | Purpose |
|---|---|---|---|
| `CHALLENGE_SECRET` | edge-proxy | yes | Signs challenge tokens and clearance cookies. Any random string locally. In Kubernetes it must come from a Secret: `kubectl create secret generic edge-proxy-secrets --from-literal=CHALLENGE_SECRET=$(openssl rand -hex 32)` |
| `AGENT_KEYS` | edge-proxy | no | Trusted agent public keys as `keyid=hexpubkey`. The compose file ships a dev key pair |
| `AGENT_KEY_SEED` | traffic-sim | no | Derives the demo agent signing key matching `AGENT_KEYS` |
| `AWS_ACCESS_KEY_ID` and `AWS_SECRET_ACCESS_KEY` | enforcement | no | Only when enabling the DynamoDB store against DynamoDB Local, any dummy values work (`local` / `local`) |

Real AWS credentials are only needed for the future cloud deployment through
`infra/terraform/envs/staging`, which is optional and documented separately
in `infra/terraform/README.md`.

## License

Apache License 2.0, see [LICENSE](LICENSE).
