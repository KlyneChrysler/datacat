package io.datacat.classifier.model;

/** Facts about one session over one window, judgment lives in signals. */
public record SessionFeatures(String sessionId, long windowEndMillis, long requestCount, double requestsPerMinute, double intervalCv, double pathEntropy, long distinctPaths, String userAgent, double verifiedShare) {
}
