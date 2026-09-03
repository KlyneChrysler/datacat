package io.datacat.policy.application;

import io.datacat.policy.domain.Rule;

/** What the service needs to announce rule changes. */
public interface RuleEventPublisher {

	void publishCreated(Rule rule);
}
