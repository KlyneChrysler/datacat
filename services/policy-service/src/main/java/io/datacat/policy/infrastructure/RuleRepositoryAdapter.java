package io.datacat.policy.infrastructure;

import io.datacat.policy.application.RuleRepository;
import io.datacat.policy.domain.Rule;
import org.springframework.stereotype.Component;

import java.util.List;
import java.util.Optional;
import java.util.UUID;

/** JPA backed implementation of the rule repository. */
@Component
public class RuleRepositoryAdapter implements RuleRepository {

	private final RuleJpaRepository jpa;
	private final RuleEntityMapper mapper;

	public RuleRepositoryAdapter(RuleJpaRepository jpa, RuleEntityMapper mapper) {
		this.jpa = jpa;
		this.mapper = mapper;
	}

	@Override
	public Rule save(Rule rule) {
		RuleEntity saved = jpa.save(mapper.toEntity(rule));

		return mapper.toDomain(saved);
	}

	@Override
	public Optional<Rule> findById(UUID id) {
		return jpa.findById(id).map(mapper::toDomain);
	}

	@Override
	public List<Rule> findAll() {
		return jpa.findAll().stream().map(mapper::toDomain).toList();
	}
}
