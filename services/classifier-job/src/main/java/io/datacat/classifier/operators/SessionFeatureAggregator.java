package io.datacat.classifier.operators;

import io.datacat.classifier.model.RequestEvent;
import org.apache.flink.api.common.functions.AggregateFunction;

/**
 * Incremental per-window accumulation: Flink folds each event into the
 * accumulator as it arrives, so window state stays small. Feature assembly
 * (which needs the key and window bounds) happens in FeatureWindowFunction.
 */
public final class SessionFeatureAggregator
		implements AggregateFunction<RequestEvent, FeatureAccumulator, FeatureAccumulator> {

	@Override
	public FeatureAccumulator createAccumulator() {
		return new FeatureAccumulator();
	}

	@Override
	public FeatureAccumulator add(RequestEvent event, FeatureAccumulator acc) {
		return acc.add(event);
	}

	@Override
	public FeatureAccumulator getResult(FeatureAccumulator acc) {
		return acc;
	}

	@Override
	public FeatureAccumulator merge(FeatureAccumulator a, FeatureAccumulator b) {
		return a.merge(b);
	}
}
