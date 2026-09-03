package io.datacat.policy.domain;

import java.util.UUID;

/** Raised when a rule id does not exist. */
public final class RuleNotFoundException extends RuntimeException {

	public RuleNotFoundException(UUID id) {
		super("rule not found: " + id);
	}
}
