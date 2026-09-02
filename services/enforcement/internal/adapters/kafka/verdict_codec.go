package kafka

import (
	"encoding/json"
	"fmt"

	"github.com/KlyneChrysler/datacat/pkg/events"
	"github.com/KlyneChrysler/datacat/services/enforcement/internal/domain"
)

// decodeVerdict converts one wire verdict into its domain form.
func decodeVerdict(value []byte) (domain.Verdict, error) {
	var wire events.Verdict
	if err := json.Unmarshal(value, &wire); err != nil {
		return domain.Verdict{}, fmt.Errorf("decode verdict: %w", err)
	}

	verdict, err := domain.NewVerdict(wire.SessionID, domain.Classification(wire.Classification), wire.Confidence)
	if err != nil {
		return domain.Verdict{}, fmt.Errorf("invalid verdict for session %q: %w", wire.SessionID, err)
	}

	return verdict, nil
}
