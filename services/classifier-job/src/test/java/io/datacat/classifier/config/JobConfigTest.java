package io.datacat.classifier.config;

import org.junit.jupiter.api.Test;

import static org.junit.jupiter.api.Assertions.assertEquals;
import static org.junit.jupiter.api.Assertions.assertThrows;

class JobConfigTest {

	@Test
	void validConfigConstructs() {
		JobConfig config = new JobConfig("broker:9092", "requests", "verdicts", "file:///tmp/cp", 30_000L);

		assertEquals("broker:9092", config.kafkaBrokers());
		assertEquals(30_000L, config.checkpointIntervalMs());
	}

	@Test
	void missingBrokerIsRejected() {
		assertThrows(IllegalStateException.class,
				() -> new JobConfig("", "requests", "verdicts", "file:///tmp/cp", 30_000L));
	}

	@Test
	void missingCheckpointUriIsRejected() {
		assertThrows(IllegalStateException.class,
				() -> new JobConfig("broker:9092", "requests", "verdicts", null, 30_000L));
	}
}
