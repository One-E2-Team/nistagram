package handler

import (
	"encoding/json"
	"net/http"
	"nistagram/auth/dto"
	"nistagram/auth/service"
)

type AuthHandler struct {
	AuthService *service.AuthService
}

func (handler *AuthHandler) LogIn(w http.ResponseWriter, r *http.Request) {
	var dto dto.LogInDTO
	err := json.NewDecoder(r.Body).Decode(&dto)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	err = handler.AuthService.LogIn(dto)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	} else {
		w.WriteHeader(http.StatusOK)
	}
	w.Header().Set("Content-Type", "application/json")
}
