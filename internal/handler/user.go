package handler

import (
	"encoding/json"
	"net/http"
)

func (h *Handler) Info(w http.ResponseWriter, r *http.Request) {
	tempUserId := r.Context().Value("userId")
	if tempUserId == "" {
		http.Error(w, "userId is empty", http.StatusUnauthorized)
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode("hola!")
}
