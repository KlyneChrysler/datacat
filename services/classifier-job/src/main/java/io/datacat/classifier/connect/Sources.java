package io.datacat.classifier.connect;

import io.datacat.classifier.config.JobConfig;
import io.datacat.classifier.model.RequestEvent;
import io.datacat.classifier.serde.RequestEventDeserializer;
import org.apache.flink.api.common.eventtime.WatermarkStrategy;
import org.apache.flink.connector.kafka.source.KafkaSource;
import org.apache.flink.connector.kafka.source.enumerator.initializer.OffsetsInitializer;
import org.apache.flink.streaming.api.datastream.DataStream;
import org.apache.flink.streaming.api.environment.StreamExecutionEnvironment;

import java.time.Duration;

/** Inbound edges of the job graph. Wire-level choices live here only. */
public final class Sources {

	private static final Duration OUT_OF_ORDERNESS = Duration.ofSeconds(10);
	private static final Duration IDLENESS = Duration.ofSeconds(30);

	private Sources() {
	}

	public static DataStream<RequestEvent> requestEvents(StreamExecutionEnvironment env, JobConfig config) {
		KafkaSource<RequestEvent> source = KafkaSource.<RequestEvent>builder()
				.setBootstrapServers(config.kafkaBrokers())
				.setTopics(config.requestsTopic())
				.setGroupId(config.consumerGroup())
				.setStartingOffsets(OffsetsInitializer.earliest())
				.setValueOnlyDeserializer(new RequestEventDeserializer())
				.build();
		return env.fromSource(source, watermarks(), "request-events");
	}

	// Event time with bounded out-of-orderness; idleness keeps quiet
	// partitions from stalling the watermark.
	private static WatermarkStrategy<RequestEvent> watermarks() {
		return WatermarkStrategy
				.<RequestEvent>forBoundedOutOfOrderness(OUT_OF_ORDERNESS)
				.withTimestampAssigner((event, recordTs) -> event.timestampMillis())
				.withIdleness(IDLENESS);
	}
}
