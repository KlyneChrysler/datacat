package io.datacat.policy.application;

import io.datacat.policy.domain.Rule;

import java.util.List;
import java.util.Optional;
import java.util.UUID;

/** What the service needs from storage. */
public interface RuleRepository {

	Rule save(Rule rule);

	Optional<Rule> findById(UUID id);

	List<Rule> findAll();
}
