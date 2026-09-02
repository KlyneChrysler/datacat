package app

import (
	"sync"
	"time"

	"github.com/KlyneChrysler/datacat/pkg/events"
)

// Gatekeeper remembers the latest decision per session with a ttl.
type Gatekeeper struct {
	ttl  time.Duration
	mu   sync.RWMutex
	byID map[string]gateEntry
}

func NewGatekeeper(ttl time.Duration) *Gatekeeper {
	return &Gatekeeper{ttl: ttl, byID: make(map[string]gateEntry)}
}

// Update records the newest decision for a session.
func (g *Gatekeeper) Update(d events.Decision) {
	g.mu.Lock()
	defer g.mu.Unlock()

	g.byID[d.SessionID] = gateEntry{action: d.Action, savedAt: time.Now()}
}

// ActionFor returns the standing action, allowing unknown or expired sessions.
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
