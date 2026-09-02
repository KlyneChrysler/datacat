package io.datacat.classifier.operators;

import io.datacat.classifier.model.FeatureMath;
import io.datacat.classifier.model.SessionFeatures;
import org.apache.flink.streaming.api.functions.windowing.ProcessWindowFunction;
import org.apache.flink.streaming.api.windowing.windows.TimeWindow;
import org.apache.flink.util.Collector;

/** Attaches key and window bounds, then turns raw samples into features. */
public final class FeatureWindowFunction
		extends ProcessWindowFunction<FeatureAccumulator, SessionFeatures, String, TimeWindow> {

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
		return new SessionFeatures(
				sessionId,
				windowEnd,
				acc.requestCount,
				FeatureMath.requestsPerMinute(acc.requestCount, windowSeconds),
				FeatureMath.intervalCv(acc.timestamps),
				FeatureMath.normalizedPathEntropy(acc.paths),
				FeatureMath.distinctPaths(acc.paths),
				acc.userAgent,
				acc.requestCount == 0 ? 0 : (double) acc.verifiedCount / acc.requestCount);
	}
}
