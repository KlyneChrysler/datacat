package classify

import "time"

const (
	maxSamples       = 128
	maxDistinctPaths = 64
)

// sessionState is the stored history for one session.
type sessionState struct {
	timestamps [maxSamples]time.Time
	sampleIdx  int
	sampleLen  int
	pathCounts map[string]int64
	pathTotal  int64
	count      int64
	verified   int64
	userAgent  string
	lastClass  string
	lastEval   time.Time
	lastSeen   time.Time
}

func newSessionState() *sessionState {
	return &sessionState{pathCounts: make(map[string]int64, 8)}
}

func (s *sessionState) sample(o Observation) {
	s.timestamps[s.sampleIdx] = o.At
	s.sampleIdx = (s.sampleIdx + 1) % maxSamples
	if s.sampleLen < maxSamples {
		s.sampleLen++
	}

	if _, known := s.pathCounts[o.Path]; known || len(s.pathCounts) < maxDistinctPaths {
		s.pathCounts[o.Path]++
		s.pathTotal++
	}

	s.count++
	if o.Verified {
		s.verified++
	}
	s.userAgent = o.UserAgent
	s.lastSeen = o.At
}
