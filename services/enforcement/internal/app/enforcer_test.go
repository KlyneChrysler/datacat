package app

import (
	"context"
	"errors"
	"testing"

	policy "github.com/KlyneChrysler/datacat/pkg/policy"
)

var errBoom = errors.New("boom")

type fakeStore struct {
	saved []policy.Decision
	err   error
}

func (f *fakeStore) Save(_ context.Context, d policy.Decision) error {
	f.saved = append(f.saved, d)
	return f.err
}

func (f *fakeStore) FindBySession(_ context.Context, sessionID string) (policy.Decision, error) {
	for _, d := range f.saved {
		if d.SessionID == sessionID {
			return d, nil
		}
	}
	return policy.Decision{}, policy.ErrDecisionNotFound
}

type fakeApplier struct {
	applied []policy.Decision
	err     error
}

func (f *fakeApplier) Apply(_ context.Context, d policy.Decision) error {
	f.applied = append(f.applied, d)
	return f.err
}

type fakePublisher struct {
	published []policy.Decision
	err       error
}

func (f *fakePublisher) PublishDecision(_ context.Context, d policy.Decision) error {
	f.published = append(f.published, d)
	return f.err
}

func abusiveVerdict(t *testing.T) policy.Verdict {
	t.Helper()
	verdict, err := policy.NewVerdict("s-1", policy.Abusive, 0.95)
	if err != nil {
		t.Fatal(err)
	}
	return verdict
}

func TestEnforcerHandleVerdict(t *testing.T) {
	tests := []struct {
		name       string
		storeErr   error
		publishErr error
		applierErr error
		wantAction policy.Action
		wantErr    bool
	}{
		{name: "abusive session gets blocked", wantAction: policy.Block},
		{name: "store failure propagates", storeErr: errBoom, wantErr: true},
		{name: "publish failure propagates", publishErr: errBoom, wantErr: true},
		{name: "applier failure propagates", applierErr: errBoom, wantErr: true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			store := &fakeStore{err: tt.storeErr}
			publisher := &fakePublisher{err: tt.publishErr}
			applier := &fakeApplier{err: tt.applierErr}
			enforcer := NewEnforcer(policy.DefaultPolicy(), store, publisher, applier, NewTally())

			err := enforcer.HandleVerdict(context.Background(), abusiveVerdict(t))

			if tt.wantErr != (err != nil) {
				t.Fatalf("err = %v, wantErr = %v", err, tt.wantErr)
			}
			if !tt.wantErr && store.saved[0].Action != tt.wantAction {
				t.Errorf("stored action = %s, want %s", store.saved[0].Action, tt.wantAction)
			}
			if !tt.wantErr && (len(publisher.published) != 1 || len(applier.applied) != 1) {
				t.Errorf("published %d, applied %d decisions, want 1 and 1",
					len(publisher.published), len(applier.applied))
			}
		})
	}
}

func TestEnforcerLookupReturnsStoredDecision(t *testing.T) {
	store := &fakeStore{}
	enforcer := NewEnforcer(policy.DefaultPolicy(), store, &fakePublisher{}, &fakeApplier{}, NewTally())
	if err := enforcer.HandleVerdict(context.Background(), abusiveVerdict(t)); err != nil {
		t.Fatal(err)
	}

	decision, err := enforcer.Lookup(context.Background(), "s-1")

	if err != nil {
		t.Fatal(err)
	}
	if decision.Action != policy.Block {
		t.Errorf("action = %s, want block", decision.Action)
	}
}
