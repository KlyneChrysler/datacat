package io.datacat.classifier.signals;

import io.datacat.classifier.model.Score;
import io.datacat.classifier.model.SessionFeatures;

import java.util.List;
import java.util.Locale;

/** A declared automation user agent is strong but spoofable evidence. */
public final class UserAgentScorer implements SignalScorer {

	private static final List<String> AUTOMATION_MARKERS = List.of("curl", "wget", "python", "go-http-client", "bot", "crawler", "spider", "scrapy", "headless", "phantom", "httpclient");

	private static final double DECLARED_AUTOMATION = 0.95;
	private static final double MISSING_UA = 0.7;
	private static final double BROWSER_LIKE = 0.15;

	@Override
	public String name() {
		return "user_agent";
	}

	@Override
	public Score score(SessionFeatures features) {
		return new Score(name(), likelihood(features.userAgent()));
	}

	private static double likelihood(String userAgent) {
		if (userAgent == null || userAgent.isBlank()) {
			return MISSING_UA;
		}

		String lowered = userAgent.toLowerCase(Locale.ROOT);
		boolean declared = AUTOMATION_MARKERS.stream().anyMatch(lowered::contains);

		return declared ? DECLARED_AUTOMATION : BROWSER_LIKE;
	}
}
