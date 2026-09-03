package io.datacat.policy.infrastructure;

import io.datacat.policy.domain.ActionType;
import io.datacat.policy.domain.Classification;
import io.datacat.policy.domain.CreateRuleCommand;
import io.datacat.policy.domain.Rule;
import org.junit.jupiter.api.Test;
import org.springframework.boot.jdbc.test.autoconfigure.AutoConfigureTestDatabase;
import org.springframework.boot.data.jpa.test.autoconfigure.DataJpaTest;
import org.springframework.context.annotation.Import;
import org.testcontainers.containers.PostgreSQLContainer;
import org.testcontainers.junit.jupiter.Container;
import org.testcontainers.junit.jupiter.Testcontainers;

import java.util.List;
import java.util.Optional;
import java.util.UUID;

import static org.junit.jupiter.api.Assertions.assertEquals;
import static org.junit.jupiter.api.Assertions.assertTrue;
import static org.springframework.boot.jdbc.test.autoconfigure.AutoConfigureTestDatabase.Replace.NONE;

/** Repository over real Postgres with Flyway applied. */
@DataJpaTest
@AutoConfigureTestDatabase(replace = NONE)
@Import({RuleRepositoryAdapter.class, RuleEntityMapper.class})
@Testcontainers
class RuleRepositoryAdapterTest {

	@Container
	@org.springframework.boot.testcontainers.service.connection.ServiceConnection
	static PostgreSQLContainer<?> postgres = new PostgreSQLContainer<>("postgres:17-alpine");

	@org.springframework.beans.factory.annotation.Autowired
	private RuleRepositoryAdapter adapter;

	@Test
	void saveThenFindByIdRoundTrip() {
		Rule rule = Rule.from(new CreateRuleCommand("throttle agents", Classification.VERIFIED_AGENT, ActionType.RATE_LIMIT, 90));

		adapter.save(rule);
		Optional<Rule> found = adapter.findById(rule.id());

		assertTrue(found.isPresent());
		assertEquals("throttle agents", found.get().name());
		assertEquals(ActionType.RATE_LIMIT, found.get().action());
	}

	@Test
	void findAllReturnsSavedRules() {
		adapter.save(Rule.from(new CreateRuleCommand("a", Classification.ABUSIVE, ActionType.BLOCK, 1)));
		adapter.save(Rule.from(new CreateRuleCommand("b", Classification.UNVERIFIED_AUTOMATION, ActionType.CHALLENGE, 1)));

		List<Rule> all = adapter.findAll();

		assertTrue(all.size() >= 2);
	}

	@Test
	void findByUnknownIdIsEmpty() {
		assertTrue(adapter.findById(UUID.randomUUID()).isEmpty());
	}
}
