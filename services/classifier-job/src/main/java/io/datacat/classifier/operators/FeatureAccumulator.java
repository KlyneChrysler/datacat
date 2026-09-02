package io.datacat.classifier.operators;

import io.datacat.classifier.model.RequestEvent;

import java.util.ArrayList;

/**
 * Mutable window accumulator (Flink POJO: public fields, no-arg constructor).
 * Samples are capped so a hostile session cannot grow window state without
 * bound; features degrade gracefully past the cap because the statistics
 * already saturate long before it.
 */
public final class FeatureAccumulator {

	public static final int MAX_SAMPLES = 512;

	public ArrayList<Long> timestamps = new ArrayList<>();
	public ArrayList<String> paths = new ArrayList<>();
	public String userAgent = "";
	public long requestCount;
	public long verifiedCount;

	public FeatureAccumulator add(RequestEvent event) {
		requestCount++;
		if (event.verifiedAgent()) {
			verifiedCount++;
		}
		userAgent = event.userAgent();
		if (timestamps.size() < MAX_SAMPLES) {
			timestamps.add(event.timestampMillis());
			paths.add(event.path());
		}
		return this;
	}

	public FeatureAccumulator merge(FeatureAccumulator other) {
		requestCount += other.requestCount;
		verifiedCount += other.verifiedCount;
		if (!other.userAgent.isEmpty()) {
			userAgent = other.userAgent;
		}
		other.timestamps.stream().limit(MAX_SAMPLES - (long) timestamps.size()).forEach(timestamps::add);
		other.paths.stream().limit(MAX_SAMPLES - (long) paths.size()).forEach(paths::add);
		return this;
	}
}
