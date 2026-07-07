package httpapi

import (
	"context"
	"net/http"
	"strings"

	"github.com/WassimBannout/gatherops/internal/security"
	"github.com/google/uuid"
)

type authContextKey struct{}

func AuthMiddleware(tokens security.TokenManager) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			authorization := r.Header.Get("Authorization")
			if authorization == "" {
				WriteError(w, r, http.StatusUnauthorized, "unauthorized", "Authentication required", nil)
				return
			}

			parts := strings.SplitN(authorization, " ", 2)
			if len(parts) != 2 || !strings.EqualFold(parts[0], "Bearer") || strings.TrimSpace(parts[1]) == "" {
				WriteError(w, r, http.StatusUnauthorized, "unauthorized", "Authentication required", nil)
				return
			}

			userID, err := tokens.ParseAccessToken(strings.TrimSpace(parts[1]))
			if err != nil {
				WriteError(w, r, http.StatusUnauthorized, "invalid_access_token", "Access token is invalid or expired", nil)
				return
			}

			ctx := context.WithValue(r.Context(), authContextKey{}, userID)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

func AuthenticatedUserID(ctx context.Context) (uuid.UUID, bool) {
	userID, ok := ctx.Value(authContextKey{}).(uuid.UUID)
	return userID, ok && userID != uuid.Nil
}
