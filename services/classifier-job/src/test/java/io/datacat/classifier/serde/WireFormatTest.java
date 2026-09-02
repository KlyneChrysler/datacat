package io.datacat.classifier.serde;

import com.fasterxml.jackson.databind.JsonNode;
import com.fasterxml.jackson.databind.ObjectMapper;
import io.datacat.classifier.model.RequestEvent;
import io.datacat.classifier.model.Verdict;
import org.junit.jupiter.api.Test;

import java.nio.charset.StandardCharsets;

import static org.junit.jupiter.api.Assertions.assertEquals;
import static org.junit.jupiter.api.Assertions.assertNull;

/** Guards compatibility with the Go services' wire format (pkg/events). */
class WireFormatTest {

	// Verbatim event produced by edge-proxy during the Kafka-phase e2e run.
	private static final String GO_PRODUCED_EVENT = """
			{"session_id":"demo-scraper-1","timestamp":"2026-09-02T07:37:28.142807Z",\
			"method":"GET","path":"/products","client_ip":"::1","user_agent":"curl/8.7.1",\
			"header_order":"c29cbedeb1df3065","tls_fingerprint":""}""";

	@Test
	void deserializesEventExactlyAsGoProducesIt() {
		RequestEvent event = new RequestEventDeserializer()
				.deserialize(GO_PRODUCED_EVENT.getBytes(StandardCharsets.UTF_8));

		assertEquals("demo-scraper-1", event.sessionId());
		assertEquals("GET", event.method());
		assertEquals("/products", event.path());
		assertEquals("curl/8.7.1", event.userAgent());
		assertEquals(1788334648142L, event.timestampMillis());
	}

	@Test
	void malformedEventIsSkippedNotFatal() {
		assertNull(new RequestEventDeserializer().deserialize("not json".getBytes(StandardCharsets.UTF_8)));
	}

	@Test
	void serializesVerdictInGoWireShape() throws Exception {
		Verdict verdict = new Verdict("s-9", 1788334658390L, "abusive", 0.97);

		JsonNode node = new ObjectMapper().readTree(VerdictSerializationSchema.toJson(verdict));

		assertEquals("s-9", node.get("session_id").asText());
		assertEquals("abusive", node.get("classification").asText());
		assertEquals(0.97, node.get("confidence").asDouble(), 1e-9);
		assertEquals("2026-09-02T07:37:38.390Z", node.get("timestamp").asText());
	}
}
