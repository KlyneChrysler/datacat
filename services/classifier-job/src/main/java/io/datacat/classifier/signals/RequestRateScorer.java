package io.datacat.classifier.signals;

import io.datacat.classifier.model.Score;
import io.datacat.classifier.model.SessionFeatures;

/** Sustained request rates far above human browsing speed indicate automation. */
public final class RequestRateScorer implements SignalScorer {

	private static final double CERTAINLY_AUTOMATED_RPM = 120.0;

	@Override
	public String name() {
		return "request_rate";
	}

	@Override
	public Score score(SessionFeatures features) {
		double likelihood = Math.min(features.requestsPerMinute() / CERTAINLY_AUTOMATED_RPM, 1.0);
		return new Score(name(), likelihood);
	}
}
