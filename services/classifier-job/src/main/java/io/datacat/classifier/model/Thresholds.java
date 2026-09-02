package io.datacat.classifier.model;

import java.io.Serializable;
import java.util.List;

/** Turns averaged scores and identity into a classification. */
public record Thresholds(double humanBelow, double abusiveFrom) implements Serializable {

	public static final String HUMAN = "human";
	public static final String VERIFIED = "verified_agent";
	public static final String UNVERIFIED = "unverified_automation";
	public static final String ABUSIVE = "abusive";

	// A session is identity verified when this share of requests is signed.
	private static final double VERIFIED_SHARE_FLOOR = 0.9;

	public Thresholds {
		if (humanBelow <= 0 || abusiveFrom <= humanBelow || abusiveFrom > 1) {
			throw new IllegalArgumentException("thresholds must satisfy 0 < humanBelow < abusiveFrom <= 1");
		}
	}

	public static Thresholds defaults() {
		return new Thresholds(0.45, 0.75);
	}

	// Verified identity wins unless behavior crosses the abusive line.
	public Verdict verdictFor(String sessionId, long atMillis, List<Score> scores, double verifiedShare) {
		double average = average(scores);

		if (verifiedShare >= VERIFIED_SHARE_FLOOR && average < abusiveFrom) {
			return new Verdict(sessionId, atMillis, VERIFIED, verifiedShare);
		}

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

	// Confidence is the distance from the nearest boundary, normalized.
	private double confidence(double average) {
		double toBoundary = Math.min(Math.abs(average - humanBelow), Math.abs(average - abusiveFrom));
		double widest = Math.max(humanBelow, Math.max(abusiveFrom - humanBelow, 1 - abusiveFrom));

		return Math.min(toBoundary / widest, 1.0);
	}

	private static double average(List<Score> scores) {
		return scores.stream().mapToDouble(Score::botLikelihood).average().orElse(0.5);
	}
}
