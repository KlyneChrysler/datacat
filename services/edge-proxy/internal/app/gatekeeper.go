package app

import (
	"sync"
	"time"

	"github.com/KlyneChrysler/datacat/pkg/events"
)

// Gatekeeper holds the latest enforcement decision per session, fed by the
// decisions topic and consulted on the hot request path. Entries expire so a
// session that stops misbehaving is eventually forgiven and memory stays
// bounded; expiry is checked lazily on read.
type Gatekeeper struct {
	ttl  time.Duration
	mu   sync.RWMutex
	byID map[string]gateEntry
}

type gateEntry struct {
	action  string
	savedAt time.Time
}

func NewGatekeeper(ttl time.Duration) *Gatekeeper {
	return &Gatekeeper{ttl: ttl, byID: make(map[string]gateEntry)}
}

// Update records a decision. Used as the decision-source handler.
func (g *Gatekeeper) Update(d events.Decision) {
	g.mu.Lock()
	defer g.mu.Unlock()
	g.byID[d.SessionID] = gateEntry{action: d.Action, savedAt: time.Now()}
}

// ActionFor returns the standing action for a session; unknown or expired
// sessions are allowed — the classifier will judge them again.
func (g *Gatekeeper) ActionFor(sessionID string) string {
	g.mu.RLock()
	entry, ok := g.byID[sessionID]
	g.mu.RUnlock()
	if !ok {
		return "allow"
	}
	if time.Since(entry.savedAt) > g.ttl {
		g.forget(sessionID)
		return "allow"
	}
	return entry.action
}

func (g *Gatekeeper) forget(sessionID string) {
	g.mu.Lock()
	defer g.mu.Unlock()
	delete(g.byID, sessionID)
}
