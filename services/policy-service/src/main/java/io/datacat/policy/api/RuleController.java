package io.datacat.policy.api;

import io.datacat.policy.application.RuleService;
import io.datacat.policy.domain.Rule;
import jakarta.validation.Valid;
import org.springframework.http.HttpStatus;
import org.springframework.web.bind.annotation.GetMapping;
import org.springframework.web.bind.annotation.PathVariable;
import org.springframework.web.bind.annotation.PostMapping;
import org.springframework.web.bind.annotation.RequestBody;
import org.springframework.web.bind.annotation.RequestMapping;
import org.springframework.web.bind.annotation.ResponseStatus;
import org.springframework.web.bind.annotation.RestController;

import java.util.List;
import java.util.UUID;

/** Thin HTTP adapter, parse then delegate then map. */
@RestController
@RequestMapping("/v1/rules")
public class RuleController {

	private final RuleService ruleService;

	public RuleController(RuleService ruleService) {
		this.ruleService = ruleService;
	}

	@PostMapping
	@ResponseStatus(HttpStatus.CREATED)
	public RuleResponse create(@Valid @RequestBody CreateRuleRequest request) {
		Rule rule = ruleService.create(request.toCommand());

		return RuleResponse.from(rule);
	}

	@GetMapping("/{id}")
	public RuleResponse get(@PathVariable UUID id) {
		return RuleResponse.from(ruleService.get(id));
	}

	@GetMapping
	public List<RuleResponse> list() {
		return ruleService.list().stream().map(RuleResponse::from).toList();
	}
}
