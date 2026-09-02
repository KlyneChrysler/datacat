# Java Standards (policy-service, classifier-job)

Java 21 language level (Flink's supported ceiling governs; verify at scaffold —
local JDK 24 can compile with `--release 21`). Full OOP here: this is where
class design, DI containers, and polymorphism are the idiom. Gradle with
version catalogs; Spotless + Checkstyle + Error Prone are CI gates.

---

# Part 1 — policy-service (Spring Boot)

## Layout

```
services/policy-service/src/main/java/io/datacat/policy/
├── PolicyServiceApplication.java
├── api/            controllers + request/response records (adapters, HTTP)
├── application/    services — use-case orchestrators
├── domain/         entities-as-behavior, value objects, domain exceptions
├── infrastructure/ JPA repositories, Kafka publishers (adapters, out)
└── config/         @ConfigurationProperties, bean wiring
```

Same hexagonal rule as Go: `api` and `infrastructure` depend on `application`
and `domain`; never the reverse. `domain` imports no Spring, no JPA annotations
where avoidable (persistence entities live in `infrastructure` and map to/from
domain objects).

## Configuration (factor III) — typed, validated at startup

No `@Value` scattered through classes. One immutable record per concern:

```java
// config/PolicyProperties.java
@ConfigurationProperties(prefix = "datacat.policy")
@Validated
public record PolicyProperties(
		@NotBlank String rulesTopic,
		@NotBlank String kafkaBrokers,
		@Positive int maxRulesPerTenant
) {}
```

```yaml
# application.yml — env indirection only, no literals that vary by deploy
datacat:
  policy:
    rules-topic: ${RULES_TOPIC}
    kafka-brokers: ${KAFKA_BROKERS}
    max-rules-per-tenant: ${MAX_RULES_PER_TENANT:100}
```

Startup fails if a required variable is missing — that is the desired behavior.

## Dependency injection

- Constructor injection ONLY. `@Autowired` on fields is forbidden.
- One constructor per class; Spring wires it without annotations.
- Collaborators are interfaces where a second implementation is plausible
  (repositories, publishers); concrete where not (YAGNI applies).

## Controller → Service → Repository

Controllers are thin adapters: parse, delegate, map. Records for DTOs.

```java
// api/RuleController.java
@RestController
@RequestMapping("/v1/rules")
public class RuleController {

	private final RuleService ruleService;

	public RuleController(RuleService ruleService) {
		this.ruleService = ruleService;
	}

	@PostMapping
	@ResponseStatus(HttpStatus.CREATED)
	public RuleResponse create(@Valid @RequestBody CreateRuleRequest request) {
		Rule rule = ruleService.create(request.toCommand());
		return RuleResponse.from(rule);
	}

	@GetMapping("/{id}")
	public RuleResponse get(@PathVariable UUID id) {
		return RuleResponse.from(ruleService.get(id));
	}
}
```

```java
// api/CreateRuleRequest.java — validation at the boundary, immutable record
public record CreateRuleRequest(
		@NotBlank String name,
		@NotNull Classification appliesTo,
		@NotNull ActionType action,
		@Min(1) @Max(10_000) int rateLimitPerMinute
) {
	public CreateRuleCommand toCommand() {
		return new CreateRuleCommand(name, appliesTo, action, rateLimitPerMinute);
	}
}
```

```java
// application/RuleService.java — orchestrator: each line delegates
@Service
public class RuleService {

	private final RuleRepository rules;
	private final RuleEventPublisher events;

	public RuleService(RuleRepository rules, RuleEventPublisher events) {
		this.rules = rules;
		this.events = events;
	}

	@Transactional
	public Rule create(CreateRuleCommand command) {
		Rule rule = Rule.from(command);
		Rule saved = rules.save(rule);
		events.publishCreated(saved);
		return saved;
	}

	@Transactional(readOnly = true)
	public Rule get(UUID id) {
		return rules.findById(id).orElseThrow(() -> new RuleNotFoundException(id));
	}
}
```

## Domain objects carry behavior, not just data

