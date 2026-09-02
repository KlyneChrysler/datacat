package io.datacat.classifier.signals;

import io.datacat.classifier.model.Score;
import io.datacat.classifier.model.SessionFeatures;

/**
 * Crawlers sweep the URL space and rarely revisit; humans loop over a few
 * pages. High normalized path entropy over enough requests reads bot-like.
 * Entropy alone is weak evidence on small samples, so it is damped by volume.
 */
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
