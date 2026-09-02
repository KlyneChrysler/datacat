package classify

import "strings"

// userAgentScorer reads a declared automation agent as strong evidence.
type userAgentScorer struct{}

var automationMarkers = []string{"curl", "wget", "python", "go-http-client", "bot", "crawler", "spider", "scrapy", "headless", "phantom", "httpclient"}

func (userAgentScorer) Score(f Features) float64 {
	if f.UserAgent == "" {
		return 0.7
	}

	lowered := strings.ToLower(f.UserAgent)
	for _, marker := range automationMarkers {
		if strings.Contains(lowered, marker) {
			return 0.95
		}
	}

	return 0.15
}
