// Package authmiddleware provides middleware for authentication.
package authmiddleware

import (
	"context"
	"net/http"

	contexthelper "github.com/estructy/server/internal/helpers/context"
	"github.com/estructy/server/internal/infra/database/repository"
	"github.com/google/uuid"
	"github.com/lestrrat-go/jwx/v3/jwk"
	"github.com/lestrrat-go/jwx/v3/jwt"
)

type AuthMiddleware struct {
	repo *repository.Queries
}

func NewAuthMiddleware(repo *repository.Queries) *AuthMiddleware {
	return &AuthMiddleware{
		repo: repo,
	}
}

func (m *AuthMiddleware) Handle(next http.Handler) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		jwkSet, err := jwk.Fetch(r.Context(), "http://localhost:4000/api/auth/jwks")
		if err != nil {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		token, err := jwt.ParseRequest(r, jwt.WithKeySet(jwkSet))
		if err != nil || token == nil {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		authID, name, email, err := getClaimsFromToken(token)
		if err != nil {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		repoCtx := context.Background()

		userID, err := m.repo.GetUserByAuthID(repoCtx, authID)
		if err != nil && err.Error() == "no rows in result set" {
			userID, err = createNewUser(m.repo, authID, name, email)
		}
		if err != nil {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		ctx := r.Context()
		userCtx := contexthelper.WithUserID(ctx, contexthelper.UserID(userID.String()))
		r = r.WithContext(userCtx)

		next.ServeHTTP(w, r)
	})
}

func getClaimsFromToken(token jwt.Token) (string, string, string, error) {
	var authID, name, email string

	err := token.Get("id", &authID)
	if err != nil {
		return "", "", "", err
	}

	err = token.Get("name", &name)
	if err != nil {
		return "", "", "", err
	}

	err = token.Get("email", &email)
	if err != nil {
		return "", "", "", err
	}

	return authID, name, email, nil
}

func createNewUser(repo *repository.Queries, authID, name, email string) (uuid.UUID, error) {
	newUUID, err := uuid.NewV7()
	if err != nil {
		return uuid.Nil, err
	}

	userID, err := repo.CreateUser(context.Background(), repository.CreateUserParams{
		UserID: newUUID,
		AuthID: authID,
		Name:   name,
		Email:  email,
	})
	if err != nil {
		return uuid.Nil, err
	}

	return userID, nil
}
