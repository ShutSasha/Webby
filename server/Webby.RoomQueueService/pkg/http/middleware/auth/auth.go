package auth

import (
	"context"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

func AuthMiddleware(jwtSecret []byte) gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		tokenStr := strings.TrimPrefix(authHeader, "Bearer ")

		token, err := jwt.Parse(tokenStr, func(token *jwt.Token) (any, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("unexpected signing method")
			}
			return jwtSecret, nil
		})

		if err != nil {
			c.AbortWithStatusJSON(
				http.StatusUnauthorized,
				map[string]string{"error": "Invalid token"},
			)
			return
		}

		if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
			if exp, ok := claims["exp"].(float64); ok {
				if time.Now().After(time.Unix(int64(exp), 0)) {
					c.AbortWithStatusJSON(
						http.StatusUnauthorized,
						map[string]string{"error": "Token expired"},
					)
					return
				}
			}

			userID, ok := claims["Id"].(string)
			if !ok || userID == "" {
				c.AbortWithStatusJSON(
					http.StatusUnauthorized,
					map[string]string{"error": "Id is missing in token"},
				)
				return
			}

			ctx := context.WithValue(c.Request.Context(), "userID", userID)
			c.Request = c.Request.WithContext(ctx)

			c.Next()
			return
		}

		c.AbortWithStatusJSON(
			http.StatusUnauthorized,
			map[string]string{"error": "Invalid token claims"},
		)
	}
}
