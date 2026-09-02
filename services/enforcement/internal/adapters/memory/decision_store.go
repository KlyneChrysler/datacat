// Package memory holds the in memory decision store for single replicas.
package memory

import (
	"context"
	"sync"

	"github.com/KlyneChrysler/datacat/pkg/policy"
	"github.com/KlyneChrysler/datacat/services/enforcement/internal/ports"
)

type DecisionStore struct {
	mu        sync.RWMutex
	bySession map[string]policy.Decision
}

var _ ports.DecisionStore = (*DecisionStore)(nil)

func NewDecisionStore() *DecisionStore {
	return &DecisionStore{bySession: make(map[string]policy.Decision)}
}

func (s *DecisionStore) Save(_ context.Context, d policy.Decision) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	s.bySession[d.SessionID] = d

	return nil
}

func (s *DecisionStore) FindBySession(_ context.Context, sessionID string) (policy.Decision, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	decision, ok := s.bySession[sessionID]
	if !ok {
		return policy.Decision{}, policy.ErrDecisionNotFound
	}

	return decision, nil
}
