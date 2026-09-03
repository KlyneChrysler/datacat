package io.datacat.policy.domain;

/** The actions a rule can apply. */
public enum ActionType {
	ALLOW,
	RATE_LIMIT,
	CHALLENGE,
	BLOCK
}
