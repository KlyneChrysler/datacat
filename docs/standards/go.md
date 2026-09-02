# Go Standards (edge-proxy, enforcement, traffic-sim, pkg/*)

Go ≥ 1.26. Idiomatic Go is the dialect: composition and small interfaces, not
class hierarchies. `gofmt` + `go vet` + `golangci-lint` are CI gates.

## Service layout

```
services/enforcement/
├── cmd/server/main.go        composition root — wiring ONLY
├── internal/
│   ├── domain/               pure logic: Verdict, Action, Session
│   ├── ports/                consumer-defined interfaces
│   ├── app/                  use-case orchestrators
│   ├── adapters/
│   │   ├── kafka/            VerdictConsumer implementation
│   │   ├── dynamo/           DecisionStore implementation
│   │   └── httpapi/          handlers, middleware, router
│   └── config/               env loading + validation
├── Dockerfile
└── go.mod                    each service is its own module
```

Shared code lives in `pkg/` at repo root (`events`, `kafkax`, `httpx`, `obsx`,
`hashx`) and is imported like any third-party module. `internal/` is never
shared. A generic helper written inside a service is a defect — it moves to
`pkg/`.

## File taxonomy (rule 2, Go form) — unforgiving

Within every package, one kind per file:

```
internal/domain/          one concept per file:
  classification.go         Classification + its constants
  verdict.go                Verdict + its constructor and methods
  action.go / decision.go / policy.go
  errors.go                 sentinel errors only
internal/adapters/httpapi/
  router.go                 wiring only
  handlers.go               behavior only — NO type declarations
  responses.go              wire shapes only
  mapper.go                 domain → wire conversion only
internal/adapters/kafka/
  verdict_source.go         behavior (consume loop bridging)
  verdict_codec.go          wire ↔ domain conversion only
internal/config/
  config.go                 the Config shape and derived accessors
  load.go                   env reading, parsing, validation
```

Binding consequences, no exceptions for unexported names:

- A `struct` with JSON/dynamodbav tags inside a behavior file is a defect —
  wire and storage shapes get shapes files.
- A function converting between two representations (`toX`, `decodeX`,
  `eventFrom`) is a mapper/codec and gets a mapper/codec file, even with a
  single caller.
- Private state shapes supporting one type (`gateEntry`) still get their
  own file: shapes are a kind.
- Private leaf steps of the file's single task and kind stay
  (`deny` inside the gate middleware, `shutdown` inside the server).

## Config — factor III, validated at startup

Env only. Missing config crashes the process at boot, not at first use.

```go
// internal/config/config.go
package config

import (
	"fmt"
	"os"
	"time"
)

type Config struct {
	Port            string
	KafkaBrokers    string
	VerdictsTopic   string
	DecisionsTable  string
	ShutdownTimeout time.Duration
}

func Load() (Config, error) {
	cfg := Config{
		Port:            os.Getenv("PORT"),
		KafkaBrokers:    os.Getenv("KAFKA_BROKERS"),
		VerdictsTopic:   os.Getenv("VERDICTS_TOPIC"),
		DecisionsTable:  os.Getenv("DECISIONS_TABLE"),
		ShutdownTimeout: 10 * time.Second,
	}
	return cfg, cfg.validate()
}

func (c Config) validate() error {
	for name, v := range map[string]string{
		"PORT":            c.Port,
		"KAFKA_BROKERS":   c.KafkaBrokers,
		"VERDICTS_TOPIC":  c.VerdictsTopic,
		"DECISIONS_TABLE": c.DecisionsTable,
	} {
		if v == "" {
			return fmt.Errorf("config: %s is required", name)
		}
	}
	return nil
}
```

## Domain — immutable values, zero imports of infra

```go
// internal/domain/verdict.go
package domain

type Classification string

const (
	Human       Classification = "human"
	VerifiedBot Classification = "verified_agent"
	Unverified  Classification = "unverified_automation"
	Abusive     Classification = "abusive"
)

// Verdict is immutable: fields are set once by the constructor and never
// mutated. "Changes" return a new value.
type Verdict struct {
	SessionID  string
	Class      Classification
	Confidence float64
}

func NewVerdict(sessionID string, class Classification, confidence float64) (Verdict, error) {
	if sessionID == "" {
		return Verdict{}, ErrEmptySessionID
	}
	if confidence < 0 || confidence > 1 {
		return Verdict{}, ErrConfidenceRange
	}
	return Verdict{SessionID: sessionID, Class: class, Confidence: confidence}, nil
}

func (v Verdict) WithClass(c Classification) Verdict {
	v.Class = c // value receiver: mutates the copy, not the original
	return v
}
```

