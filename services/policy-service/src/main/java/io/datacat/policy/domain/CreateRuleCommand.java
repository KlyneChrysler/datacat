package io.datacat.policy.domain;

/** The validated intent to create one rule. */
public record CreateRuleCommand(String name, Classification appliesTo, ActionType action, int rateLimitPerMinute) {
}
