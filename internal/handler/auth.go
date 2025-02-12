package handler

import (
	"Backend-trainee-assignment-winter-2025/internal/config/serverConfig"
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

	encryptedPassword, err := bcrypt.GenerateFromPassword(
		[]byte(input.Password),
		bcrypt.DefaultCost,
	)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		//http.Error(w, Errors.JsonError(err.Error()), http.StatusInternalServerError)
		return
	}

	user := models.User{
		Username: input.Username,
		Password: string(encryptedPassword),
	}

	if err = h.services.Authorization.CheckIfUserExists(r.Context(), user.Username); err != nil { //TODO check if user does not exist
		err = h.services.Authorization.SignUp(r.Context(), &user)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
	}

	if bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(input.Password)) != nil {
		http.Error(w, "Wrong password", http.StatusBadRequest) //TODO change
		return
	}

	token := os.Getenv(serverConfig.ACCESS_TOKEN_SECRET)
	if token == "" {
		log.Fatalf(serverConfig.ACCESS_TOKEN_SECRET + " environment variable not set") //TODO change
	}

	expiry, _ := strconv.ParseInt(os.Getenv(serverConfig.ACCESS_TOKEN_EXPIRY_HOUR), 10, 64)
	if expiry == 0 {
		log.Fatalf(serverConfig.ACCESS_TOKEN_EXPIRY_HOUR + " environment variable not set")
	}

	accessToken, err := h.services.Authorization.SignIn(&user, token, int(expiry))
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
