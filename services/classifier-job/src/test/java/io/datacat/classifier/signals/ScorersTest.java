package io.datacat.classifier.signals;

import io.datacat.classifier.model.SessionFeatures;
import org.junit.jupiter.api.Test;

import static org.junit.jupiter.api.Assertions.assertEquals;
import static org.junit.jupiter.api.Assertions.assertTrue;

class ScorersTest {

	private static SessionFeatures features(double intervalCv, double rpm, double entropy,
			long count, String userAgent) {
		return new SessionFeatures("s-1", 0L, count, rpm, intervalCv, entropy, count, userAgent);
	}

	@Test
	void metronomeTimingReadsAutomated() {
		double score = new TimingRegularityScorer()
				.score(features(0.02, 10, 0.5, 20, "x")).botLikelihood();

		assertTrue(score > 0.9, "cv near zero should read automated, got " + score);
	}

	@Test
	void jitteryTimingReadsHuman() {
		double score = new TimingRegularityScorer()
				.score(features(1.4, 10, 0.5, 20, "x")).botLikelihood();

		assertEquals(0.0, score, 1e-9);
	}

	@Test
	void insufficientTimingEvidenceIsNeutral() {
		double score = new TimingRegularityScorer()
				.score(features(Double.NaN, 10, 0.5, 2, "x")).botLikelihood();

		assertEquals(0.5, score, 1e-9);
	}

	@Test
	void extremeRateSaturates() {
		double score = new RequestRateScorer()
				.score(features(1, 500, 0.5, 500, "x")).botLikelihood();

		assertEquals(1.0, score, 1e-9);
	}

	@Test
	void entropyIsDampedOnSmallSamples() {
		double small = new PathEntropyScorer().score(features(1, 1, 1.0, 4, "x")).botLikelihood();
		double large = new PathEntropyScorer().score(features(1, 1, 1.0, 40, "x")).botLikelihood();

		assertTrue(small < large, "few requests must damp entropy evidence");
		assertEquals(1.0, large, 1e-9);
	}

	@Test
	void declaredAutomationUserAgentScoresHigh() {
		double curl = new UserAgentScorer().score(features(1, 1, 0, 5, "curl/8.7.1")).botLikelihood();
		double browser = new UserAgentScorer()
				.score(features(1, 1, 0, 5, "Mozilla/5.0 (Macintosh) Safari/605.1")).botLikelihood();

		assertTrue(curl > 0.9);
		assertTrue(browser < 0.2);
	}

	@Test
	void registryContainsEverySignalOnce() {
		long distinctNames = Scorers.all().stream().map(SignalScorer::name).distinct().count();

		assertEquals(Scorers.all().size(), distinctNames);
		assertEquals(4, Scorers.all().size());
	}
}
