package handler

import (
	"encoding/json"
	"net/http"
	"nistagram/profile/dto"
	"nistagram/profile/service"
)

type Handler struct {
	ProfileService *service.ProfileService
}

func (handler *Handler) Register(w http.ResponseWriter, r *http.Request) {
	var dto dto.RegistrationDto
	err := json.NewDecoder(r.Body).Decode(&dto)

	if err != nil{
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	err = handler.ProfileService.Register(dto)

	if err != nil{
		w.WriteHeader(http.StatusBadRequest)
		return
	}else {
		w.WriteHeader(http.StatusCreated)
	}

	w.Header().Set("Content-Type", "application/json")
}