// Package domain holds pure business logic. It imports zero infrastructure.
// One concept per file (file taxonomy, standards rule 2).
package domain

type Classification string

const (
	Human       Classification = "human"
	VerifiedBot Classification = "verified_agent"
	Unverified  Classification = "unverified_automation"
	Abusive     Classification = "abusive"
)
