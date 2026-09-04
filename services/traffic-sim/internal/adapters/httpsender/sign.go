package httpsender

import (
	"crypto/ed25519"
	"encoding/base64"
	"net/http"
	"strconv"
	"time"

	"github.com/KlyneChrysler/datacat/pkg/agentsig"
	"github.com/KlyneChrysler/datacat/services/traffic-sim/internal/domain"
)

// sign attaches the trusted agent signature headers when the request has a credential.
func sign(httpReq *http.Request, req domain.Request) {
	if req.Credential == nil {
		return
	}

	issued := time.Now().Unix()
	base := agentsig.Base(httpReq.Method, httpReq.URL.Path, req.SessionID, issued)

	httpReq.Header.Set(agentsig.HeaderKey, req.Credential.KeyID)
	httpReq.Header.Set(agentsig.HeaderTimestamp, strconv.FormatInt(issued, 10))
	httpReq.Header.Set(agentsig.HeaderSignature, base64.StdEncoding.EncodeToString(ed25519.Sign(req.Credential.Key, []byte(base))))
}
