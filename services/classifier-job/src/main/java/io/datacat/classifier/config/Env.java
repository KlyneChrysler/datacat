package io.datacat.classifier.config;

/**
 * Generic environment readers (utility kind, file taxonomy rule 2) —
 * separate from the JobConfig shape they populate.
 */
public final class Env {

	private Env() {
	}

	public static String string(String name, String fallback) {
		String raw = System.getenv(name);
		return raw == null || raw.isBlank() ? fallback : raw;
	}

	public static long longValue(String name, long fallback) {
		String raw = System.getenv(name);
		return raw == null || raw.isBlank() ? fallback : Long.parseLong(raw);
	}
}
