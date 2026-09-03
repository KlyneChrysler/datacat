package io.datacat.policy.api;

import io.datacat.policy.domain.Rule;

import java.time.Instant;
import java.util.UUID;

/** Wire shape of one rule. */
public record RuleResponse(UUID id, String name, String appliesTo, String action, int rateLimitPerMinute, Instant createdAt) {

	public static RuleResponse from(Rule rule) {
		return new RuleResponse(rule.id(), rule.name(), rule.appliesTo().name(), rule.action().name(), rule.rateLimitPerMinute(), rule.createdAt());
	}
}
