package io.datacat.policy.api;

import io.datacat.policy.domain.ActionType;
import io.datacat.policy.domain.Classification;
import io.datacat.policy.domain.CreateRuleCommand;
import jakarta.validation.constraints.Max;
import jakarta.validation.constraints.Min;
import jakarta.validation.constraints.NotBlank;
import jakarta.validation.constraints.NotNull;

/** Wire shape of a rule creation, validated at the boundary. */
public record CreateRuleRequest(@NotBlank String name, @NotNull Classification appliesTo, @NotNull ActionType action, @Min(1) @Max(10_000) int rateLimitPerMinute) {

	public CreateRuleCommand toCommand() {
		return new CreateRuleCommand(name, appliesTo, action, rateLimitPerMinute);
	}
}
