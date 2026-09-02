package io.datacat.classifier.model;

/** Mirror of the Go wire schema for one observed request. */
public record RequestEvent(String sessionId, long timestampMillis, String method, String path, String clientIp, String userAgent, String headerOrder, String tlsFingerprint, boolean verifiedAgent) {
}
