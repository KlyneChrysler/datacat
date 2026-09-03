package io.datacat.policy.domain;

/** Raised when a rule pairs a class with a forbidden action. */
public final class IncompatibleRuleException extends RuntimeException {

	public IncompatibleRuleException(Classification appliesTo, ActionType action) {
		super("action " + action + " is not allowed for " + appliesTo);
	}
}
