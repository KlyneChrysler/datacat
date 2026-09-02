# Performance Standards — Big-O and Allocation Discipline

Binding for all languages. The rules distinguish hot paths (unforgiving)
from cold paths (strict). "It's fast enough on my machine" is not an
argument; the complexity class is.

## What is a hot path

Code executed **per request, per event/record, or per stream element**:

- edge-proxy: the request path (gate lookup, observe, proxy), the recorder
  drain loop
- enforcement: `HandleVerdict`, decision lookup, Kafka record handling
- classifier-job: `processElement`, every `AggregateFunction.add`, serde
- pkg/kafkax, pkg/httpx middleware: everything
- dashboard: anything running per render frame or per poll tick

Everything else (window `getResult`, startup, admin endpoints, tests,
Terraform) is a cold path.

## Hot-path rules (unforgiving)

1. **O(1) or O(log n) per item.** A per-item cost of O(n) in session count,
   rule count, or history length is a defect. If it cannot be avoided, the
   function carries a written justification AND a bound on n — or it does
   not merge.
2. **Membership is a map/set lookup, never a scan.** `for` + compare to
   find one item is banned on hot paths.
3. **No allocations that scale with traffic where avoidable**: pre-size
   slices/maps when length is known, reuse buffers, no per-request
   compilation of regexes or templates (compile once at construction).
4. **No lock held across I/O.** Take the lock, read/write memory, release.
   Prefer RWMutex read locks or ownership transfer over exclusive locks.
5. **Bounded everything.** Every queue, cache, and accumulator on a hot
   path has an explicit cap with a defined overflow policy (drop-and-count,
   evict, backpressure). Unbounded growth is a defect even when "unlikely".
6. **String assembly uses a builder** (strings.Builder, StringBuilder,
   array join) — never `+=` in a loop.

## Cold-path rules (strict)

1. **No accidental quadratic.** Nested iteration over collections that can
   grow together is a defect unless the inner collection is provably
   bounded by a constant (state the constant in a comment).
2. **Every super-linear function carries a complexity comment** —
   `// O(n log n): sorts the window sample (bounded at 512)` — stating the
   class and the bound. No comment, no merge.
3. **Sorting is O(n log n) once, not O(n² ) by insertion**; top-k uses a
   heap, not a full sort, when k << n.

## Data-structure selection (all paths)

| Need | Use | Never |
|---|---|---|
| Membership / lookup by key | map / set / dict | linear scan |
| Latest-value per key | map with overwrite | append + scan-back |
| Top-k of n | heap of size k | full sort when k << n |
| FIFO with cap | bounded channel / ring buffer | unbounded slice |
| Counting by key | map[key]count | repeated filtering |
| Ordered iteration by time | already-ordered source, or one sort at the boundary | re-sorting per read |

## Verification

- Reviews check the complexity class of every touched hot-path function —
  it is part of correctness, not polish.
- Known bounds live in code: caps are named constants
  (`MAX_SAMPLES = 512`), never magic numbers.
- When a hot-path structure is contended or measured slow, the fix is an
  algorithmic or structural one first; micro-tuning without a complexity
  argument is rejected.

## Current codebase register (kept honest, update when it changes)

- `Gatekeeper.ActionFor`: O(1) map read under RLock — compliant.
- `Recorder.Record`: O(1) non-blocking enqueue, bounded 1024, drop-and-count — compliant.
- `FeatureMath.intervalCv` / `normalizedPathEntropy`: O(n log n)/O(n) on a
  window sample bounded at `FeatureAccumulator.MAX_SAMPLES = 512`; runs per
  window emission (cold), not per event — compliant with comment.
- `SessionFeatureAggregator.add`: O(1) append under cap — compliant.
- In-memory `DecisionStore`: O(1) map; DynamoDB adapter is O(1) point ops.
