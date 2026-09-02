package httpapi

import "github.com/KlyneChrysler/datacat/services/enforcement/internal/domain"

// Domain → wire conversion only (file taxonomy, standards rule 2).

func toDecisionResponse(d domain.Decision) decisionResponse {
	return decisionResponse{SessionID: d.SessionID, Class: string(d.Class), Action: string(d.Action)}
}

func toTrafficSummaryResponse(s domain.TrafficSummary) trafficSummaryResponse {
	return trafficSummaryResponse{
		WindowMinutes: s.WindowMinutes,
		Human:         s.Human,
		VerifiedAgent: s.VerifiedAgent,
		Unverified:    s.Unverified,
		Abusive:       s.Abusive,
		Other:         s.Other,
		Total:         s.Total,
	}
}
