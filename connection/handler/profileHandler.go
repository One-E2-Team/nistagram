package handler

import (
	"encoding/json"
	"github.com/gorilla/mux"
	"net/http"
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
