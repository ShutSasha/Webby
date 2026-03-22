package auth

import (
	"context"
	"fmt"
	"net/http"
	"strings"
	"time"

	"webby/pkg/http/render"

	"github.com/golang-jwt/jwt/v5"
)

func AuthMiddleware(jwtSecret []byte) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			authHeader := r.Header.Get("Authorization")

			tokenStr := strings.TrimPrefix(authHeader, "Bearer ")

			token, err := jwt.Parse(tokenStr, func(token *jwt.Token) (any, error) {
				if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
					return nil, fmt.Errorf("unexpected signing method")
				}
				return jwtSecret, nil
			})
			if err != nil {
				render.Encode(w, r, http.StatusUnauthorized, map[string]string{"error": "Invalid token"})
				return
			}

			if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
				if exp, ok := claims["exp"].(float64); ok {
					if time.Now().After(time.Unix(int64(exp), 0)) {
						render.Encode(w, r, http.StatusUnauthorized, map[string]string{"error": "Token expired"})
						return
					}
				}

				userID, ok := claims["Id"].(string)
				if !ok || userID == "" {
					render.Encode(w, r, http.StatusUnauthorized, map[string]string{"error": "Id is missing in token"})
					return
				}

				ctx := context.WithValue(r.Context(), "userID", userID)
				next.ServeHTTP(w, r.WithContext(ctx))
				return
			}

			render.Encode(w, r, http.StatusUnauthorized, map[string]string{"error": "Invalid token claims"})
		})
	}
}
