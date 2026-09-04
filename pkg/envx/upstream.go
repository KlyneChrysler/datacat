package envx

import (
	"fmt"
	"net/url"
	"os"
)

// Upstream reads and parses the UPSTREAM_URL env var, required.
func Upstream() (*url.URL, error) {
	raw := os.Getenv("UPSTREAM_URL")
	if raw == "" {
		return nil, fmt.Errorf("config: UPSTREAM_URL is required")
	}

	upstream, err := url.Parse(raw)
	if err != nil {
		return nil, fmt.Errorf("config: UPSTREAM_URL invalid: %w", err)
	}

	return upstream, nil
}
