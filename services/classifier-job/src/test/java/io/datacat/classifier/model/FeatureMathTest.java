package io.datacat.classifier.model;

import org.junit.jupiter.api.Test;

import java.util.List;

import static org.junit.jupiter.api.Assertions.assertEquals;
import static org.junit.jupiter.api.Assertions.assertTrue;

class FeatureMathTest {

	@Test
	void regularIntervalsHaveNearZeroCv() {
		List<Long> metronome = List.of(0L, 500L, 1000L, 1500L, 2000L);

		assertEquals(0.0, FeatureMath.intervalCv(metronome), 1e-9);
	}

	@Test
	void jitteryIntervalsHaveHighCv() {
		List<Long> human = List.of(0L, 200L, 4200L, 4500L, 30000L);

		assertTrue(FeatureMath.intervalCv(human) > 0.8, "human jitter should score high cv");
	}

	@Test
	void tooFewSamplesYieldNaN() {
		assertTrue(Double.isNaN(FeatureMath.intervalCv(List.of(0L, 100L))));
		assertTrue(Double.isNaN(FeatureMath.intervalCv(List.of())));
	}

	@Test
	void unsortedTimestampsAreHandled() {
		List<Long> outOfOrder = List.of(1000L, 0L, 2000L, 500L, 1500L);

		assertEquals(0.0, FeatureMath.intervalCv(outOfOrder), 1e-9);
	}

	@Test
	void neverRevisitingPathsMaximizesEntropy() {
		List<String> crawl = List.of("/a", "/b", "/c", "/d");

		assertEquals(1.0, FeatureMath.normalizedPathEntropy(crawl), 1e-9);
	}

	@Test
	void singlePathHasZeroEntropy() {
		List<String> hammer = List.of("/login", "/login", "/login");

		assertEquals(0.0, FeatureMath.normalizedPathEntropy(hammer), 1e-9);
	}

	@Test
	void requestsPerMinuteScalesByWindow() {
		assertEquals(60.0, FeatureMath.requestsPerMinute(60, 60), 1e-9);
		assertEquals(12.0, FeatureMath.requestsPerMinute(60, 300), 1e-9);
	}
}
