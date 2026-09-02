// Package agentauth is the inbound adapter verifying trusted-agent request
// signatures and exposing the result to downstream adapters via context.
package agentauth

import "context"

type contextKey struct{}

func withVerified(ctx context.Context) context.Context {
	return context.WithValue(ctx, contextKey{}, true)
}

// Verified reports whether this request carried a valid trusted-agent
// signature.
func Verified(ctx context.Context) bool {
	verified, ok := ctx.Value(contextKey{}).(bool)
	return ok && verified
}
