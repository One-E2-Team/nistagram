package handler

import (
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"gopkg.in/go-playground/validator.v9"
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

	fmt.Print(dto)

	v := validator.New()
	err = v.Struct(dto)

	if err != nil {
		for _, e := range err.(validator.ValidationErrors) {
			w.Write([]byte(e.Field()))
			w.Write([]byte(" "))
		}
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

func (handler *Handler) Search (w http.ResponseWriter, r *http.Request){
	vars := mux.Vars(r)
	result := handler.ProfileService.Search(vars["username"])
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(result)
	w.WriteHeader(http.StatusOK)
}