```go
// internal/domain/errors.go — sentinel errors live with the domain
package domain

import "errors"

var (
	ErrEmptySessionID  = errors.New("session id must not be empty")
	ErrConfidenceRange = errors.New("confidence must be within [0,1]")
	ErrVerdictNotFound = errors.New("verdict not found")
)
```

## Ports — small, consumer-defined

The interface lives in the package that *needs* the capability. Implementations
are elsewhere and are never referenced by name outside `main.go`.

```go
// internal/ports/ports.go
package ports

import (
	"context"

	"datacat/services/enforcement/internal/domain"
)

type DecisionStore interface {
	Save(ctx context.Context, d domain.Decision) error
	FindBySession(ctx context.Context, sessionID string) (domain.Decision, error)
}

type ActionApplier interface {
	Apply(ctx context.Context, d domain.Decision) error
}

type VerdictSource interface {
	// Consume blocks, invoking handle for each verdict until ctx is done.
	Consume(ctx context.Context, handle func(context.Context, domain.Verdict) error) error
}
```

## App — the orchestrator layer

Every public method reads like prose: one delegating call per step.

```go
// internal/app/enforcer.go
package app

import (
	"context"
	"fmt"

	"datacat/services/enforcement/internal/domain"
	"datacat/services/enforcement/internal/ports"
)

type Enforcer struct {
	policy  domain.Policy
	store   ports.DecisionStore
	applier ports.ActionApplier
}

func NewEnforcer(policy domain.Policy, store ports.DecisionStore, applier ports.ActionApplier) *Enforcer {
	return &Enforcer{policy: policy, store: store, applier: applier}
}

// HandleVerdict is an orchestrator: every line delegates.
func (e *Enforcer) HandleVerdict(ctx context.Context, v domain.Verdict) error {
	decision := e.policy.Decide(v)
	if err := e.store.Save(ctx, decision); err != nil {
		return fmt.Errorf("save decision for session %s: %w", v.SessionID, err)
	}
	if err := e.applier.Apply(ctx, decision); err != nil {
		return fmt.Errorf("apply %s for session %s: %w", decision.Action, v.SessionID, err)
	}
	return nil
}
```

## Adapters — one per backing service (factor IV)

```go
// internal/adapters/dynamo/decision_store.go
package dynamo

// DecisionStore implements ports.DecisionStore against DynamoDB. Compile-time
// check keeps the contract honest without coupling the port to this package.
type DecisionStore struct {
	client *dynamodb.Client
	table  string
}

var _ ports.DecisionStore = (*DecisionStore)(nil)

func NewDecisionStore(client *dynamodb.Client, table string) *DecisionStore {
	return &DecisionStore{client: client, table: table}
}

func (s *DecisionStore) Save(ctx context.Context, d domain.Decision) error {
	item, err := marshalDecision(d)
	if err != nil {
		return fmt.Errorf("marshal decision: %w", err)
	}
	if err := s.putItem(ctx, item); err != nil {
		return fmt.Errorf("dynamo put %s: %w", s.table, err)
	}
	return nil
}
```

## Errors

- Wrap with context at every boundary: `fmt.Errorf("doing X for %s: %w", id, err)`.
- Compare with `errors.Is` / extract with `errors.As` — never string matching.
- Sentinels in `domain/errors.go`; adapters translate infra errors into them
  (e.g. DynamoDB `ResourceNotFoundException` → `domain.ErrVerdictNotFound`).
- `_ = err` and empty error branches are forbidden. If an error is truly
  ignorable, a comment must say why.

## HTTP — handlers are thin adapters

Handlers parse, delegate, respond. Business logic never lives in a handler.

```go
// internal/adapters/httpapi/router.go
package httpapi

func NewRouter(h *Handlers, log *slog.Logger) http.Handler {
	mux := http.NewServeMux()
	mux.HandleFunc("GET /healthz", h.Health)
	mux.HandleFunc("GET /readyz", h.Ready)
	mux.HandleFunc("GET /v1/decisions/{sessionID}", h.GetDecision)
	return httpx.WithMiddleware(mux, httpx.RequestID(), httpx.Logging(log), httpx.Recover(log))
}
```

```go
// internal/adapters/httpapi/decisions.go
func (h *Handlers) GetDecision(w http.ResponseWriter, r *http.Request) {
	sessionID := r.PathValue("sessionID")
	decision, err := h.enforcer.Lookup(r.Context(), sessionID)
	switch {
	case errors.Is(err, domain.ErrVerdictNotFound):
		httpx.Error(w, http.StatusNotFound, "no decision for session")
	case err != nil:
		httpx.InternalError(w, r, err) // logs detail, returns generic message
	default:
		httpx.JSON(w, http.StatusOK, toDecisionResponse(decision))
	}
}
```

