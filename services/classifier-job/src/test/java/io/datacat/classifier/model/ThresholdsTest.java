package io.datacat.classifier.model;

import org.junit.jupiter.api.Test;

import java.util.List;

import static org.junit.jupiter.api.Assertions.assertEquals;
import static org.junit.jupiter.api.Assertions.assertThrows;
import static org.junit.jupiter.api.Assertions.assertTrue;

class ThresholdsTest {

	private static Verdict verdictWith(double... likelihoods) {
		return verdictWithShare(0, likelihoods);
	}

	private static Verdict verdictWithShare(double verifiedShare, double... likelihoods) {
		List<Score> scores = java.util.Arrays.stream(likelihoods)
				.mapToObj(l -> new Score("test", l))
				.toList();
		return Thresholds.defaults().verdictFor("s-1", 42L, scores, verifiedShare);
	}

	@Test
	void lowAverageIsHuman() {
		Verdict verdict = verdictWith(0.1, 0.2);

		assertEquals(Thresholds.HUMAN, verdict.classification());
	}

	@Test
	void midAverageIsUnverifiedAutomation() {
		Verdict verdict = verdictWith(0.5, 0.6);

		assertEquals(Thresholds.UNVERIFIED, verdict.classification());
	}

	@Test
	void highAverageIsAbusive() {
		Verdict verdict = verdictWith(0.9, 0.95);

		assertEquals(Thresholds.ABUSIVE, verdict.classification());
	}

	@Test
	void confidenceGrowsAwayFromBoundaries() {
		double onBoundary = verdictWith(0.45).confidence();
		double deepAbusive = verdictWith(1.0).confidence();

		assertEquals(0.0, onBoundary, 1e-9);
		assertTrue(deepAbusive > 0.4, "clear-cut case should be confident");
	}

	@Test
	void invalidThresholdOrderingIsRejected() {
		assertThrows(IllegalArgumentException.class, () -> new Thresholds(0.8, 0.5));
	}

	@Test
	void verifiedIdentityYieldsVerifiedAgent() {
		Verdict verdict = verdictWithShare(0.95, 0.5, 0.6); // bot-ish behavior, signed

		assertEquals(Thresholds.VERIFIED, verdict.classification());
	}

	@Test
	void abusiveBehaviorTrumpsVerifiedIdentity() {
		Verdict verdict = verdictWithShare(1.0, 0.9, 0.95);

		assertEquals(Thresholds.ABUSIVE, verdict.classification());
	}

	@Test
	void partiallySignedSessionIsNotVerified() {
		Verdict verdict = verdictWithShare(0.5, 0.5, 0.6);

		assertEquals(Thresholds.UNVERIFIED, verdict.classification());
	}
}
