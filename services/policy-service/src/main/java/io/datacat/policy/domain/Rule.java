package io.datacat.policy.domain;

import java.time.Instant;
import java.util.UUID;

/** One immutable enforcement rule with its invariants enforced at birth. */
public record Rule(UUID id, String name, Classification appliesTo, ActionType action, int rateLimitPerMinute, Instant createdAt) {

	public static Rule from(CreateRuleCommand command) {
		requireCompatible(command.appliesTo(), command.action());

		return new Rule(UUID.randomUUID(), command.name(), command.appliesTo(), command.action(), command.rateLimitPerMinute(), Instant.now());
	}

	private static void requireCompatible(Classification appliesTo, ActionType action) {
		if (appliesTo == Classification.HUMAN && action == ActionType.BLOCK) {
			throw new IncompatibleRuleException(appliesTo, action);
		}
	}
}
