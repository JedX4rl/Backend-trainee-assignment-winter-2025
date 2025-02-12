package handler

import (
	"Backend-trainee-assignment-winter-2025/internal/config/serverConfig"
	accessTkn "Backend-trainee-assignment-winter-2025/internal/jwt"
	"Backend-trainee-assignment-winter-2025/internal/models"
	"encoding/json"
	"golang.org/x/crypto/bcrypt"
	"log"
	"net/http"
	"os"
	"strconv"
)

func (h *Handler) Auth(w http.ResponseWriter, r *http.Request) {
	var input models.User

	err := json.NewDecoder(r.Body).Decode(&input) //TODO add validations and errors
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	//TODO add user validations

	user, err := h.services.Authorization.GetUserByUsername(r.Context(), input.Username)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	if user != nil {
		if bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(input.Password)) != nil {
			http.Error(w, "Invalid password", http.StatusUnauthorized)
			return
		}
	} else {
		encryptedPassword, err := bcrypt.GenerateFromPassword(
			[]byte(input.Password),
			bcrypt.DefaultCost,
		)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		input.Password = string(encryptedPassword)
		if err = h.services.Authorization.SignUp(r.Context(), &input); err != nil {
			http.Error(w, "Failed to create user", http.StatusInternalServerError)
			return
		}
	}

	expiry, err := strconv.Atoi(os.Getenv(serverConfig.ACCESS_TOKEN_EXPIRY_HOUR))
	if err != nil || expiry <= 0 {
		log.Println("Warning: ACCESS_TOKEN_EXPIRY_HOUR is not set or invalid")
		http.Error(w, "Server error", http.StatusInternalServerError)
		return
	}

	accessToken, err := h.services.Authorization.SignIn(&input, string(accessTkn.JwtSecretKey), expiry)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	err = json.NewEncoder(w).Encode(accessToken)
	if err != nil {
		//TODO add logs
	}
	return //TODO check if necessary
}
