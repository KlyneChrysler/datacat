package io.datacat.classifier.signals;

import java.util.List;

/** The registry: the single line a new signal touches. */
public final class Scorers {

	private Scorers() {
	}

	public static List<SignalScorer> all() {
		return List.of(
				new TimingRegularityScorer(),
				new RequestRateScorer(),
				new PathEntropyScorer(),
				new UserAgentScorer());
	}
}
