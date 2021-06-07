package handler

import (
	"encoding/json"
	"github.com/gorilla/mux"
	"net/http"
	"nistagram/connection/model"
	"nistagram/connection/service"
	"nistagram/util"
	"strconv"
)

type Handler struct {
	ConnectionService *service.ConnectionService
}

func (handler *Handler) AddProfile(w http.ResponseWriter, r *http.Request){
	vars := mux.Vars(r)
	id, _ := strconv.ParseUint(vars["id"],10,32)
	profile, ok := handler.ConnectionService.AddProfile(uint(id))
	if ok {
		w.WriteHeader(http.StatusOK)
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(profile)
	} else {
		w.WriteHeader(http.StatusInternalServerError)
	}
}

func (handler *Handler) GetConnection(w http.ResponseWriter, r *http.Request){
	vars := mux.Vars(r)
	id1, e1 := strconv.ParseUint(vars["followerId"],10,32)
	id2, e2 := strconv.ParseUint(vars["profileId"],10,32)
	if e1!=nil || e2!=nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	connection := handler.ConnectionService.GetConnection(uint(id1), uint(id2))
	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(connection)
}

func (handler *Handler) GetConnectionPublic(w http.ResponseWriter, r *http.Request){
	vars := mux.Vars(r)
	id1, e1 := strconv.ParseUint(vars["profileId"],10,32)
	if e1!=nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	id := util.GetLoggedUserIDFromToken(r)
	if id == 0 {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}
	connection := handler.ConnectionService.GetConnection(uint(id1), id)
	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(connection)
}

func (handler *Handler) FollowRequest(w http.ResponseWriter, r *http.Request){
	vars := mux.Vars(r)
	id2, e2 := strconv.ParseUint(vars["profileId"],10,32)
	if e2!=nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	id := util.GetLoggedUserIDFromToken(r)
	if id == 0 {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	connection, ok := handler.ConnectionService.FollowRequest(id, uint(id2))
	if ok {
		w.WriteHeader(http.StatusOK)
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(connection)
	} else {
		w.WriteHeader(http.StatusInternalServerError)
	}
}

func (handler *Handler) FollowApprove(w http.ResponseWriter, r *http.Request){
	vars := mux.Vars(r)
	id1, e1 := strconv.ParseUint(vars["profileId"],10,32)
	if e1!=nil{
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	id := util.GetLoggedUserIDFromToken(r)
	if id == 0 {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	connection, ok := handler.ConnectionService.ApproveConnection(uint(id1), id)
	if ok {
		w.WriteHeader(http.StatusOK)
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(connection)
	} else {
		w.WriteHeader(http.StatusInternalServerError)
	}
}

func (handler *Handler) GetFollowedProfilesNotMuted(w http.ResponseWriter, r *http.Request){
	vars := mux.Vars(r)
	id, _ := strconv.ParseUint(vars["id"],10,32)
	conn := model.Connection{
		PrimaryProfile:    uint(id),
		SecondaryProfile:  0,
		Muted:             false,
		CloseFriend:       false,
		NotifyPost:        false,
		NotifyStory:       false,
		NotifyMessage:     false,
		NotifyComment:     false,
		ConnectionRequest: false,
		Approved:          true,
		MessageRequest:    false,
		MessageConnected:  false,
		Block:             false,
	}
	profiles := handler.ConnectionService.GetConnectedProfiles(conn, true)
	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(*profiles)
}

func (handler *Handler) GetFollowedProfiles(w http.ResponseWriter, r *http.Request){
	vars := mux.Vars(r)
	id, _ := strconv.ParseUint(vars["id"],10,32)
	conn := model.Connection{
		PrimaryProfile:    uint(id),
		SecondaryProfile:  0,
		Muted:             false,
		CloseFriend:       false,
		NotifyPost:        false,
		NotifyStory:       false,
		NotifyMessage:     false,
		NotifyComment:     false,
		ConnectionRequest: false,
		Approved:          true,
		MessageRequest:    false,
		MessageConnected:  false,
		Block:             false,
	}
	profiles := handler.ConnectionService.GetConnectedProfiles(conn, false)
	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(*profiles)
}

func (handler *Handler) UpdateConnection(w http.ResponseWriter, r *http.Request) {
	var dto model.Connection
	err := json.NewDecoder(r.Body).Decode(&dto)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
	} else {
		id := util.GetLoggedUserIDFromToken(r)
		ret, ok := handler.ConnectionService.UpdateConnection(id, dto)
		if !ok {
			w.WriteHeader(http.StatusInternalServerError)
		} else {
			w.WriteHeader(http.StatusOK)
			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(*ret)
		}
	}
}

func (handler *Handler) DeclineFollowRequest(writer http.ResponseWriter, request *http.Request) {
	vars := mux.Vars(request)
	followerId, e1 := strconv.ParseUint(vars["followerId"],10,32)
	if e1 != nil {
		writer.WriteHeader(http.StatusInternalServerError)
	}
	id := util.GetLoggedUserIDFromToken(request)
	_, ok := handler.ConnectionService.DeleteConnection(uint(followerId), id)
	if !ok {
		writer.WriteHeader(http.StatusInternalServerError)
	} else {
		writer.WriteHeader(http.StatusOK)
		writer.Header().Set("Content-Type", "application/json")
		writer.Write([]byte("{\"status\":\"ok\"}"))
	}
}

func (handler *Handler) GetAllFollowRequests(writer http.ResponseWriter, request *http.Request) {
	
}
