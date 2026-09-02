package io.datacat.classifier.config;

/**
 * All job parameters come from the environment (twelve-factor III) and are
 * validated here; a missing variable fails the job at submission, not
 * mid-run. The record owns its invariants; generic env reading lives in Env.
 */
public record JobConfig(String kafkaBrokers, String requestsTopic, String verdictsTopic,
		String consumerGroup, String checkpointUri, long checkpointIntervalMs,
		long windowSeconds, long slideSeconds) {

	public JobConfig {
		require(kafkaBrokers, "KAFKA_BROKERS");
		require(requestsTopic, "REQUESTS_TOPIC");
		require(verdictsTopic, "VERDICTS_TOPIC");
		require(consumerGroup, "CONSUMER_GROUP");
		require(checkpointUri, "CHECKPOINT_URI");
		requirePositive(windowSeconds, "WINDOW_SECONDS");
		requirePositive(slideSeconds, "SLIDE_SECONDS");
	}

	public static JobConfig fromEnv() {
		return new JobConfig(
				System.getenv("KAFKA_BROKERS"),
				System.getenv("REQUESTS_TOPIC"),
				System.getenv("VERDICTS_TOPIC"),
				Env.string("CONSUMER_GROUP", "classifier"),
				System.getenv("CHECKPOINT_URI"),
				Env.longValue("CHECKPOINT_INTERVAL_MS", 30_000L),
				Env.longValue("WINDOW_SECONDS", 300L),
				Env.longValue("SLIDE_SECONDS", 30L));
	}

	private static void require(String value, String name) {
		if (value == null || value.isBlank()) {
			throw new IllegalStateException("config: " + name + " is required");
		}
	}

	private static void requirePositive(long value, String name) {
		if (value <= 0) {
			throw new IllegalStateException("config: " + name + " must be positive");
		}
	}
}
