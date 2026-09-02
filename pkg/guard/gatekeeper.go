package guard

import (
	"sync"
	"time"
)

// Gatekeeper remembers the latest action per session with a ttl.
type Gatekeeper struct {
	ttl  time.Duration
	mu   sync.RWMutex
	byID map[string]gateEntry
}

func NewGatekeeper(ttl time.Duration) *Gatekeeper {
	return &Gatekeeper{ttl: ttl, byID: make(map[string]gateEntry)}
}

// Update records the newest action for a session.
func (g *Gatekeeper) Update(sessionID, action string) {
	g.mu.Lock()
	defer g.mu.Unlock()

	g.byID[sessionID] = gateEntry{action: action, savedAt: time.Now()}
}

// ActionFor returns the standing action, empty for unknown or expired sessions.
func (g *Gatekeeper) ActionFor(sessionID string) string {
	g.mu.RLock()
	entry, ok := g.byID[sessionID]
	g.mu.RUnlock()

	if !ok {
		return ""
	}
	if time.Since(entry.savedAt) > g.ttl {
		g.forget(sessionID)
		return ""
	}

	return entry.action
}

func (g *Gatekeeper) forget(sessionID string) {
	g.mu.Lock()
	defer g.mu.Unlock()

	delete(g.byID, sessionID)
}
