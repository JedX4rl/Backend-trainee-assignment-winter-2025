package accessToken

import (
	"Backend-trainee-assignment-winter-2025/internal/config/serverConfig"
	"fmt"
	"github.com/golang-jwt/jwt/v4"
	"golang.org/x/net/context"
	"net/http"
	"os"
	"strconv"
	"strings"
) //TODO improve

type Jwt struct {
	Username string `json:"username"`
	Id       string `json:"id"`
	jwt.StandardClaims
}

func JwtAuthMiddleware() func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			secret := os.Getenv(serverConfig.ACCESS_TOKEN_SECRET)
			if secret == "" {
				http.Error(w, "Missing secret key", http.StatusInternalServerError)
				return
			}

			authToken, err := extractToken(r)
			if err != nil {
				http.Error(w, err.Error(), http.StatusUnauthorized)
				return
			}

			userId, err := extractIDFromToken(authToken, secret)
			if err != nil {
				http.Error(w, "Invalid token", http.StatusUnauthorized)
				return
			}
			ctx := context.WithValue(r.Context(), "userId", userId)
			next.ServeHTTP(w, r.WithContext(ctx))
			return
		})
	}
}

func extractToken(r *http.Request) (string, error) {
	authHeader := r.Header.Get("Authorization")
	if authHeader == "" {
		return "", fmt.Errorf("missing Authorization header")
	}
	t := strings.Split(authHeader, " ")
	if len(t) != 2 || t[0] != "Bearer" {
		return "", fmt.Errorf("invalid Authorization format")
	}
	return t[1], nil
}

func extractIDFromToken(requestToken string, secret string) (int, error) {
	token, err := jwt.Parse(requestToken, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(secret), nil
	})

	if err != nil {
		return 0, err
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok || !token.Valid {
		return 0, fmt.Errorf("invalid Token")
	}

	idStr, ok := claims["id"].(string)
	if !ok {
		return 0, fmt.Errorf("ID not found in token")
	}

	id, err := strconv.Atoi(idStr)
	if err != nil {
		return 0, fmt.Errorf("invalid ID format")
	}

	return id, nil
}
