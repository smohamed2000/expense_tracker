package middleware

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"strings"
	"time"

	jwt "github.com/golang-jwt/jwt/v4"
)

// Define a custom type for context keys
type contextKey string

const userIDKey contextKey = "userID"

// GenerateJWT generates a JWT token for a given userID
func GenerateJWT(userID uint) (string, error) {
	secret := os.Getenv("JWT_SECRET")
	claims := jwt.MapClaims{
		"user_id": userID,
		"exp":     time.Now().Add(24 * time.Hour).Unix(),
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(secret))
}

// VerifyJWT verifies the JWT token string
func VerifyJWT(tokenString string) (*jwt.Token, error) {
	secret := os.Getenv("JWT_SECRET")
	return jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(secret), nil
	})
}

// AuthMiddleware validates JWT and adds userID to request context
func AuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			http.Error(w, "authorization header required", http.StatusUnauthorized)
			return
		}

		tokenString := strings.TrimPrefix(authHeader, "Bearer ")
		token, err := VerifyJWT(tokenString)
		if err != nil || !token.Valid {
			http.Error(w, "invalid or expired token", http.StatusUnauthorized)
			return
		}

		claims, ok := token.Claims.(jwt.MapClaims)
		if !ok {
			http.Error(w, "invalid token claims", http.StatusUnauthorized)
			return
		}

		userIDFloat, ok := claims["user_id"].(float64)
		if !ok {
			http.Error(w, "user id not found in token", http.StatusUnauthorized)
			return
		}

		// Use custom type for context key
		ctx := context.WithValue(r.Context(), userIDKey, uint(userIDFloat))
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

// GetUserIDFromToken retrieves the userID from request context
func GetUserIDFromToken(r *http.Request) (uint, error) {
	userID, ok := r.Context().Value(userIDKey).(uint)
	if !ok {
		return 0, fmt.Errorf("user id not found in context")
	}
	return userID, nil
}
