package observe

import (
	"net/http"
	"slices"
	"strings"
	"time"

	"github.com/KlyneChrysler/datacat/pkg/edge/agentauth"
	"github.com/KlyneChrysler/datacat/pkg/edge/ident"
	"github.com/KlyneChrysler/datacat/pkg/events"
	"github.com/KlyneChrysler/datacat/pkg/hashx"
)

// eventFrom converts one request into its wire event.
func eventFrom(r *http.Request) events.RequestEvent {
	return events.RequestEvent{
		SessionID:      ident.SessionID(r),
		Timestamp:      time.Now().UTC(),
		Method:         r.Method,
		Path:           r.URL.Path,
		ClientIP:       ident.ClientIP(r),
		UserAgent:      r.UserAgent(),
		HeaderOrder:    headerOrderHash(r),
		TLSFingerprint: "",
		VerifiedAgent:  agentauth.Verified(r.Context()),
	}
}

// headerOrderHash fingerprints the sorted header name set.
func headerOrderHash(r *http.Request) string {
	names := make([]string, 0, len(r.Header))
	for name := range r.Header {
		names = append(names, name)
	}

	slices.Sort(names)

	return hashx.Short(strings.Join(names, ","))
}
