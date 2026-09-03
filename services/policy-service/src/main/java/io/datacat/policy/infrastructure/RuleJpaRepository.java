package io.datacat.policy.infrastructure;

import org.springframework.data.jpa.repository.JpaRepository;

import java.util.UUID;

/** Spring Data access to the rules table. */
public interface RuleJpaRepository extends JpaRepository<RuleEntity, UUID> {
}
