package middleware

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestJSONResponse(t *testing.T) {
	responseData := map[string]string{"message": "Hello, World!"}

	rr := httptest.NewRecorder()

	JSONResponse(rr, http.StatusOK, responseData)

	assert.Equal(t, http.StatusOK, rr.Code)

	assert.Equal(t, "application/json", rr.Header().Get("Content-Type"))

	expectedResponse := `{"message":"Hello, World!"}`
	assert.JSONEq(t, expectedResponse, rr.Body.String())
}
