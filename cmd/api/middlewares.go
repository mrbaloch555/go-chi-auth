package main

import (
	"errors"
	"net/http"
	"os"
	"strings"

	"github.com/dgrijalva/jwt-go"
	"github.com/mrbaloch555/go-chi-auth/common"
)

func (app *Config) Middleware(roles ...string) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			authHeader := r.Header.Get("Authorization")
			if authHeader == "" {
				app.errorJSON(w, errors.New("not authorized"), http.StatusUnauthorized)
				return
			}

			tokenParts := strings.Split(authHeader, "Bearer ")
			if len(tokenParts) != 2 {
				app.errorJSON(w, errors.New("invalid authorization format, should be Bearer token"), http.StatusUnauthorized)
				return
			}

			tokenValue := strings.TrimSpace(tokenParts[1])

			claims := &common.Claims{}

			tkn, err := jwt.ParseWithClaims(tokenValue, claims, func(t *jwt.Token) (interface{}, error) {
				return []byte(os.Getenv("JWT_SECRET")), nil
			})
			if err != nil || !tkn.Valid || !hasRequiredRole(claims, roles) {
				app.errorJSON(w, errors.New("not authorized"), http.StatusUnauthorized)
				return
			}

			next.ServeHTTP(w, r)
		})
	}
}
func hasRequiredRole(claims *common.Claims, roles []string) bool {
	requiredRole := "user"
	for _, role := range roles {
		if role == requiredRole {
			return true
		}
	}
	return false
}
