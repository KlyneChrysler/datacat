package domain

type Action string

const (
	Allow     Action = "allow"
	Challenge Action = "challenge"
	RateLimit Action = "rate_limit"
	Block     Action = "block"
)
