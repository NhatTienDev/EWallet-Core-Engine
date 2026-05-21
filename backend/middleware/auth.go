package middleware

import (
	"fmt"
	"net/http"
	"context"
	"encoding/json"
	"strings"

	"github.com/golang-jwt/jwt/v5"
)

// Create a custom type to avoid key collision in context
type contextKey string
const UserIDKey contextKey = "user_id"

func AuthMiddleware(jwtSecretKey []byte) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// Get header authorization
			authHeader := r.Header.Get("Authorization")
			if authHeader == "" {
				writeJSONWithError(w, http.StatusUnauthorized, "Missing Authorization header")
				return
			}

			// Seperate "Bearer" from the token string
			parts := strings.Split(authHeader, " ")
			if len(parts) != 2 || parts[0] != "Bearer" {
				writeJSONWithError(w, http.StatusUnauthorized, "Invalid Authorization header format")
				return
			}

			// Parse and verify the token signature
			tokenString := parts[1]
			token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
				if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
					return nil, fmt.Errorf("Unexpected signing method: %v", token.Header["alg"])
				}
				return jwtSecretKey, nil
			})

			if err != nil || !token.Valid {
				writeJSONWithError(w, http.StatusUnauthorized, "Invalid or expired token")
				return
			}

			// Extract user_id from payload
			claims, ok := token.Claims.(jwt.MapClaims)
			if !ok || claims["user_id"] == nil {
				writeJSONWithError(w, http.StatusUnauthorized, "Invalid token claims")
				return
			}

			userIDFloat, ok := claims["user_id"].(float64)
			if !ok {
				writeJSONWithError(w, http.StatusUnauthorized, "User ID not found in the token")
				return
			}
			userID := int64(userIDFloat)

			// Insert userID into context and allow the request to proceed
			ctx := context.WithValue(r.Context(), UserIDKey, userID)
			reqWithCtx := r.WithContext(ctx)

			next.ServeHTTP(w, reqWithCtx)
		})
	}
}

func writeJSONWithError(w http.ResponseWriter, status int, message string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(map[string]string{"error": message})
}