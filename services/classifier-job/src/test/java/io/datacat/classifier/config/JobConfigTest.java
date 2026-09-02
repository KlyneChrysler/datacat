package io.datacat.classifier.config;

import org.junit.jupiter.api.Test;

import static org.junit.jupiter.api.Assertions.assertEquals;
import static org.junit.jupiter.api.Assertions.assertThrows;

class JobConfigTest {

	private static JobConfig valid() {
		return new JobConfig("broker:9092", "requests", "verdicts", "classifier",
				"file:///tmp/cp", 30_000L, 300L, 30L);
	}

	@Test
	void validConfigConstructs() {
		JobConfig config = valid();

		assertEquals("broker:9092", config.kafkaBrokers());
		assertEquals(300L, config.windowSeconds());
	}

	@Test
	void missingBrokerIsRejected() {
		assertThrows(IllegalStateException.class,
				() -> new JobConfig("", "requests", "verdicts", "classifier",
						"file:///tmp/cp", 30_000L, 300L, 30L));
	}

	@Test
	void missingCheckpointUriIsRejected() {
		assertThrows(IllegalStateException.class,
				() -> new JobConfig("broker:9092", "requests", "verdicts", "classifier",
						null, 30_000L, 300L, 30L));
	}

	@Test
	void nonPositiveWindowIsRejected() {
		assertThrows(IllegalStateException.class,
				() -> new JobConfig("broker:9092", "requests", "verdicts", "classifier",
						"file:///tmp/cp", 30_000L, 0L, 30L));
	}
}
