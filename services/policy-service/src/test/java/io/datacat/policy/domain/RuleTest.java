package io.datacat.policy.domain;

import org.junit.jupiter.api.Test;

import static org.junit.jupiter.api.Assertions.assertEquals;
import static org.junit.jupiter.api.Assertions.assertNotNull;
import static org.junit.jupiter.api.Assertions.assertThrows;

class RuleTest {

	@Test
	void fromAssignsIdentityAndTimestamp() {
		Rule rule = Rule.from(new CreateRuleCommand("slow agents", Classification.VERIFIED_AGENT, ActionType.RATE_LIMIT, 60));

		assertNotNull(rule.id());
		assertNotNull(rule.createdAt());
		assertEquals("slow agents", rule.name());
		assertEquals(ActionType.RATE_LIMIT, rule.action());
	}

	@Test
	void blockingHumansIsRejected() {
		CreateRuleCommand command = new CreateRuleCommand("bad rule", Classification.HUMAN, ActionType.BLOCK, 60);

		assertThrows(IncompatibleRuleException.class, () -> Rule.from(command));
	}
}
