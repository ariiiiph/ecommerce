package middleware

import (
	"net/http"
)

func RequireRole(requiredRole string) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			claims, ok := GetClaims(r.Context())

			if !ok {
				http.Error(
					w,
					"unauthorized",
					http.StatusUnauthorized,
				)
				return
			}

			if claims.Role != requiredRole {
				http.Error(
					w,
					"forbidden",
					http.StatusForbidden,
				)
				return
			}

			next.ServeHTTP(w, r)
		})
	}
}
