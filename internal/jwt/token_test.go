package accessToken

import (
	"Backend-trainee-assignment-winter-2025/internal/models"
	"github.com/golang-jwt/jwt/v4"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

func TestSetSecretKey(t *testing.T) {
	secret := "mysecret"
	err := SetSecretKey(secret)

	assert.Nil(t, err)

	assert.Equal(t, []byte(secret), JwtSecretKey)

	err = SetSecretKey("")
	assert.EqualError(t, err, "missing secret key")
}

func TestSetTokenExpiry(t *testing.T) {
	expiry := 24
	err := SetTokenExpiry(expiry)

	assert.Nil(t, err)

	assert.Equal(t, expiry, TokenExpiry)

	err = SetTokenExpiry(0)
	assert.EqualError(t, err, "missing token expiry")
}

func TestGenerateJWT(t *testing.T) {
	user := &models.User{
		Username: "testuser",
		Id:       1,
	}

	_ = SetSecretKey("mysecret")
	_ = SetTokenExpiry(24)

	token, err := GenerateJWT(user, TokenExpiry)

	assert.Nil(t, err)

	assert.NotEmpty(t, token)

	parsedToken, err := jwt.Parse(token, func(token *jwt.Token) (interface{}, error) {
		return []byte("mysecret"), nil
	})

	assert.Nil(t, err)

	claims, ok := parsedToken.Claims.(jwt.MapClaims)
	assert.True(t, ok)

	assert.Equal(t, "1", claims["id"].(string))
}

func TestExtractToken(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/test", nil)
	req.Header.Set("Authorization", "Bearer mytoken")

	token, err := extractToken(req)

	assert.Nil(t, err)
	assert.Equal(t, "mytoken", token)

	req = httptest.NewRequest(http.MethodGet, "/test", nil)
	token, err = extractToken(req)
	assert.Equal(t, "", token)
	assert.EqualError(t, err, "missing Authorization header")

	req.Header.Set("Authorization", "BearerToken mytoken")
	token, err = extractToken(req)
	assert.Equal(t, "", token)
	assert.EqualError(t, err, "invalid Authorization format")
}

func TestExtractIDFromToken(t *testing.T) {
	_ = SetSecretKey("mysecret")

	user := &models.User{
		Username: "testuser",
		Id:       1,
	}
	token, err := GenerateJWT(user, 24)
	assert.Nil(t, err)

	id, err := extractIDFromToken(token, "mysecret")

	assert.Nil(t, err)
	assert.Equal(t, 1, id)

	invalidToken := "invalid_token"
	_, err = extractIDFromToken(invalidToken, "mysecret")
	assert.EqualError(t, err, "token contains an invalid number of segments")

	// Тестируем случай с токеном, который не содержит ID
	claims := jwt.MapClaims{
		"username": "testuser",
		"exp":      time.Now().Add(time.Hour * 24).Unix(),
	}
	tokenWithoutID := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	signedToken, err := tokenWithoutID.SignedString([]byte("mysecret"))
	assert.Nil(t, err)

	_, err = extractIDFromToken(signedToken, "mysecret")
	assert.EqualError(t, err, "ID not found in token")
}
