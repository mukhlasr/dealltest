package httphandler

import (
	"net/http"
	"strings"
)

func MiddlewareAllowRoles(tokenParser func(string) (JWTClaims, error), roles ...string) func(http.Handler) http.HandlerFunc {
	return func(h http.Handler) http.HandlerFunc {
		return func(w http.ResponseWriter, r *http.Request) {
			authHeader := r.Header.Get("Authorization")
			if authHeader == "" || !strings.HasPrefix(authHeader, "Bearer ") {
				HandleUnauthorized()(w, r)
				return
			}

			jwtToken := strings.TrimPrefix(authHeader, "Bearer ")
			claims, err := tokenParser(jwtToken)
			if err != nil {
				HandleUnauthorized()(w, r)
				return
			}

			for _, role := range roles {
				if role == claims.Role {
					h.ServeHTTP(w, r)
					return
				}
			}

			HandleForbidden()(w, r)
		}
	}
}
