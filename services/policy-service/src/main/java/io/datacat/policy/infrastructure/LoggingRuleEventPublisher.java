package io.datacat.policy.infrastructure;

import io.datacat.policy.application.RuleEventPublisher;
import io.datacat.policy.domain.Rule;
import org.slf4j.Logger;
import org.slf4j.LoggerFactory;
import org.springframework.stereotype.Component;

/** Logs rule events, the Kafka publisher replaces this later. */
@Component
public class LoggingRuleEventPublisher implements RuleEventPublisher {

	private static final Logger log = LoggerFactory.getLogger(LoggingRuleEventPublisher.class);

	@Override
	public void publishCreated(Rule rule) {
		log.info("rule created id={} name={} appliesTo={} action={}", rule.id(), rule.name(), rule.appliesTo(), rule.action());
	}
}
