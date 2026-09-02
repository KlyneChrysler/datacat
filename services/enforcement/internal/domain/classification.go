// Package domain holds pure business logic with zero infra imports.
package domain

type Classification string

const (
	Human       Classification = "human"
	VerifiedBot Classification = "verified_agent"
	Unverified  Classification = "unverified_automation"
	Abusive     Classification = "abusive"
)
