// Package agentauth verifies trusted agent request signatures.
package agentauth

import "context"

type contextKey struct{}

func withVerified(ctx context.Context) context.Context {
	return context.WithValue(ctx, contextKey{}, true)
}

// Verified reports whether this request carried a valid agent signature.
func Verified(ctx context.Context) bool {
	verified, ok := ctx.Value(contextKey{}).(bool)

	return ok && verified
}
