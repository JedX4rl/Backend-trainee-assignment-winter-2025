package handler

import (
	customErrors "Backend-trainee-assignment-winter-2025/internal/errors"
	"Backend-trainee-assignment-winter-2025/internal/handler/middleware"
	"github.com/go-chi/chi/v5"
	"log/slog"
	"net/http"
)

func (h *Handler) GetInfo(w http.ResponseWriter, r *http.Request) {

	slog.Debug("get info request received")

	userId := r.Context().Value("userId")
	if userId == nil {
		slog.Error("userId is nil")
		middleware.JSONResponse(w, http.StatusUnauthorized, "unauthorized")
		return
	}

	infoResponse, err := h.services.User.GetInfo(r.Context(), userId.(int))
	if err != nil {
		slog.Error("getInfo error", "error", err)
		middleware.JSONResponse(w, http.StatusInternalServerError, customErrors.ErrInternal)
		return
	}
	middleware.JSONResponse(w, http.StatusOK, infoResponse)
	slog.Debug("get info response sent")
}

func (h *Handler) BuyItem(w http.ResponseWriter, r *http.Request) {
	slog.Debug("buy item request received")
	userId := r.Context().Value("userId")
	if userId == nil {
		slog.Error("userId is nil", "function", "BuyItem")
		middleware.JSONResponse(w, http.StatusUnauthorized, customErrors.ErrUnauthorized)
		return
	}
	item := chi.URLParam(r, "item")
	if item == "" {
		slog.Error("urlParam is empty", "function", "BuyItem")
		middleware.JSONResponse(w, http.StatusBadRequest, customErrors.ErrBindFailed)
		return
	}

	if err := h.services.Shop.BuyItem(r.Context(), userId.(int), item); err != nil {
		slog.Error("buy item error", "error", err)
		middleware.JSONResponse(w, http.StatusInternalServerError, err.Error())
		return
	}

	w.WriteHeader(http.StatusOK)
	slog.Debug("buy item response sent")
}
