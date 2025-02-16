package accessToken

import (
	customErrors "Backend-trainee-assignment-winter-2025/internal/errors"
	errorResponse "Backend-trainee-assignment-winter-2025/internal/handler/middleware"
	"golang.org/x/net/context"
	"log/slog"
	"net/http"
) //TODO improve

func JwtAuthMiddleware() func(http.Handler) http.Handler {
	slog.Debug("JWT Auth Middleware")
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

			authToken, err := extractToken(r)
			if err != nil {
				slog.Error("Error extracting token", "error", err)
				errorResponse.JSONResponse(w, http.StatusUnauthorized, customErrors.ErrUnauthorized)
				return
			}

			userId, err := extractIDFromToken(authToken, string(JwtSecretKey))
			if err != nil {
				slog.Error("Error extracting user id", "error", err)
				errorResponse.JSONResponse(w, http.StatusUnauthorized, customErrors.ErrUnauthorized)
				return
			}
			ctx := context.WithValue(r.Context(), "userId", userId)
			next.ServeHTTP(w, r.WithContext(ctx))
			slog.Debug("JWT Auth Middleware", "userId", userId)
			return
		})
	}
}
