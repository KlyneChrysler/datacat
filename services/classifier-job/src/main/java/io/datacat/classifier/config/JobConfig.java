package io.datacat.classifier.config;

/**
 * All job parameters come from the environment (twelve-factor III) and are
 * validated here; a missing variable fails the job at submission, not mid-run.
 */
public record JobConfig(String kafkaBrokers, String requestsTopic, String verdictsTopic,
		String checkpointUri, long checkpointIntervalMs) {

	public JobConfig {
		require(kafkaBrokers, "KAFKA_BROKERS");
		require(requestsTopic, "REQUESTS_TOPIC");
		require(verdictsTopic, "VERDICTS_TOPIC");
		require(checkpointUri, "CHECKPOINT_URI");
	}

	public static JobConfig fromEnv() {
		return new JobConfig(
				System.getenv("KAFKA_BROKERS"),
				System.getenv("REQUESTS_TOPIC"),
				System.getenv("VERDICTS_TOPIC"),
				System.getenv("CHECKPOINT_URI"),
				longFromEnv("CHECKPOINT_INTERVAL_MS", 30_000L));
	}

	private static void require(String value, String name) {
		if (value == null || value.isBlank()) {
			throw new IllegalStateException("config: " + name + " is required");
		}
	}

	private static long longFromEnv(String name, long fallback) {
		String raw = System.getenv(name);
		return raw == null || raw.isBlank() ? fallback : Long.parseLong(raw);
	}
}
