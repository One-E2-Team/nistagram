package handler

import (
	"encoding/json"
	"fmt"
	"net/http"
	"nistagram/agent/service"
	"nistagram/agent/dto"
)

type AuthHandler struct {
	AuthService *service.AuthService
}

func (handler *AuthHandler) Register(w http.ResponseWriter, r *http.Request) {
	var registerDTO dto.RegisterDTO
	err := json.NewDecoder(r.Body).Decode(&registerDTO)
	if err != nil {
		fmt.Println(err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	err = handler.AuthService.Register(registerDTO)
	if err != nil {
		fmt.Println(err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	w.WriteHeader(http.StatusCreated)
	_, _ = w.Write([]byte("{\"success\":\"ok\"}"))
	w.Header().Set("Content-Type", "application/json")
}
