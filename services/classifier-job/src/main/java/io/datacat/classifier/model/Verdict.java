package io.datacat.classifier.model;

/** Mirror of the Go wire schema in pkg/events (Verdict). */
public record Verdict(String sessionId, long timestampMillis, String classification,
		double confidence) {
}
