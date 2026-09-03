package io.datacat.policy.domain;

/** The session classes a rule can target. */
public enum Classification {
	HUMAN,
	VERIFIED_AGENT,
	UNVERIFIED_AUTOMATION,
	ABUSIVE
}
