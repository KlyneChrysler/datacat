package domain

import (
	"math/rand/v2"
	"time"
)

// Persona is one synthetic traffic actor.
type Persona struct {
	Name       string
	SessionID  string
	UserAgent  string
	BaseDelay  time.Duration
	Jitter     time.Duration
	Paths      PathSource
	Credential *AgentCredential
}

func (p Persona) NextRequest() Request {
	return Request{SessionID: p.SessionID, Path: p.Paths.Next(), UserAgent: p.UserAgent, Credential: p.Credential}
}

// NextDelay is the base delay plus uniform jitter.
func (p Persona) NextDelay() time.Duration {
	if p.Jitter <= 0 {
		return p.BaseDelay
	}

	return p.BaseDelay + rand.N(p.Jitter)
}
