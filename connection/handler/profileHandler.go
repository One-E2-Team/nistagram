package handler

import (
	"encoding/json"
	"github.com/gorilla/mux"
	"net/http"
	"nistagram/connection/dto"
	"nistagram/connection/model"
	"nistagram/util"
	"strconv"
)

func (handler *Handler) AddProfile(w http.ResponseWriter, r *http.Request) {
	method := "nistagram/connection/handler.AddProfile"
	vars := mux.Vars(r)
	id, _ := strconv.ParseUint(vars["id"], 10, 32)
	profile, ok := handler.ConnectionService.AddOrUpdateProfile(model.ProfileVertex{ProfileID: uint(id)})
	if ok {
		w.WriteHeader(http.StatusOK)
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(profile)
		util.Logging(util.INFO, method, "", "Added user: "+util.Uint2String(uint(id)), "connection")
	} else {
		w.WriteHeader(http.StatusInternalServerError)
		util.Logging(util.INFO, method, "", "Failed to add user: "+util.Uint2String(uint(id)), "connection")
	}
}

func (handler *Handler) DeleteProfile(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, _ := strconv.ParseUint(vars["id"], 10, 32)
	profile, ok := handler.ConnectionService.AddOrUpdateProfile(model.ProfileVertex{ProfileID: uint(id), Deleted: true})
	if ok {
		w.WriteHeader(http.StatusOK)
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(profile)
	} else {
		w.WriteHeader(http.StatusInternalServerError)
	}
}

func (handler *Handler) ReActivateProfile(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, _ := strconv.ParseUint(vars["id"], 10, 32)
	profile, ok := handler.ConnectionService.AddOrUpdateProfile(model.ProfileVertex{ProfileID: uint(id), Deleted: false})
	if ok {
		w.WriteHeader(http.StatusOK)
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(profile)
	} else {
		w.WriteHeader(http.StatusInternalServerError)
	}
}


func (handler *Handler) RecommendProfiles(writer http.ResponseWriter, request *http.Request) {
	profileId := util.GetLoggedUserIDFromToken(request)
	if profileId == 0 {
		writer.Write([]byte("[{\"status\":\"error\"}]"))
		writer.WriteHeader(http.StatusBadRequest)
		return
	}

	var users *[]dto.ProfileRecommendationDTO
	users, ok := handler.ConnectionService.GetRecommendations(profileId)

	if !ok {
		writer.Write([]byte("[{\"status\":\"error\"}]"))
		writer.WriteHeader(http.StatusInternalServerError)
		return
	}

	writer.WriteHeader(http.StatusOK)
	writer.Header().Set("Content-Type", "application/json")
	json.NewEncoder(writer).Encode(users)
}