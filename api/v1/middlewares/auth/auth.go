// Package authmiddleware provides middleware for authentication.
package authmiddleware

import (
	"net/http"

	contexthelper "github.com/nahtann/controlriver.com/internal/helpers/context"
)

type AuthMiddleware struct{}

func NewAuthMiddleware() *AuthMiddleware {
	return &AuthMiddleware{}
}

func (m *AuthMiddleware) Handle(next http.Handler) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// @todo: replace with real authentication logic
		// Extract user ID from headers
		userID := r.Header.Get("X-User-ID")
		if userID == "" {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		ctx := r.Context()
		userCtx := contexthelper.WithUserID(ctx, contexthelper.UserID(userID))
		r = r.WithContext(userCtx)

		next.ServeHTTP(w, r)
	})
}
