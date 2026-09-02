package domain

import "time"

// DefaultPersonas is the registry of standard actors, one entry per persona.
func DefaultPersonas(agentCred *AgentCredential) []Persona {
	return []Persona{
		{
			Name:      "human",
			SessionID: "sim-human",
			UserAgent: "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/605.1.15 (KHTML, like Gecko) Version/17.4 Safari/605.1.15",
			BaseDelay: 1200 * time.Millisecond,
			Jitter:    2500 * time.Millisecond,
			Paths:     NewLoopingPaths([]string{"/home", "/products", "/products", "/cart", "/home"}),
		},
		{
			Name:       "polite-agent",
			SessionID:  "sim-agent",
			UserAgent:  "datacat-agent/1.0 (+https://github.com/KlyneChrysler/datacat)",
			BaseDelay:  800 * time.Millisecond,
			Jitter:     100 * time.Millisecond,
			Paths:      NewLoopingPaths([]string{"/api/catalog", "/api/prices"}),
			Credential: agentCred,
		},
		{
			Name:      "scraper",
			SessionID: "sim-scraper",
			UserAgent: "python-requests/2.32",
			BaseDelay: 350 * time.Millisecond,
			Jitter:    30 * time.Millisecond,
			Paths:     NewCrawlingPaths("/catalog/item-"),
		},
	}
}
