package io.datacat.classifier.config;

import org.apache.flink.configuration.CheckpointingOptions;
import org.apache.flink.configuration.Configuration;
import org.apache.flink.configuration.ExternalizedCheckpointRetention;
import org.apache.flink.core.execution.CheckpointingMode;
import org.apache.flink.streaming.api.environment.StreamExecutionEnvironment;


/** Non-negotiable checkpoint settings (docs/standards/java.md part 2). */
public final class Environments {

	private Environments() {
	}

	public static StreamExecutionEnvironment checkpointed(JobConfig config) {
		StreamExecutionEnvironment env = StreamExecutionEnvironment.getExecutionEnvironment(storage(config));
		env.enableCheckpointing(config.checkpointIntervalMs(), CheckpointingMode.EXACTLY_ONCE);
		env.getCheckpointConfig().setMinPauseBetweenCheckpoints(config.checkpointIntervalMs() / 2);
		env.getCheckpointConfig()
				.setExternalizedCheckpointRetention(ExternalizedCheckpointRetention.RETAIN_ON_CANCELLATION);
		return env;
	}

	private static Configuration storage(JobConfig config) {
		Configuration conf = new Configuration();
		conf.set(CheckpointingOptions.CHECKPOINT_STORAGE, "filesystem");
		conf.set(CheckpointingOptions.CHECKPOINTS_DIRECTORY, config.checkpointUri());
		return conf;
	}
}
