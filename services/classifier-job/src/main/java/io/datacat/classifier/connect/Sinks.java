package io.datacat.classifier.connect;

import io.datacat.classifier.config.JobConfig;
import io.datacat.classifier.model.Verdict;
import io.datacat.classifier.serde.VerdictSerializationSchema;
import org.apache.flink.connector.kafka.sink.KafkaSink;

/** Outbound edges of the job graph, at least once by design. */
public final class Sinks {

	private Sinks() {
	}

	public static KafkaSink<Verdict> verdicts(JobConfig config) {
		return KafkaSink.<Verdict>builder()
				.setBootstrapServers(config.kafkaBrokers())
				.setRecordSerializer(new VerdictSerializationSchema(config.verdictsTopic()))
				.build();
	}
}
