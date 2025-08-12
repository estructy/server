// Package contexthelper provides utility functions to manage user IDs in a context.
package contexthelper

import (
	"context"

	"github.com/google/uuid"
)

type (
	UserID    string
	AccountID string
)

// WithUserID adds a user ID to the context.
func WithUserID(ctx context.Context, userID UserID) context.Context {
	return context.WithValue(ctx, UserID("userID"), userID)
}

func WithAccountID(ctx context.Context, accountID AccountID) context.Context {
	return context.WithValue(ctx, AccountID("accountID"), accountID)
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

// AccountIDFromContext retrieves the account ID from the context.
func AccountIDFromContext(ctx context.Context) (uuid.UUID, bool) {
	rawAccountID, ok := ctx.Value(AccountID("accountID")).(AccountID)
	if !ok {
		return uuid.UUID{}, false
	}

	accountID, err := uuid.Parse(string(rawAccountID))
	if err != nil {
		return uuid.UUID{}, false
	}

	return accountID, true
}
