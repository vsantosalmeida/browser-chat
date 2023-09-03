package midleware

import (
	"context"
	"net/http"
	"strings"

	"github.com/vsantosalmeida/browser-chat/pkg/auth"
)

func AuthMiddleware(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		bearer := r.Header.Get("Authorization")
		token := strings.TrimPrefix(bearer, "Bearer ")

		if token != "" {
			user, err := auth.ValidateJWTToken(token)
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
