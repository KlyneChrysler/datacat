package io.datacat.policy.infrastructure;

import jakarta.persistence.Column;
import jakarta.persistence.Entity;
import jakarta.persistence.Id;
import jakarta.persistence.Table;

import java.time.Instant;
import java.util.UUID;

/** Storage shape of one rule, mapped to and from the domain record. */
@Entity
@Table(name = "rules")
public class RuleEntity {

	@Id
	private UUID id;

	@Column(nullable = false)
	private String name;

	@Column(name = "applies_to", nullable = false)
	private String appliesTo;

	@Column(nullable = false)
	private String action;

	@Column(name = "rate_limit_per_minute", nullable = false)
	private int rateLimitPerMinute;

	@Column(name = "created_at", nullable = false)
	private Instant createdAt;

	protected RuleEntity() {
	}

	public RuleEntity(UUID id, String name, String appliesTo, String action, int rateLimitPerMinute, Instant createdAt) {
		this.id = id;
		this.name = name;
		this.appliesTo = appliesTo;
		this.action = action;
		this.rateLimitPerMinute = rateLimitPerMinute;
		this.createdAt = createdAt;
	}

	public UUID getId() {
		return id;
	}

	public String getName() {
		return name;
	}

	public String getAppliesTo() {
		return appliesTo;
	}

	public String getAction() {
		return action;
	}

	public int getRateLimitPerMinute() {
		return rateLimitPerMinute;
	}

	public Instant getCreatedAt() {
		return createdAt;
	}
}
