# datacat Code Standards

This directory is the contract for every line of code in this repository.
The contract is unforgiving: a violation is a defect, regardless of whether
the code works, and "it's small" or "it's related" is never a defense. If
code conflicts with these documents, the code is wrong.

| Document | Covers |
|---|---|
| [twelve-factor.md](twelve-factor.md) | How datacat implements all 12 factors of https://12factor.net - binding |
| [performance.md](performance.md) | The Big-O and allocation discipline - binding, hot paths are unforgiving |
| [go.md](go.md) | Go services: layout, hexagonal layers, config, errors, HTTP, concurrency, testing |
| [java.md](java.md) | Spring Boot service and Flink job: OOP, DI, JPA, operators, checkpoints |
| [react.md](react.md) | Dashboard SPA: structure, components, hooks, server state |
| [infra.md](infra.md) | Terraform, Docker, Kubernetes, CI/CD |

## Universal rules (all languages)

### 1. One task per unit

Every function/method does exactly one thing. Two responsibilities means two
functions - always, immediately, not "when it grows". The test: you can name
it without "and". Depth of extraction is never an argument against
extraction; layers of small, named units are the goal, not a cost.

### 2. One kind per file (file taxonomy)

Every declaration has a KIND, and a file holds exactly one kind, serving
exactly one purpose. A declaration of a different kind than the file's
purpose MUST move to its own file - exported or unexported, large or small.
"It's only used here" and "it's related" do not exempt it.

The kinds, and where each lives:

| Kind | Examples | Lives in |
|---|---|---|
| Domain concept | `Verdict`, `Decision`, `Classification` | one file per concept in `domain/` (model layer) |
| Wire shape / DTO | `decisionResponse`, request records, event schemas | its own shapes file (`responses.go`, `*_dto`, schema file) - never inside behavior files |
| Interface / port | `DecisionStore`, `SignalScorer` | the ports/interface file of its consumer |
| Mapper / codec | `toDecisionResponse`, `decodeVerdict`, `eventFrom` | its own mapper/codec file - conversion is a kind, not a helper |
| Behavior / implementation | handlers, stores, middleware flows, operators | one implementation per file, named after it |
| Configuration | config shape, env loading | config files - shape and loading in separate files |
| Errors | sentinel errors, error types | the errors file of the layer |
| Registry / constants | `Scorers.all()`, action tables | its own registry file |
| Wiring | composition roots, routers, module calls | `main`, router, env roots - wiring only, nothing else |

The single exception: a private leaf function that is a *step of the file's
one task and the same kind* stays (e.g. a `deny` step inside the gate
middleware). The moment a private function does conversion, defines a shape,
or serves another purpose, it is a different kind and moves. There is no
"unexported" loophole anywhere in these standards.

### 3. Orchestrators delegate, leaves compute

- **Orchestrators** coordinate a use case. Every line delegates to a named
  step. No inline computation, no nested conditionals. No exceptions.
- **Leaves** do the work: ≤ 15 lines target, 30 the absolute cap. A leaf at
  16 lines gets split before merge, not after.

### 4. Hexagonal layers, one-way dependencies

```
domain/     pure business logic, zero framework/infra imports
ports/      interfaces only - "what I need", not "how it's done"
app/        use cases (the orchestrators)
adapters/   Kafka, DynamoDB, HTTP - implementations of ports
```

Dependency rule: `adapters → ports → domain`, never the reverse. Domain
code MUST compile and unit-test with zero infrastructure. One infra import
in `domain/` is a defect, not a shortcut.

### 5. Dependency injection everywhere

All collaborators arrive through the constructor. No globals, no
singletons, no service locators, no `init()` side effects. Concrete types
meet ONLY in composition roots (`main.go`, Spring context, `main.jsx`,
Terraform env roots).

### 6. Small interfaces, defined at the consumer

Interfaces have 1–3 methods and live with the code that *uses* them.
A 4-method interface requires a written justification in its doc comment;
5 is the cap.

### 7. Immutability by default

Domain objects are immutable. Transformations return new values. Mutation
is permitted only inside a single function's local scope or where the
language idiom demands it (documented per language, per case).

### 8. Extraction is mandatory, duplication is forbidden

