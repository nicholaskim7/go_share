package middleware

import (
	"context"
	"net/http"

	"github.com/nicholaskim7/go_share/internal/auth"
)

type ContextKey string

const UserIDKey ContextKey = "userID"

func AuthMiddleware(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// get the cookie
		cookie, err := r.Cookie("auth_token")
		if err != nil {
			// no cookie found
			http.Error(w, "Unauthorized: Please log in", http.StatusUnauthorized)
			return
		}
		// validate token
		userID, err := auth.ValidateToken(cookie.Value)
		if err != nil {
			// token invalid
			http.Error(w, "Unauthorized: Invalid session", http.StatusUnauthorized)
			return
		}
		// store userID in the request context
		// tells next handler who this is
		ctx := context.WithValue(r.Context(), UserIDKey, userID)
		next.ServeHTTP(w, r.WithContext(ctx))
	}
}
