package httpapi

import "github.com/KlyneChrysler/datacat/pkg/policy"

// Domain to wire conversion.

func toDecisionResponse(d policy.Decision) decisionResponse {
	return decisionResponse{SessionID: d.SessionID, Class: string(d.Class), Action: string(d.Action)}
}

func toTrafficSummaryResponse(s policy.TrafficSummary) trafficSummaryResponse {
	return trafficSummaryResponse{WindowMinutes: s.WindowMinutes, Human: s.Human, VerifiedAgent: s.VerifiedAgent, Unverified: s.Unverified, Abusive: s.Abusive, Other: s.Other, Total: s.Total}
}
