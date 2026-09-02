package io.datacat.classifier.model;

/**
 * Mirror of the Go wire schema in pkg/events (request_event.go). Field
 * mapping to/from snake_case JSON lives in the serde package only.
 */
public record RequestEvent(String sessionId, long timestampMillis, String method, String path,
		String clientIp, String userAgent, String headerOrder, String tlsFingerprint) {
}
