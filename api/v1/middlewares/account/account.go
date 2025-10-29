// Package accountmiddleware provides middleware for account-related operations.
package accountmiddleware

import (
	"net/http"

	contexthelper "github.com/estructy/server/internal/helpers/context"
)

type AccountMiddleware struct{}

func NewAccountMiddleware() *AccountMiddleware {
	return &AccountMiddleware{}
}

func (m *AccountMiddleware) Handle(next http.Handler) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// @todo: Implement cache lookup for account ID. Consider using memcached for this operation.
		// Extract account ID from headers
		accountID := r.Header.Get("X-Account-ID")
		if accountID == "" {
			http.Error(w, "Account ID is required", http.StatusBadRequest)
			return
		}

		ctx := r.Context()
		accountCtx := contexthelper.WithAccountID(ctx, contexthelper.AccountID(accountID))
		r = r.WithContext(accountCtx)

		next.ServeHTTP(w, r)
	})
}
