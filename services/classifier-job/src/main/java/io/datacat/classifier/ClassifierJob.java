package io.datacat.classifier;

import io.datacat.classifier.config.Environments;
import io.datacat.classifier.config.JobConfig;
import org.apache.flink.streaming.api.environment.StreamExecutionEnvironment;

/**
 * Job graph assembly only (the composition root). Sources, operators, and
 * sinks are added in the Flink phase; this skeleton proves config loading and
 * checkpoint setup compile against the pinned Flink version.
 */
public final class ClassifierJob {

	private ClassifierJob() {
	}

	public static void main(String[] args) {
		JobConfig config = JobConfig.fromEnv();
		StreamExecutionEnvironment env = Environments.checkpointed(config);
		throw new UnsupportedOperationException(
				"pipeline arrives in the Flink phase; env configured: " + env.getCheckpointConfig());
	}
}
