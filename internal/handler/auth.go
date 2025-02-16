package handler

import (
	customErrors "Backend-trainee-assignment-winter-2025/internal/errors"
	"Backend-trainee-assignment-winter-2025/internal/handler/middleware"
	"Backend-trainee-assignment-winter-2025/internal/models"
	structValidator "Backend-trainee-assignment-winter-2025/internal/validator"
	"encoding/json"
	"log/slog"
	"net/http"
)

func (h *Handler) Auth(w http.ResponseWriter, r *http.Request) {

	slog.Debug("auth handler request received")

	var authRequest models.AuthRequest

	err := json.NewDecoder(r.Body).Decode(&authRequest)
	if err != nil {
		slog.Error("error decoding body", "error", err)
		middleware.JSONResponse(w, http.StatusBadRequest, customErrors.ErrBindFailed)
		return
	}
	if err = structValidator.ValidateStruct(authRequest); err != nil {
		slog.Error("error validating struct", "error", err)
		middleware.JSONResponse(w, http.StatusBadRequest, customErrors.ErrValidate)
		return
	}

	accessToken, err := h.services.User.Auth(r.Context(), authRequest.Username, authRequest.Password)
	if err != nil {
		slog.Error("error authenticating user", "error", err)
		middleware.JSONResponse(w, http.StatusUnauthorized, err)
		return
	}
	middleware.JSONResponse(w, http.StatusOK, accessToken)
	slog.Debug("auth response sent")
}
