package io.datacat.classifier.signals;

import io.datacat.classifier.model.Score;
import io.datacat.classifier.model.SessionFeatures;

/** Never revisiting paths reads as crawling, damped on small samples. */
public final class PathEntropyScorer implements SignalScorer {

	private static final long FULL_WEIGHT_REQUESTS = 20;

	@Override
	public String name() {
		return "path_entropy";
	}

	@Override
	public Score score(SessionFeatures features) {
		double volumeWeight = Math.min((double) features.requestCount() / FULL_WEIGHT_REQUESTS, 1.0);

		return new Score(name(), features.pathEntropy() * volumeWeight);
	}
}
