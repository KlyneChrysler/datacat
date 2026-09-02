// Package httpsender implements ports.Sender over HTTP: session cookie,
// declared User-Agent, body drained and closed so connections are reused.
package httpsender

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/KlyneChrysler/datacat/pkg/events"
	"github.com/KlyneChrysler/datacat/services/traffic-sim/internal/domain"
	"github.com/KlyneChrysler/datacat/services/traffic-sim/internal/ports"
)

// maxDrainBytes caps how much of a response body is read for connection
// reuse — a huge body must not stall the persona's cadence.
const maxDrainBytes = 64 * 1024

type Sender struct {
	client  *http.Client
	baseURL string
}

var _ ports.Sender = (*Sender)(nil)

func NewSender(baseURL string) *Sender {
	return &Sender{
		client:  &http.Client{Timeout: 10 * time.Second},
		baseURL: baseURL,
	}
}

func (s *Sender) Send(ctx context.Context, req domain.Request) (int, error) {
	httpReq, err := s.build(ctx, req)
	if err != nil {
		return 0, err
	}
	resp, err := s.client.Do(httpReq)
	if err != nil {
		return 0, fmt.Errorf("send %s: %w", req.Path, err)
	}
	defer resp.Body.Close()
	_, _ = io.Copy(io.Discard, io.LimitReader(resp.Body, maxDrainBytes)) // drain for reuse; content is irrelevant
	return resp.StatusCode, nil
}

func (s *Sender) build(ctx context.Context, req domain.Request) (*http.Request, error) {
	httpReq, err := http.NewRequestWithContext(ctx, http.MethodGet, s.baseURL+req.Path, nil)
	if err != nil {
		return nil, fmt.Errorf("build request %s: %w", req.Path, err)
	}
	httpReq.Header.Set("User-Agent", req.UserAgent)
	httpReq.AddCookie(&http.Cookie{Name: events.SessionCookie, Value: req.SessionID})
	sign(httpReq, req)
	return httpReq, nil
}
