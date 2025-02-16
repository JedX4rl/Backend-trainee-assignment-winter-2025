package accessToken

import (
	"Backend-trainee-assignment-winter-2025/internal/models"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestJwtAuthMiddleware_Success(t *testing.T) {
	user := &models.User{Username: "testuser", Id: 1}
	secretKey := "testsecret"
	SetSecretKey(secretKey)
	SetTokenExpiry(1)

	token, err := GenerateJWT(user, TokenExpiry)
	assert.NoError(t, err)

	req := httptest.NewRequest(http.MethodGet, "/test", nil)
	req.Header.Set("Authorization", "Bearer "+token)

	rr := httptest.NewRecorder()

	handler := JwtAuthMiddleware()(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		userId := r.Context().Value("userId")
		assert.Equal(t, 1, userId)
		w.WriteHeader(http.StatusOK)
	}))

	handler.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusOK, rr.Code)
}

func TestJwtAuthMiddleware_InvalidToken(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/test", nil)
	req.Header.Set("Authorization", "Bearer invalidtoken")

	rr := httptest.NewRecorder()

	handler := JwtAuthMiddleware()(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		t.Fatal("Request should not have passed through middleware")
	}))

	handler.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusUnauthorized, rr.Code)

	expectedError := `{"errors":"user unauthorized"}`
	assert.JSONEq(t, expectedError, rr.Body.String())
}

func TestJwtAuthMiddleware_MissingToken(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/test", nil)

	rr := httptest.NewRecorder()

	handler := JwtAuthMiddleware()(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		t.Fatal("Request should not have passed through middleware")
	}))

	handler.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusUnauthorized, rr.Code)

	expectedError := `{"errors":"user unauthorized"}`
	assert.JSONEq(t, expectedError, rr.Body.String())
}
