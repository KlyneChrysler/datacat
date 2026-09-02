package httpapi

import "github.com/KlyneChrysler/datacat/services/enforcement/internal/domain"

// Domain → wire conversion only (file taxonomy, standards rule 2).

func toDecisionResponse(d domain.Decision) decisionResponse {
	return decisionResponse{SessionID: d.SessionID, Class: string(d.Class), Action: string(d.Action)}
}
