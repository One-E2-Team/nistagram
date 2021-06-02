package handler

import (
	"encoding/json"
	"github.com/gorilla/mux"
	"net/http"
	"nistagram/connection/service"
	"strconv"
)

type Handler struct {
	ConnectionService *service.ConnectionService
}

func (handler *Handler) AddProfile(w http.ResponseWriter, r *http.Request){
	vars := mux.Vars(r)
	id, _ := strconv.ParseUint(vars["id"],10,32)
	profile, ok := handler.ConnectionService.AddProfile(uint(id), vars["username"])
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(profile)
	if ok {
		w.WriteHeader(http.StatusOK)
	} else {
		w.WriteHeader(http.StatusInternalServerError)
	}
}