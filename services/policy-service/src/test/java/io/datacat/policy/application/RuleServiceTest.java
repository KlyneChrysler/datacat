package io.datacat.policy.application;

import io.datacat.policy.domain.ActionType;
import io.datacat.policy.domain.Classification;
import io.datacat.policy.domain.CreateRuleCommand;
import io.datacat.policy.domain.Rule;
import io.datacat.policy.domain.RuleNotFoundException;
import org.junit.jupiter.api.Test;

import java.util.ArrayList;
import java.util.List;
import java.util.Optional;
import java.util.UUID;

import static org.junit.jupiter.api.Assertions.assertEquals;
import static org.junit.jupiter.api.Assertions.assertThrows;

class RuleServiceTest {

	static final class FakeRepository implements RuleRepository {

		final List<Rule> saved = new ArrayList<>();

		@Override
		public Rule save(Rule rule) {
			saved.add(rule);
			return rule;
		}

		@Override
		public Optional<Rule> findById(UUID id) {
			return saved.stream().filter(r -> r.id().equals(id)).findFirst();
		}

		@Override
		public List<Rule> findAll() {
			return List.copyOf(saved);
		}
	}

	static final class FakePublisher implements RuleEventPublisher {

		final List<Rule> published = new ArrayList<>();

		@Override
		public void publishCreated(Rule rule) {
			published.add(rule);
		}
	}

	@Test
	void createSavesThenPublishes() {
		FakeRepository repository = new FakeRepository();
		FakePublisher publisher = new FakePublisher();
		RuleService service = new RuleService(repository, publisher);

		Rule created = service.create(new CreateRuleCommand("r", Classification.ABUSIVE, ActionType.BLOCK, 60));

		assertEquals(1, repository.saved.size());
		assertEquals(List.of(created), publisher.published);
	}

	@Test
	void getUnknownIdRaisesNotFound() {
		RuleService service = new RuleService(new FakeRepository(), new FakePublisher());

		assertThrows(RuleNotFoundException.class, () -> service.get(UUID.randomUUID()));
	}
}
