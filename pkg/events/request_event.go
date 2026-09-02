// Package events is the single source of truth for event schemas flowing
// through Kafka. Every producer and consumer (Go services, Flink job via
// mirrored Java types) must match these shapes exactly. One schema per file.
package events

import "time"

// RequestEvent is emitted by edge-proxy for every observed request.
// Partition key: SessionID.
type RequestEvent struct {
	SessionID      string    `json:"session_id"`
	Timestamp      time.Time `json:"timestamp"`
	Method         string    `json:"method"`
	Path           string    `json:"path"`
	ClientIP       string    `json:"client_ip"`
	UserAgent      string    `json:"user_agent"`
	HeaderOrder    string    `json:"header_order"`    // hash of header names in wire order
	TLSFingerprint string    `json:"tls_fingerprint"` // e.g. JA4, empty if unavailable
	VerifiedAgent  bool      `json:"verified_agent"`  // request carried a valid trusted-agent signature
}
