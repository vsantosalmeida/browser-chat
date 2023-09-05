package midleware

import (
	"context"
	"github.com/vsantosalmeida/browser-chat/pkg/auth"
	"net/http"
)

// AuthMiddleware middleware to validate an AuthenticatedUser and pass through context.
func AuthMiddleware(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		token, tok := r.URL.Query()["bearer"]
		if tok && len(token) == 1 {
			user, err := auth.ValidateJWTToken(token[0])
			if err != nil {
				http.Error(w, "forbidden", http.StatusForbidden)

			} else {
				ctx := context.WithValue(r.Context(), auth.UserContextKey, user)
				next.ServeHTTP(w, r.WithContext(ctx))
			}

		} else {
			w.WriteHeader(http.StatusUnauthorized)
			w.Write([]byte("authentication required"))
		}
	}
}
