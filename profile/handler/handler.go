package handler

import (
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"gopkg.in/go-playground/validator.v9"
	"net/http"
	"nistagram/profile/dto"
	"nistagram/profile/service"
	"nistagram/util"
)

type Handler struct {
	ProfileService *service.ProfileService
}

func (handler *Handler) Register(w http.ResponseWriter, r *http.Request) {
	var dto dto.RegistrationDto
	err := json.NewDecoder(r.Body).Decode(&dto)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

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

	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
	} else {
		w.WriteHeader(http.StatusCreated)
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte("{\"success\":\"ok\"}"))
	}
	return
}

func (handler *Handler) Search(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	result := handler.ProfileService.Search(vars["username"])
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(result)
	w.WriteHeader(http.StatusOK)
}

func (handler *Handler) GetProfileByUsername(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	result := handler.ProfileService.GetProfileByUsername(vars["username"])
	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(&result)
}

func (handler *Handler) ChangeProfileSettings(w http.ResponseWriter, r *http.Request) {
	var dto dto.ProfileSettingsDTO
	err := json.NewDecoder(r.Body).Decode(&dto)
	if err != nil {
		fmt.Println(err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	userId := util.GetLoggedUserIDFromToken(r)
	err = handler.ProfileService.ChangeProfileSettings(dto, userId)
	if err != nil {
		fmt.Println(err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("{\"success\":\"ok\"}"))
	w.Header().Set("Content-Type", "application/json")
}

func (handler *Handler) ChangePersonalData(w http.ResponseWriter, r *http.Request) {
	var dto dto.PersonalDataDTO
	err := json.NewDecoder(r.Body).Decode(&dto)
	if err != nil {
		fmt.Println(err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	userId := util.GetLoggedUserIDFromToken(r)
	err = handler.ProfileService.ChangePersonalData(dto, userId)
	if err != nil {
		fmt.Println(err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("{\"success\":\"ok\"}"))
	w.Header().Set("Content-Type", "application/json")
}

func (handler *Handler) Test(w http.ResponseWriter, r *http.Request){
	var key string
	err := json.NewDecoder(r.Body).Decode(&key)
	if err != nil {
		fmt.Println(err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	err = handler.ProfileService.Test(key)
	if err != nil {
		fmt.Println(err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	w.WriteHeader(http.StatusOK)
}

func (handler *Handler) GetAllInterests(w http.ResponseWriter, r *http.Request){
	interests, err := handler.ProfileService.GetAllInterests()
	if err != nil {
		fmt.Println(err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(interests)
}
