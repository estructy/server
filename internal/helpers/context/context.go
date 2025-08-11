// Package contexthelper provides utility functions to manage user IDs in a context.
package contexthelper

import (
	"context"

	"github.com/google/uuid"
)

type UserID string

// WithUserID adds a user ID to the context.
func WithUserID(ctx context.Context, userID UserID) context.Context {
	return context.WithValue(ctx, UserID("userID"), userID)
}

// UserIDFromContext retrieves the user ID from the context.
func UserIDFromContext(ctx context.Context) (uuid.UUID, bool) {
	rawUserID, ok := ctx.Value(UserID("userID")).(UserID)
	if !ok {
		return uuid.UUID{}, false
	}

	userID, err := uuid.Parse(string(rawUserID))
	if err != nil {
		return uuid.UUID{}, false
	}

	return userID, true
}
