package classify

import (
	"sync"
	"time"

	policy "github.com/KlyneChrysler/datacat/pkg/policy"
)

const (
	maxSessions  = 10_000
	evalInterval = 2 * time.Second
	rateWindow   = time.Minute
	idleTTL      = 10 * time.Minute
)

// Tracker keeps per session history and reclassifies at a bounded pace.
// Sampling is O(1), evaluation is O(128) at most once per 2s per session.
type Tracker struct {
	scorers  []Scorer
	mu       sync.Mutex
	sessions map[string]*sessionState
}

func NewTracker() *Tracker {
	return &Tracker{scorers: defaultScorers(), sessions: make(map[string]*sessionState)}
}

// Observe records one request, returning a verdict when the class changes.
func (t *Tracker) Observe(o Observation) (policy.Verdict, bool) {
	t.mu.Lock()
	defer t.mu.Unlock()

	state := t.track(o.SessionID, o.At)
	if state == nil {
		return policy.Verdict{}, false
	}

	state.sample(o)
	if o.At.Sub(state.lastEval) < evalInterval {
		return policy.Verdict{}, false
	}
	state.lastEval = o.At

	return t.evaluate(o.SessionID, state, o.At)
}

// track finds or registers the session, unknown means untracked at cap.
func (t *Tracker) track(sessionID string, now time.Time) *sessionState {
	state, ok := t.sessions[sessionID]
	if ok {
		return state
	}

	if len(t.sessions) >= maxSessions {
		t.sweep(now)
	}
	if len(t.sessions) >= maxSessions {
		return nil
	}

	state = newSessionState()
	t.sessions[sessionID] = state

	return state
}

// sweep drops idle sessions, runs only at the cap.
func (t *Tracker) sweep(now time.Time) {
	for id, state := range t.sessions {
		if now.Sub(state.lastSeen) > idleTTL {
			delete(t.sessions, id)
		}
	}
}

func (t *Tracker) evaluate(sessionID string, state *sessionState, now time.Time) (policy.Verdict, bool) {
	features := featuresOf(state, now)
	scores := make([]float64, 0, len(t.scorers))
	for _, scorer := range t.scorers {
		scores = append(scores, scorer.Score(features))
	}

	class, confidence := classifyScores(scores, features.VerifiedShare)
	if string(class) == state.lastClass {
		return policy.Verdict{}, false
	}
	state.lastClass = string(class)

	verdict, err := policy.NewVerdict(sessionID, class, confidence)
	return verdict, err == nil
}

// featuresOf turns buffered samples into features.
func featuresOf(state *sessionState, now time.Time) Features {
	samples := orderedSamples(state)

	return Features{
		RequestCount:      state.count,
		RequestsPerMinute: recentPerMinute(samples, now),
		IntervalCv:        intervalCv(samples),
		PathEntropy:       normalizedEntropy(state.pathCounts, state.pathTotal),
		UserAgent:         state.userAgent,
		VerifiedShare:     verifiedShareOf(state),
	}
}

// orderedSamples reads the ring oldest first.
func orderedSamples(state *sessionState) []time.Time {
	samples := make([]time.Time, 0, state.sampleLen)
	start := (state.sampleIdx - state.sampleLen + maxSamples) % maxSamples
	for i := 0; i < state.sampleLen; i++ {
		samples = append(samples, state.timestamps[(start+i)%maxSamples])
	}

	return samples
}

func recentPerMinute(samples []time.Time, now time.Time) float64 {
	cutoff := now.Add(-rateWindow)
	recent := 0
	for _, ts := range samples {
		if ts.After(cutoff) {
			recent++
		}
	}

	return float64(recent)
}

func verifiedShareOf(state *sessionState) float64 {
	if state.count == 0 {
		return 0
	}

	return float64(state.verified) / float64(state.count)
}