```java
// domain/Rule.java — invariants enforced in the factory, object is immutable
public record Rule(UUID id, String name, Classification appliesTo,
		ActionType action, int rateLimitPerMinute, Instant createdAt) {

	public static Rule from(CreateRuleCommand cmd) {
		requireCompatible(cmd.appliesTo(), cmd.action());
		return new Rule(UUID.randomUUID(), cmd.name(), cmd.appliesTo(),
				cmd.action(), cmd.rateLimitPerMinute(), Instant.now());
	}

	public boolean matches(Verdict verdict) {
		return verdict.classification() == appliesTo;
	}

	private static void requireCompatible(Classification c, ActionType a) {
		if (c == Classification.HUMAN && a == ActionType.BLOCK) {
			throw new IncompatibleRuleException(c, a);
		}
	}
}
```

## Errors — one global handler, no leaked internals

```java
// api/ApiExceptionHandler.java
@RestControllerAdvice
public class ApiExceptionHandler {

	private static final Logger log = LoggerFactory.getLogger(ApiExceptionHandler.class);

	@ExceptionHandler(RuleNotFoundException.class)
	@ResponseStatus(HttpStatus.NOT_FOUND)
	public ErrorResponse notFound(RuleNotFoundException e) {
		return ErrorResponse.of("rule not found");
	}

	@ExceptionHandler(Exception.class)
	@ResponseStatus(HttpStatus.INTERNAL_SERVER_ERROR)
	public ErrorResponse unexpected(Exception e) {
		log.error("unhandled exception", e);       // detail to logs (factor XI)
		return ErrorResponse.of("internal error"); // generic to client
	}
}
```

## JPA rules

- Persistence entities live in `infrastructure/`, mapped to/from domain records
  by a dedicated mapper class — JPA annotations never touch `domain/`.
- No bidirectional relationships unless proven necessary; prefer IDs across
  aggregates.
- `FetchType.LAZY` everywhere; queries that need more use explicit
  `@EntityGraph` or JPQL fetch joins. N+1 is a review-blocking defect.
- Flyway owns the schema; `ddl-auto: validate` only. Migrations run as one-off
  Kubernetes Jobs (factor XII).

## Testing

- Domain: plain JUnit 5, no Spring context, TDD.
- Application: JUnit 5 + Mockito for ports (Java idiom permits Mockito, unlike Go).
- API: `@WebMvcTest` slices; persistence: `@DataJpaTest` against Testcontainers
  Postgres. No full-context `@SpringBootTest` unless testing wiring itself.

```java
@Test
void createPublishesEventAfterSave() {
	RuleRepository rules = mock(RuleRepository.class);
	RuleEventPublisher events = mock(RuleEventPublisher.class);
	when(rules.save(any())).thenAnswer(inv -> inv.getArgument(0));
	RuleService service = new RuleService(rules, events);

	Rule created = service.create(validCommand());

	InOrder inOrder = inOrder(rules, events);
	inOrder.verify(rules).save(any(Rule.class));
	inOrder.verify(events).publishCreated(created);
}
```

---

# Part 2 — classifier-job (Apache Flink)

## Layout — one class per operator

```
services/classifier-job/src/main/java/io/datacat/classifier/
├── ClassifierJob.java        job graph assembly ONLY (the composition root)
├── config/                   env-driven job parameters
├── model/                    RequestEvent, SessionFeatures, Verdict (POJOs)
├── operators/
│   ├── SessionFeatureAggregator.java    keyed window aggregate
│   ├── TimingRegularityScorer.java      one signal, one class
│   ├── NavigationEntropyScorer.java
│   ├── FingerprintConsistencyScorer.java
│   └── VerdictAssembler.java            combines scores → verdict
├── serde/                    Kafka (de)serialization schemas
└── sinks/                    verdict topic, DynamoDB, S3 archive
```

The signal-scorer pattern from the Go standard repeats here: adding a detection
signal = adding one operator class + one line in `ClassifierJob`. Existing
operators are never edited (open/closed).

## Job assembly reads like the dataflow diagram

