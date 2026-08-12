package middleware

import (
	"context"
	"net/http"
	"strings"

	"github.com/ariiiiph/ecommerce/internal/config"
	"github.com/ariiiiph/ecommerce/internal/utils"
)

type contextKey string

const claimsKey contextKey = "claims"

func AuthMiddleware(
	jwtConfig config.JWTConfig,
) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			authHeader := r.Header.Get("Authorization")

			if authHeader == "" {
				http.Error(
					w,
					"authorization header is required",
					http.StatusUnauthorized,
				)
				return
			}

			parts := strings.SplitN(authHeader, " ", 2)

			if len(parts) != 2 ||
				parts[0] != "Bearer" ||
				parts[1] == "" {
				http.Error(
					w,
					"invalid authorization header",
					http.StatusUnauthorized,
				)
				return
			}

			claims, err := utils.ValidateAccessToken(
				jwtConfig,
				parts[1],
			)
			if err != nil {
				http.Error(
					w,
					"invalid or expired access token",
					http.StatusUnauthorized,
				)
				return
			}

			ctx := context.WithValue(
				r.Context(),
				claimsKey,
				claims,
			)

			next.ServeHTTP(
				w,
				r.WithContext(ctx),
			)
		})
	}
}

func GetClaims(ctx context.Context) (*utils.Claims, bool) {
	claims, ok := ctx.Value(claimsKey).(*utils.Claims)

	return claims, ok
}
