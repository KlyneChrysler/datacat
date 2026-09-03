package io.datacat.policy.application;

import io.datacat.policy.domain.CreateRuleCommand;
import io.datacat.policy.domain.Rule;
import io.datacat.policy.domain.RuleNotFoundException;
import org.springframework.stereotype.Service;
import org.springframework.transaction.annotation.Transactional;

import java.util.List;
import java.util.UUID;

/** Rule use cases, every public method delegates. */
@Service
public class RuleService {

	private final RuleRepository rules;
	private final RuleEventPublisher events;

	public RuleService(RuleRepository rules, RuleEventPublisher events) {
		this.rules = rules;
		this.events = events;
	}

	@Transactional
	public Rule create(CreateRuleCommand command) {
		Rule rule = Rule.from(command);
		Rule saved = rules.save(rule);

		events.publishCreated(saved);

		return saved;
	}

	@Transactional(readOnly = true)
	public Rule get(UUID id) {
		return rules.findById(id).orElseThrow(() -> new RuleNotFoundException(id));
	}

	@Transactional(readOnly = true)
	public List<Rule> list() {
		return rules.findAll();
	}
}
