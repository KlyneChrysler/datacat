package kafka

import (
	"encoding/json"
	"fmt"

	"github.com/KlyneChrysler/datacat/pkg/events"
	policy "github.com/KlyneChrysler/datacat/pkg/policy"
)

// decodeVerdict converts one wire verdict into its domain form.
func decodeVerdict(value []byte) (policy.Verdict, error) {
	var wire events.Verdict
	if err := json.Unmarshal(value, &wire); err != nil {
		return policy.Verdict{}, fmt.Errorf("decode verdict: %w", err)
	}

	verdict, err := policy.NewVerdict(wire.SessionID, policy.Classification(wire.Classification), wire.Confidence)
	if err != nil {
		return policy.Verdict{}, fmt.Errorf("invalid verdict for session %q: %w", wire.SessionID, err)
	}

	return verdict, nil
}
