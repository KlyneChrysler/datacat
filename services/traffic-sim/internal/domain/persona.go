package domain

import (
	"math/rand/v2"
	"time"
)

// Persona describes one synthetic traffic actor: who it claims to be, how
// fast it moves, and how it navigates. Fields are set once; mutable
// navigation state lives inside its PathSource.
type Persona struct {
	Name       string
	SessionID  string
	UserAgent  string
	BaseDelay  time.Duration
	Jitter     time.Duration // uniform extra delay in [0, Jitter)
	Paths      PathSource
	Credential *AgentCredential // nil: persona does not sign
}

func (p Persona) NextRequest() Request {
	return Request{
		SessionID:  p.SessionID,
		Path:       p.Paths.Next(),
		UserAgent:  p.UserAgent,
		Credential: p.Credential,
	}
}

// NextDelay returns the pause before the next request: BaseDelay plus
// uniform jitter. A metronome persona simply sets Jitter near zero.
func (p Persona) NextDelay() time.Duration {
	if p.Jitter <= 0 {
		return p.BaseDelay
	}
	return p.BaseDelay + rand.N(p.Jitter)
}
