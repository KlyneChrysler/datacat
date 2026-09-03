package io.datacat.policy.infrastructure;

import io.datacat.policy.domain.ActionType;
import io.datacat.policy.domain.Classification;
import io.datacat.policy.domain.Rule;
import org.springframework.stereotype.Component;

/** Domain to storage conversion. */
@Component
public class RuleEntityMapper {

	public RuleEntity toEntity(Rule rule) {
		return new RuleEntity(rule.id(), rule.name(), rule.appliesTo().name(), rule.action().name(), rule.rateLimitPerMinute(), rule.createdAt());
	}

	public Rule toDomain(RuleEntity entity) {
		return new Rule(entity.getId(), entity.getName(), Classification.valueOf(entity.getAppliesTo()), ActionType.valueOf(entity.getAction()), entity.getRateLimitPerMinute(), entity.getCreatedAt());
	}
}