Shared logic lives in the shared layer (`pkg/` in Go, common packages in
Java, `lib/` in React). A generic helper inside a service is a defect: if
it is generic, it belongs to the shared layer; if it is not generic, it
must be named for its specific task. (YAGNI still bounds *speculative*
abstraction: do not build for imagined callers - but organizing existing
code by kind is never speculative.)

**No duplicate function or method, anywhere.** Two functions that do the
same thing are one function in the wrong number of places. This is
unforgiving and has no size exception:

- **Byte-identical bodies** are a build-stopping defect. The second copy is
  deleted and every caller points at the one that remains.
- **Same logic, cosmetic differences** (renamed locals, reordered
  independent statements, a different literal that should be a parameter)
  counts as duplicate. Extract one function; pass the difference as an
  argument.
- **Same shape across languages that share a contract** (the Go signature
  base and the Flink signature base, a wire struct and its mirror) is
  allowed ONLY because the languages cannot import each other; a test must
  pin them together (see the wire-format tests), and a comment on each names
  its twin.
- **Placement of the one copy** follows the layer rules: used by two
  services -> `pkg/`; used across one service's packages -> that service's
  lowest shared package; used by one type -> a private method on it. Never
  copy to avoid an import.
- A test helper repeated across test files is duplication too: it moves to a
  shared `_test` helper, one definition.

The check before writing any function: does a function that already does
this exist? If yes, call it. If it *almost* does, make the difference a
parameter. Writing the second copy is never the answer.

### 9. Errors are handled, wrapped, and never swallowed

Every error is handled meaningfully or propagated with added context.
Empty catch blocks, `_ = err`, and log-and-forget on errors that matter
are defects. An intentionally ignored error carries a comment saying why,
every time.

### 10. Size budgets - caps, not suggestions

| Unit | Target | Hard cap |
|---|---|---|
| Leaf function | 15 lines | 30 lines |
| File | 200–400 lines | 800 lines |
| React component | 100 lines | 150 lines |
| Interface | 1–3 methods | 5 methods |
| Function parameters | ≤ 3 | 5 (then introduce a shape) |
| Nesting depth | 2 | 4 |

Exceeding a hard cap blocks the merge. Exceeding a target requires the
split to be the next edit in the same change.

### 11. Performance is a standard, not an afterthought

[performance.md](performance.md) is binding. Hot paths are O(1)/O(log n)
per item with no avoidable allocations; accidental quadratic behavior is a
defect anywhere; every super-linear function carries a complexity comment.

### 12. Tests are part of the definition of done

- TDD for `domain/` and `app/` layers: test first, red, green, refactor.
- Table-driven tests (Go), slice tests (Spring), Testing Library (React).
- Coverage floor: 80% on domain and app layers.
- Fakes are hand-written implementations of ports (Go); Mockito is the
  Java idiom; MSW at the fetch boundary in React.
- Adapters get integration tests against local containers.

### 13. Naming

- Names say what, types say how. Booleans read as predicates.
- Packages/modules are nouns, functions are verbs; no invented
  abbreviations (industry-standard ID, URL, TLS are fine).
- One concept per file and the file is named after it - see rule 2.

### 14. Comments: one line, plain words

Every comment is exactly one line, high level, in simple words. It says what
the thing is for, never how it works. No em dashes anywhere in the repo,
code or docs. A comment that needs a second line means the code needs a
better name or a smaller function. Rewrite the code, not the comment.

### 15. Declarations on one line, air inside methods

Function and method signatures, record headers, and constructor calls never
wrap parameters onto continuation lines, whatever the length. A call too
long for one line gets named locals first, then a one line call. Builder
chains and struct literals with named fields may stack. Inside any method
longer than three lines, blank lines separate setup, action, and result.

### 16. Git

- Commits: `<type>: <description>` - lowercase imperative, ≤ 72 chars.
  Types: `feat`, `fix`, `refactor`, `docs`, `test`, `chore`, `perf`, `ci`.
- Branches: `<author>/<ticket>-<description>`.
- Every PR: Goal, Summary, Design decisions, Edge cases, Files changed, Test plan.
- No secrets in the repo, ever. gitleaks gates every push.
