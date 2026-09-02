package observe

import (
	"net/http"
	"slices"
	"strings"
	"time"

	"github.com/KlyneChrysler/datacat/pkg/events"
	"github.com/KlyneChrysler/datacat/pkg/hashx"
	"github.com/KlyneChrysler/datacat/services/edge-proxy/internal/adapters/agentauth"
	"github.com/KlyneChrysler/datacat/services/edge-proxy/internal/adapters/ident"
)

// Request → wire conversion only (file taxonomy, standards rule 2).

func eventFrom(r *http.Request) events.RequestEvent {
	return events.RequestEvent{
		SessionID:      ident.SessionID(r),
		Timestamp:      time.Now().UTC(),
		Method:         r.Method,
		Path:           r.URL.Path,
		ClientIP:       ident.ClientIP(r),
		UserAgent:      r.UserAgent(),
		HeaderOrder:    headerOrderHash(r),
		TLSFingerprint: "", // requires raw ClientHello capture; later phase
		VerifiedAgent:  agentauth.Verified(r.Context()),
	}
}

// headerOrderHash approximates a header fingerprint. net/http stores headers
// in a map, losing wire order, so this hashes the sorted name set; true
// wire-order capture needs a lower-level listener (later phase).
// O(h log h) in header count, bounded by the server's header limits.
func headerOrderHash(r *http.Request) string {
	names := make([]string, 0, len(r.Header))
	for name := range r.Header {
		names = append(names, name)
	}
	slices.Sort(names)
	return hashx.Short(strings.Join(names, ","))
}
