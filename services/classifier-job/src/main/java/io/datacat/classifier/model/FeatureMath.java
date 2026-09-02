package io.datacat.classifier.model;

import java.util.ArrayList;
import java.util.HashMap;
import java.util.List;
import java.util.Map;

/** Pure helpers turning raw window samples into features. No state, no I/O. */
public final class FeatureMath {

	private FeatureMath() {
	}

	/**
	 * Coefficient of variation of inter-request intervals. Low values mean
	 * metronome-regular traffic. Returns NaN when fewer than 3 requests make
	 * the statistic meaningless — scorers must treat NaN as "no evidence".
	 */
	public static double intervalCv(List<Long> timestampsMillis) {
		List<Long> intervals = sortedIntervals(timestampsMillis);
		if (intervals.size() < 2) {
			return Double.NaN;
		}
		double mean = intervals.stream().mapToLong(Long::longValue).average().orElse(0);
		if (mean == 0) {
			return 0;
		}
		return Math.sqrt(variance(intervals, mean)) / mean;
	}

	/**
	 * Shannon entropy of the path distribution, normalized to [0,1] by the
	 * maximum entropy for the sample count. 1 = never revisits a path
	 * (crawler-like); 0 = hammers a single path.
	 */
	public static double normalizedPathEntropy(List<String> paths) {
		if (paths.size() < 2) {
			return 0;
		}
		double entropy = 0;
		for (long count : countByPath(paths).values()) {
			double p = (double) count / paths.size();
			entropy -= p * (Math.log(p) / Math.log(2));
		}
		double max = Math.log(paths.size()) / Math.log(2);
		return max == 0 ? 0 : Math.min(entropy / max, 1.0);
	}

	public static long distinctPaths(List<String> paths) {
		return paths.stream().distinct().count();
	}

	public static double requestsPerMinute(long count, long windowSeconds) {
		return count * 60.0 / windowSeconds;
	}

	private static List<Long> sortedIntervals(List<Long> timestampsMillis) {
		List<Long> sorted = new ArrayList<>(timestampsMillis);
		sorted.sort(Long::compareTo);
		List<Long> intervals = new ArrayList<>(Math.max(sorted.size() - 1, 0));
		for (int i = 1; i < sorted.size(); i++) {
			intervals.add(sorted.get(i) - sorted.get(i - 1));
		}
		return intervals;
	}

	private static double variance(List<Long> values, double mean) {
		double sum = 0;
		for (long v : values) {
			sum += (v - mean) * (v - mean);
		}
		return sum / values.size();
	}

	private static Map<String, Long> countByPath(List<String> paths) {
		Map<String, Long> counts = new HashMap<>();
		for (String path : paths) {
			counts.merge(path, 1L, Long::sum);
		}
		return counts;
	}
}