```java
// ClassifierJob.java — orchestrator: every statement is one pipeline step
public final class ClassifierJob {

	public static void main(String[] args) throws Exception {
		JobConfig config = JobConfig.fromEnv();
		StreamExecutionEnvironment env = Environments.checkpointed(config);

		DataStream<RequestEvent> events = Sources.requestEvents(env, config);

		DataStream<SessionFeatures> features = events
				.keyBy(RequestEvent::sessionId)
				.window(SlidingEventTimeWindows.of(Duration.ofMinutes(5), Duration.ofSeconds(30)))
				.aggregate(new SessionFeatureAggregator());

		DataStream<Verdict> verdicts = features
				.keyBy(SessionFeatures::sessionId)
				.process(new VerdictAssembler(Scorers.all()));

		verdicts.sinkTo(Sinks.verdictTopic(config));
		verdicts.sinkTo(Sinks.verdictStore(config));
		events.sinkTo(Sinks.rawArchive(config)); // S3, Parquet

		env.execute("datacat-classifier");
	}
}
```

## Operator rules

- Every operator is its own class with a single responsibility and its own
  unit tests (operators are plain classes — test them without a cluster).
- Keyed state is declared in `open()`, typed via `ValueStateDescriptor` /
  `MapStateDescriptor`, and every state has a TTL (sessions end; state must
  not grow forever).
- Event time + watermarks, never processing time, for anything user-visible:
  `WatermarkStrategy.forBoundedOutOfOrderness(Duration.ofSeconds(10))`.
- Late events go to a side output and are counted, not silently dropped.

```java
// operators/VerdictAssembler.java
public final class VerdictAssembler extends KeyedProcessFunction<String, SessionFeatures, Verdict> {

	private final List<SignalScorer> scorers;
	private transient ValueState<Classification> lastClass;

	public VerdictAssembler(List<SignalScorer> scorers) {
		this.scorers = List.copyOf(scorers);
	}

	@Override
	public void open(OpenContext ctx) {
		ValueStateDescriptor<Classification> descriptor =
				new ValueStateDescriptor<>("last-class", Classification.class);
		descriptor.enableTimeToLive(StateTtlConfig.newBuilder(Duration.ofHours(1)).build());
		lastClass = getRuntimeContext().getState(descriptor);
	}

	@Override
	public void processElement(SessionFeatures features, Context ctx, Collector<Verdict> out)
			throws Exception {
		ScoreCard scores = scoreAll(features);
		Verdict verdict = scores.toVerdict(features.sessionId());
		emitIfChanged(verdict, out);
	}

	private ScoreCard scoreAll(SessionFeatures features) {
		return scorers.stream()
				.map(scorer -> scorer.score(features))
				.collect(ScoreCard.collector());
	}

	private void emitIfChanged(Verdict verdict, Collector<Verdict> out) throws Exception {
		if (verdict.classification() != lastClass.value()) {
			lastClass.update(verdict.classification());
			out.collect(verdict);
		}
	}
}
```

## Checkpointing — non-negotiable settings

```java
// config/Environments.java
public static StreamExecutionEnvironment checkpointed(JobConfig config) {
	StreamExecutionEnvironment env = StreamExecutionEnvironment.getExecutionEnvironment();
	env.enableCheckpointing(config.checkpointIntervalMs(), CheckpointingMode.EXACTLY_ONCE);
	CheckpointConfig cp = env.getCheckpointConfig();
	cp.setCheckpointStorage(config.checkpointUri()); // s3://... in AWS, file:// locally
	cp.setMinPauseBetweenCheckpoints(config.checkpointIntervalMs() / 2);
	cp.setExternalizedCheckpointRetention(RETAIN_ON_CANCELLATION);
	return env;
}
```

Kafka source uses committed offsets + checkpoints (exactly-once within Flink);
sinks to external stores are idempotent (verdict keyed by session ID), so
end-to-end behavior is effectively-once.

## Testing Flink

- Operators: plain unit tests + `KeyedOneInputStreamOperatorTestHarness` for
  state/timer behavior.
- Job topology: `MiniClusterWithClientResource` integration test with bounded
  test sources — runs in CI, no real Kafka needed.
- A golden-path test replays a recorded human session and a scraper session
  and asserts the classifier separates them; this doubles as the accuracy
  regression gate.
