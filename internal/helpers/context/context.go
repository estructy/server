// Package contexthelper provides utility functions to manage user IDs in a context.
package contexthelper

import "context"

type UserID string

// WithUserID adds a user ID to the context.
func WithUserID(ctx context.Context, userID UserID) context.Context {
	return context.WithValue(ctx, UserID("userID"), userID)
}

// UserIDFromContext retrieves the user ID from the context.
func UserIDFromContext(ctx context.Context) (UserID, bool) {
	userID, ok := ctx.Value(UserID("userID")).(UserID)
	if !ok {
		return "", false
	}
	return userID, true
}
