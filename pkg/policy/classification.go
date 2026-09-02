// Package domain holds pure business logic with zero infra imports.
package policy

type Classification string

const (
	Human       Classification = "human"
	VerifiedBot Classification = "verified_agent"
	Unverified  Classification = "unverified_automation"
	Abusive     Classification = "abusive"
)
