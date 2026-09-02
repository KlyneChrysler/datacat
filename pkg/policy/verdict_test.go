package policy

import (
	"errors"
	"testing"
)

func TestNewVerdict(t *testing.T) {
	tests := []struct {
		name       string
		sessionID  string
		confidence float64
		wantErr    error
	}{
		{name: "valid verdict", sessionID: "s-1", confidence: 0.9},
		{name: "empty session id rejected", sessionID: "", confidence: 0.9, wantErr: ErrEmptySessionID},
		{name: "confidence above one rejected", sessionID: "s-1", confidence: 1.1, wantErr: ErrConfidenceRange},
		{name: "negative confidence rejected", sessionID: "s-1", confidence: -0.1, wantErr: ErrConfidenceRange},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := NewVerdict(tt.sessionID, Abusive, tt.confidence)
			if !errors.Is(err, tt.wantErr) {
				t.Fatalf("err = %v, want %v", err, tt.wantErr)
			}
		})
	}
}

func TestVerdictWithClassIsImmutable(t *testing.T) {
	original, err := NewVerdict("s-1", Unverified, 0.5)
	if err != nil {
		t.Fatal(err)
	}

	changed := original.WithClass(Abusive)

	if original.Class != Unverified {
		t.Errorf("original mutated: class = %s", original.Class)
	}
	if changed.Class != Abusive {
		t.Errorf("copy not updated: class = %s", changed.Class)
	}
}
