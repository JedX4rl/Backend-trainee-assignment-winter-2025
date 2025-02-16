package accessToken

import (
	"Backend-trainee-assignment-winter-2025/internal/models"
	"fmt"
	"github.com/golang-jwt/jwt/v4"
	"log/slog"
	"net/http"
	"strconv"
	"strings"
	"time"
)

type Jwt struct {
	Username string `json:"username"`
	Id       string `json:"id"`
	jwt.StandardClaims
}

var (
	JwtSecretKey []byte
	TokenExpiry  int
)

func SetSecretKey(secretKey string) error {
	if secretKey == "" {
		return fmt.Errorf("missing secret key")
	}
	JwtSecretKey = []byte(secretKey)
	return nil
}

func SetTokenExpiry(expiry int) error {
	if expiry <= 0 {
		return fmt.Errorf("missing token expiry")
	}
	TokenExpiry = expiry
	return nil
}

func GenerateJWT(user *models.User, expiry int) (string, error) {
	slog.Debug("generate JWT started successfully")
	exp := time.Now().Add(time.Hour * time.Duration(expiry)).Unix()
	claims := &Jwt{
		Username: user.Username,
		Id:       strconv.Itoa(user.Id),
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: exp,
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	t, err := token.SignedString(JwtSecretKey)
	if err != nil {
		slog.Error("Failed to generate JWT", "error", err)
		return "", fmt.Errorf("failed to generate token")
	}
	slog.Debug("generate JWT finished successfully")
	return t, nil
}

func extractToken(r *http.Request) (string, error) {
	slog.Debug("Extracting token")
	authHeader := r.Header.Get("Authorization")
	if authHeader == "" {
		slog.Error("Missing Authorization header")
		return "", fmt.Errorf("missing Authorization header")
	}
	t := strings.Split(authHeader, " ")
	if len(t) != 2 || t[0] != "Bearer" {
		slog.Error("Invalid Authorization header", "authHeader", authHeader)
		return "", fmt.Errorf("invalid Authorization format")
	}
	slog.Debug("Extracting token finished successfully")
	return t[1], nil
}

func extractIDFromToken(requestToken string, secret string) (int, error) {
	slog.Debug("Extracting id from token")
	token, err := jwt.Parse(requestToken, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			slog.Error("Unexpected signing method")
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		slog.Debug("Extracting id from token finished successfully")
		return []byte(secret), nil
	})

	if err != nil {
		slog.Error("Error extracting id from token", "error", err)
		return 0, err
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok || !token.Valid {
		slog.Error("Error extracting id from token", "error", err)
		return 0, fmt.Errorf("invalid Token")
	}

	idStr, ok := claims["id"].(string)
	if !ok {
		slog.Error("Error extracting id from token", "error", err)
		return 0, fmt.Errorf("ID not found in token")
	}

	id, err := strconv.Atoi(idStr)
	if err != nil {
		slog.Error("Error extracting id from token", "error", err)
		return 0, fmt.Errorf("invalid ID format")
	}

	slog.Debug("Extracting id finished successfully")
	return id, nil
}
