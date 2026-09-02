// Package events is the single source of truth for Kafka event schemas.
package events

import "time"

// RequestEvent is one observed request, keyed by session.
type RequestEvent struct {
	SessionID      string    `json:"session_id"`
	Timestamp      time.Time `json:"timestamp"`
	Method         string    `json:"method"`
	Path           string    `json:"path"`
	ClientIP       string    `json:"client_ip"`
	UserAgent      string    `json:"user_agent"`
	HeaderOrder    string    `json:"header_order"`    // hash of header names
	TLSFingerprint string    `json:"tls_fingerprint"` // empty until captured
	VerifiedAgent  bool      `json:"verified_agent"`  // valid trusted agent signature
}