`pkg/httpx` owns the response envelope (`{"ok": bool, "data": ..., "error": ...}`),
middleware, and request validation helpers — written once, reused by every service.

## Composition root + graceful shutdown (factor IX)

`main.go` is wiring and lifecycle only. Nothing else.

```go
// cmd/server/main.go
package main

func main() {
	if err := run(); err != nil {
		slog.Error("fatal", "err", err)
		os.Exit(1)
	}
}

func run() error {
	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	cfg, err := config.Load()
	if err != nil {
		return err
	}
	log := obsx.NewLogger("enforcement") // slog JSON → stdout (factor XI)

	store := dynamo.NewDecisionStore(mustDynamoClient(ctx, cfg), cfg.DecisionsTable)
	applier := actions.NewApplier(log)
	enforcer := app.NewEnforcer(domain.DefaultPolicy(), store, applier)
	consumer := kafka.NewVerdictConsumer(cfg.KafkaBrokers, cfg.VerdictsTopic, log)
	server := httpx.NewServer(cfg.Port, httpapi.NewRouter(httpapi.New(enforcer), log))

	g, ctx := errgroup.WithContext(ctx)
	g.Go(func() error { return consumer.Consume(ctx, enforcer.HandleVerdict) })
	g.Go(func() error { return server.ListenAndServe(ctx, cfg.ShutdownTimeout) })
	return g.Wait()
}
```

Rules: trap SIGTERM via `signal.NotifyContext`; every long-running component
takes `ctx` and drains on cancellation; `httpx.NewServer` wraps
`http.Server.Shutdown` with the timeout.

## Concurrency

- Every blocking call takes a `context.Context` as the first parameter.
- `errgroup.Group` for fan-out with error propagation; raw `go func()` only
  with an explicit lifecycle owner.
- Channels are owned by the sender: the goroutine that writes, closes.
- Bounded everything: worker pools have a size, channels have a capacity,
  retries have a cap with jittered backoff (`pkg/kafkax` owns retry policy).
- No mutable shared state without a mutex or, preferably, ownership transfer.

```go
// traffic-sim: bounded fan-out, one task per function
func (s *Simulator) Run(ctx context.Context, personas []Persona) error {
	g, ctx := errgroup.WithContext(ctx)
	g.SetLimit(s.maxConcurrent)
	for _, p := range personas {
		g.Go(func() error { return s.runPersona(ctx, p) })
	}
	return g.Wait()
}
```

## Logging (factor XI)

`log/slog` JSON to stdout, built once in `pkg/obsx`, injected everywhere.
Contextual fields, never string interpolation:

```go
log.InfoContext(ctx, "decision applied",
	"session_id", d.SessionID, "action", d.Action, "class", d.Class)
```

Never log secrets, credentials, or full request bodies.

## Testing

- Table-driven tests, `t.Run` subtests, `t.Parallel()` where isolation allows.
- Fakes are hand-written structs implementing ports — no mock frameworks.
- Domain and app layers: TDD, ≥ 80% coverage, zero infrastructure needed.
- Adapters: integration tests against local containers (Redpanda, DynamoDB
  Local), tagged `//go:build integration`.

```go
// internal/app/enforcer_test.go
type fakeStore struct{ saved []domain.Decision; err error }

func (f *fakeStore) Save(_ context.Context, d domain.Decision) error {
	f.saved = append(f.saved, d)
	return f.err
}

func TestEnforcerHandleVerdict(t *testing.T) {
	tests := []struct {
		name       string
		verdict    domain.Verdict
		storeErr   error
		wantAction domain.Action
		wantErr    bool
	}{
		{name: "abusive session gets blocked", verdict: abusiveVerdict(t), wantAction: domain.Block},
		{name: "store failure propagates", verdict: abusiveVerdict(t), storeErr: errBoom, wantErr: true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			store := &fakeStore{err: tt.storeErr}
			e := app.NewEnforcer(domain.DefaultPolicy(), store, &fakeApplier{})

			err := e.HandleVerdict(context.Background(), tt.verdict)

			if tt.wantErr != (err != nil) {
				t.Fatalf("err = %v, wantErr = %v", err, tt.wantErr)
			}
			if !tt.wantErr && store.saved[0].Action != tt.wantAction {
				t.Errorf("action = %s, want %s", store.saved[0].Action, tt.wantAction)
			}
		})
	}
}
```

## Dependency versions

Never guess versions. At scaffold time, add dependencies with `go get` and let
the toolchain resolve the latest stable; `go.sum` pins them. Candidate
libraries (verify at scaffold): `twmb/franz-go` (Kafka), `aws-sdk-go-v2`,
`golang.org/x/sync/errgroup`.
