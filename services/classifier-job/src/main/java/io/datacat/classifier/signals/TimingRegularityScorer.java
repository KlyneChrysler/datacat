package io.datacat.classifier.signals;

import io.datacat.classifier.model.Score;
import io.datacat.classifier.model.SessionFeatures;

/**
 * Humans are jittery; schedulers are metronomes. A low coefficient of
 * variation across inter-request intervals is strong automation evidence.
 */
public final class TimingRegularityScorer implements SignalScorer {

	private static final double HUMAN_JITTER_CV = 1.0;
	private static final double NEUTRAL = 0.5;

	@Override
	public String name() {
		return "timing_regularity";
	}

	@Override
	public Score score(SessionFeatures features) {
		double cv = features.intervalCv();
		if (Double.isNaN(cv)) {
			return new Score(name(), NEUTRAL); // too few requests: no evidence
		}
		double likelihood = 1.0 - Math.min(cv / HUMAN_JITTER_CV, 1.0);
		return new Score(name(), likelihood);
	}
}
