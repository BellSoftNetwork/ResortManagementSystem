package audit

import (
	"context"
)

// contextKey is a private type for context keys to avoid collisions
type contextKey string

const userContextKey contextKey = "audit_user_context"

// SetUserContext sets user context in the request context for audit logging
func SetUserContext(ctx context.Context, userID *uint, username string) context.Context {
	userCtx := &UserContext{
		Username: username,
	}
	if userID != nil {
		userCtx.UserID = *userID
	}
	return context.WithValue(ctx, userContextKey, userCtx)
}

// GetUserContext retrieves user context from the request context
// Returns empty UserContext if no user context is found
func GetUserContext(ctx context.Context) *UserContext {
	if userCtx, ok := ctx.Value(userContextKey).(*UserContext); ok {
		return userCtx
	}
	return &UserContext{} // Return empty context if not found
}

// HasUserContext checks if user context exists in the request context
func HasUserContext(ctx context.Context) bool {
	_, ok := ctx.Value(userContextKey).(*UserContext)
	return ok
}
