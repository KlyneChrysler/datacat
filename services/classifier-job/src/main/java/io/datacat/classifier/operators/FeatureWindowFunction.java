package io.datacat.classifier.operators;

import io.datacat.classifier.model.FeatureMath;
import io.datacat.classifier.model.SessionFeatures;
import org.apache.flink.streaming.api.functions.windowing.ProcessWindowFunction;
import org.apache.flink.streaming.api.windowing.windows.TimeWindow;
import org.apache.flink.util.Collector;

/** Turns one window of raw samples into session features. */
public final class FeatureWindowFunction extends ProcessWindowFunction<FeatureAccumulator, SessionFeatures, String, TimeWindow> {

	private final long windowSeconds;

	public FeatureWindowFunction(long windowSeconds) {
		this.windowSeconds = windowSeconds;
	}

	@Override
	public void process(String sessionId, Context context, Iterable<FeatureAccumulator> accumulators, Collector<SessionFeatures> out) {
		FeatureAccumulator acc = accumulators.iterator().next();

		out.collect(toFeatures(sessionId, context.window().getEnd(), acc));
	}

	private SessionFeatures toFeatures(String sessionId, long windowEnd, FeatureAccumulator acc) {
		double perMinute = FeatureMath.requestsPerMinute(acc.requestCount, windowSeconds);
		double intervalCv = FeatureMath.intervalCv(acc.timestamps);
		double entropy = FeatureMath.normalizedPathEntropy(acc.paths);
		long distinct = FeatureMath.distinctPaths(acc.paths);
		double verifiedShare = acc.requestCount == 0 ? 0 : (double) acc.verifiedCount / acc.requestCount;

		return new SessionFeatures(sessionId, windowEnd, acc.requestCount, perMinute, intervalCv, entropy, distinct, acc.userAgent, verifiedShare);
	}
}
