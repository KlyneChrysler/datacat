package domain

import "errors"

var (
	ErrEmptySessionID   = errors.New("session id must not be empty")
	ErrConfidenceRange  = errors.New("confidence must be within [0,1]")
	ErrVerdictNotFound  = errors.New("verdict not found")
	ErrDecisionNotFound = errors.New("decision not found")
)
