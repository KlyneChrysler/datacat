package io.datacat.classifier.model;

/** Mirror of the Go wire schema for one classification. */
public record Verdict(String sessionId, long timestampMillis, String classification, double confidence) {
}
