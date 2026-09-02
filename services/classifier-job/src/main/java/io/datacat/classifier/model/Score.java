package io.datacat.classifier.model;

/** One signal's judgment: 0 = certainly human, 1 = certainly automated. */
public record Score(String signal, double botLikelihood) {

	public Score {
		if (botLikelihood < 0 || botLikelihood > 1) {
			throw new IllegalArgumentException(
					"bot likelihood must be within [0,1], got " + botLikelihood + " from " + signal);
		}
	}
}
