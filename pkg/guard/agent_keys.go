package guard

import (
	"crypto/ed25519"
	"encoding/hex"
	"fmt"
	"strings"
)

// ParseAgentKeys reads keyid=hexpubkey pairs, failing on any bad entry.
func ParseAgentKeys(raw string) (map[string]ed25519.PublicKey, error) {
	keys := make(map[string]ed25519.PublicKey)
	if raw == "" {
		return keys, nil
	}

	for _, entry := range strings.Split(raw, ";") {
		keyID, hexKey, found := strings.Cut(entry, "=")
		if !found {
			return nil, fmt.Errorf("config: AGENT_KEYS entry %q is not keyid=hexpubkey", entry)
		}
		pub, err := hex.DecodeString(hexKey)
		if err != nil || len(pub) != ed25519.PublicKeySize {
			return nil, fmt.Errorf("config: AGENT_KEYS %q has an invalid public key", keyID)
		}
		keys[keyID] = ed25519.PublicKey(pub)
	}

	return keys, nil
}
