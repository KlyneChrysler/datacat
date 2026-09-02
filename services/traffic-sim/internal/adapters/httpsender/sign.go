package httpsender

import (
	"crypto/ed25519"
	"encoding/base64"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/KlyneChrysler/datacat/services/traffic-sim/internal/domain"
)

// Signature headers, must mirror the edge proxy agentauth adapter.
const (
	headerKey       = "X-Agent-Key"
	headerSignature = "X-Agent-Signature"
	headerTimestamp = "X-Agent-Timestamp"
)

// sign adds the trusted agent headers when the request has a credential.
func sign(httpReq *http.Request, req domain.Request) {
	if req.Credential == nil {
		return
	}

	issued := time.Now().Unix()
	base := strings.Join([]string{httpReq.Method, httpReq.URL.Path, req.SessionID, strconv.FormatInt(issued, 10)}, "\n")

	httpReq.Header.Set(headerKey, req.Credential.KeyID)
	httpReq.Header.Set(headerTimestamp, strconv.FormatInt(issued, 10))
	httpReq.Header.Set(headerSignature, base64.StdEncoding.EncodeToString(ed25519.Sign(req.Credential.Key, []byte(base))))
}
