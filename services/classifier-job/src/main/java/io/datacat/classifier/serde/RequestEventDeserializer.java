package io.datacat.classifier.serde;

import com.fasterxml.jackson.databind.JsonNode;
import com.fasterxml.jackson.databind.ObjectMapper;
import io.datacat.classifier.model.RequestEvent;
import org.apache.flink.api.common.serialization.DeserializationSchema;
import org.apache.flink.api.common.typeinfo.TypeInformation;
import org.slf4j.Logger;
import org.slf4j.LoggerFactory;

import java.time.Instant;

/**
 * Decodes the Go wire format (pkg/events RequestEvent, snake_case JSON).
 * A malformed record is logged and skipped (returns null) — one poison
 * message must never crash the job.
 */
public final class RequestEventDeserializer implements DeserializationSchema<RequestEvent> {

	private static final Logger LOG = LoggerFactory.getLogger(RequestEventDeserializer.class);
	private static final ObjectMapper MAPPER = new ObjectMapper();

	@Override
	public RequestEvent deserialize(byte[] message) {
		try {
			return fromJson(MAPPER.readTree(message));
		} catch (Exception e) {
			LOG.error("skipping malformed request event", e);
			return null;
		}
	}

	private static RequestEvent fromJson(JsonNode node) {
		return new RequestEvent(
				node.path("session_id").asText(),
				Instant.parse(node.path("timestamp").asText()).toEpochMilli(),
				node.path("method").asText(),
				node.path("path").asText(),
				node.path("client_ip").asText(),
				node.path("user_agent").asText(),
				node.path("header_order").asText(),
				node.path("tls_fingerprint").asText(),
				node.path("verified_agent").asBoolean(false));
	}

	@Override
	public boolean isEndOfStream(RequestEvent nextElement) {
		return false;
	}

	@Override
	public TypeInformation<RequestEvent> getProducedType() {
		return TypeInformation.of(RequestEvent.class);
	}
}
