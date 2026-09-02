# datacat Code Standards

This directory is the contract for every line of code in this repository. If code
conflicts with these documents, the code is wrong.

| Document | Covers |
|---|---|
| [twelve-factor.md](twelve-factor.md) | How datacat implements all 12 factors of https://12factor.net — binding, not aspirational |
| [go.md](go.md) | Go services: layout, hexagonal layers, config, errors, HTTP, concurrency, testing |
| [java.md](java.md) | Spring Boot service and Flink job: OOP, DI, JPA, operators, checkpoints |
| [react.md](react.md) | Dashboard SPA: structure, components, hooks, server state |
| [infra.md](infra.md) | Terraform, Docker, Kubernetes, CI/CD |

## Universal rules (all languages)

### 1. One task per unit

Every function/method does exactly one thing. Two responsibilities means two
functions. The test: you can name it without "and".

### 2. Orchestrators delegate, leaves compute

Two kinds of functions exist, and each has its own budget:

- **Orchestrators** coordinate a use case. Every line delegates to a named step.
  They read like prose. No inline computation, no nested conditionals.
- **Leaves** do the actual work. Target ≤ 15 lines, hard cap 30. A leaf that
  outgrows its budget gets split, and the split parts get extracted.

Do not shred leaves into one-liners for their own sake — a computation that
belongs together stays together, capped at 30 lines.

### 3. Hexagonal layers, one-way dependencies

```
domain/     pure business logic, zero framework/infra imports
ports/      interfaces only — "what I need", not "how it's done"
app/        use cases (the orchestrators)
adapters/   Kafka, DynamoDB, HTTP, filesystem — implementations of ports
```

Dependency rule: `adapters → ports → domain`. Never the reverse. Domain code
must compile and unit-test with zero infrastructure.

### 4. Dependency injection everywhere

All collaborators arrive through the constructor. No globals, no singletons, no
service locators, no `init()` side effects. The composition root (`main.go`,
Spring context, React entry) is the ONLY place where concrete types meet.

### 5. Small interfaces, defined at the consumer

Interfaces have 1–3 methods and live next to the code that *uses* them, not the
code that implements them. This is what keeps everything swappable.

### 6. Immutability by default

Domain objects are immutable. Transformations return new values. Mutation is
permitted only inside a single function's local scope or where the language
idiom demands it (documented per-language).

### 7. Extract helpers into dedicated packages — on the second use

Shared logic lives in `pkg/` (Go), `common` packages (Java), `lib/` (React).
Extraction is mandatory when a function does two things or when a second caller
appears. Do NOT pre-build abstractions for imagined callers (YAGNI).

### 8. Errors are handled, wrapped, and never swallowed

Every error is either handled meaningfully or propagated with added context.
Empty catch blocks and `_ = err` are forbidden. User-facing surfaces return
friendly messages; logs carry the detail.

### 9. Size budgets

| Unit | Target | Hard cap |
|---|---|---|
| Leaf function | 15 lines | 30 lines |
| File | 200–400 lines | 800 lines |
| React component | 100 lines | 150 lines |
| Interface | 1–3 methods | 5 methods |
| Function parameters | ≤ 3 | 5 (then introduce a struct/record) |
| Nesting depth | 2 | 4 |

### 10. Tests are part of the definition of done

- TDD for `domain/` and `app/` layers: test first, red, green, refactor.
- Table-driven tests (Go), slice tests (Spring), Testing Library (React).
- Coverage floor: 80% on domain and app layers.
- Tests use hand-written fakes implementing ports — not mocking frameworks —
  except where the ecosystem idiom is otherwise (documented per-language).

### 11. Naming

- Names say what, types say how. `verdictStore`, not `dynamoClient2`.
- Booleans read as predicates: `isVerified`, `hasExpired`.
- Packages/modules are nouns, functions are verbs, no abbreviations that
  aren't industry-standard (ID, URL, TLS are fine).
- One exported concept per file; the file is named after it.

### 12. Git

- Commits: `<type>: <description>` — lowercase imperative, ≤ 72 chars.
  Types: `feat`, `fix`, `refactor`, `docs`, `test`, `chore`, `perf`, `ci`.
- Branches: `<author>/<ticket>-<description>`.
- Every PR: Goal, Summary, Design decisions, Edge cases, Files changed, Test plan.
- No secrets in the repo, ever. Pre-commit scan enforced in CI.
