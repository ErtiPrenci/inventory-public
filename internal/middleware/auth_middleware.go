package middleware

import (
	"context"
	"inventory-backend/internal/utils/auth"
	"inventory-backend/internal/utils/response"
	"net/http"
	"strings"
)

func AuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			response.Unauthorized(w, "Unauthorized", "Missing Authorization header")
			return
		}

		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || parts[0] != "Bearer" {
			response.Unauthorized(w, "Unauthorized", "Invalid Authorization header format")
			return
		}

		tokenString := parts[1]

		claims, err := auth.ValidateToken(tokenString)
		if err != nil {
			response.Unauthorized(w, "Unauthorized", "Invalid or expired token")
			return
		}

		// Set user ID in context
		ctx := context.WithValue(r.Context(), "user_id", claims.UserID)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
