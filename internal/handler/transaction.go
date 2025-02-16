package handler

import (
	customErrors "Backend-trainee-assignment-winter-2025/internal/errors"
	"Backend-trainee-assignment-winter-2025/internal/handler/middleware"
	"Backend-trainee-assignment-winter-2025/internal/models"
	"encoding/json"
	"log/slog"
	"net/http"
)

func (h *Handler) SendCoin(w http.ResponseWriter, r *http.Request) {
	slog.Debug("send coin request received")
	userId := r.Context().Value("userId")
	if userId == nil {
		slog.Error("userId is nil", "function", "SendCoin")
		middleware.JSONResponse(w, http.StatusUnauthorized, "unauthorized")
		return
	}
	var sendCoin models.SendCoinRequest

	err := json.NewDecoder(r.Body).Decode(&sendCoin)
	if err != nil {
		slog.Error("sendCoin json decode error", "error", err)
		middleware.JSONResponse(w, http.StatusInternalServerError, customErrors.ErrInternal)
		return
	}

	if err = h.services.Transaction.SendMoney(r.Context(), userId.(int), sendCoin.ToUser, sendCoin.Amount); err != nil {
		slog.Error("sendCoin error", "error", err)
		middleware.JSONResponse(w, http.StatusInternalServerError, err.Error())
		return
	}
	w.WriteHeader(http.StatusOK)
	slog.Debug("send coin response sent")
}
