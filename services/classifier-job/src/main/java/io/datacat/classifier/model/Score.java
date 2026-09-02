package io.datacat.classifier.model;

/** One signal judgment, zero human to one automated. */
public record Score(String signal, double botLikelihood) {

	public Score {
		if (botLikelihood < 0 || botLikelihood > 1) {
			throw new IllegalArgumentException("bot likelihood must be within [0,1], got " + botLikelihood + " from " + signal);
		}
	}
}
