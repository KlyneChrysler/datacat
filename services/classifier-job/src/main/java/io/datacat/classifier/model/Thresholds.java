package io.datacat.classifier.model;

import java.io.Serializable;
import java.util.List;

/**
 * Turns averaged signal scores into a classification. Wire strings match the
 * Go domain (services/enforcement/internal/domain/verdict.go).
 */
public record Thresholds(double humanBelow, double abusiveFrom) implements Serializable {

	public static final String HUMAN = "human";
	public static final String UNVERIFIED = "unverified_automation";
	public static final String ABUSIVE = "abusive";

	public Thresholds {
		if (humanBelow <= 0 || abusiveFrom <= humanBelow || abusiveFrom > 1) {
			throw new IllegalArgumentException(
					"thresholds must satisfy 0 < humanBelow < abusiveFrom <= 1");
		}
	}

	public static Thresholds defaults() {
		return new Thresholds(0.45, 0.75);
	}

	public Verdict verdictFor(String sessionId, long atMillis, List<Score> scores) {
		double average = average(scores);
		return new Verdict(sessionId, atMillis, classify(average), confidence(average));
	}

	private String classify(double average) {
		if (average < humanBelow) {
			return HUMAN;
		}
		if (average < abusiveFrom) {
			return UNVERIFIED;
		}
		return ABUSIVE;
	}

	// Confidence is how far the average sits from the nearest boundary,
	// normalized to [0,1]; a score on a threshold is maximally uncertain.
	private double confidence(double average) {
		double toBoundary = Math.min(Math.abs(average - humanBelow), Math.abs(average - abusiveFrom));
		double widest = Math.max(humanBelow, Math.max(abusiveFrom - humanBelow, 1 - abusiveFrom));
		return Math.min(toBoundary / widest, 1.0);
	}

	private static double average(List<Score> scores) {
		return scores.stream().mapToDouble(Score::botLikelihood).average().orElse(0.5);
	}
}
