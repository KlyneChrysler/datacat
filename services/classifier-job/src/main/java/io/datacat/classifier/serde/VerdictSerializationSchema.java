package io.datacat.classifier.serde;

import com.fasterxml.jackson.databind.ObjectMapper;
import com.fasterxml.jackson.databind.node.ObjectNode;
import io.datacat.classifier.model.Verdict;
import org.apache.flink.connector.kafka.sink.KafkaRecordSerializationSchema;
import org.apache.kafka.clients.producer.ProducerRecord;

import java.nio.charset.StandardCharsets;
import java.time.Instant;

/**
 * Encodes verdicts in the Go wire format (pkg/events Verdict, snake_case,
 * RFC3339 timestamp), keyed by session ID for per-session ordering.
 */
public final class VerdictSerializationSchema implements KafkaRecordSerializationSchema<Verdict> {

	private static final ObjectMapper MAPPER = new ObjectMapper();

	private final String topic;

	public VerdictSerializationSchema(String topic) {
		this.topic = topic;
	}

	@Override
	public ProducerRecord<byte[], byte[]> serialize(Verdict verdict, KafkaSinkContext context, Long timestamp) {
		byte[] key = verdict.sessionId().getBytes(StandardCharsets.UTF_8);
		return new ProducerRecord<>(topic, key, toJson(verdict));
	}

	static byte[] toJson(Verdict verdict) {
		ObjectNode node = MAPPER.createObjectNode();
		node.put("session_id", verdict.sessionId());
		node.put("timestamp", Instant.ofEpochMilli(verdict.timestampMillis()).toString());
		node.put("classification", verdict.classification());
		node.put("confidence", verdict.confidence());
		return node.toString().getBytes(StandardCharsets.UTF_8);
	}
}
