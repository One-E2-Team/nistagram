package handler

import (
	"encoding/json"
	"github.com/gorilla/mux"
	"net/http"
	"nistagram/connection/model"
	"nistagram/connection/service"
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

func (handler *Handler) FollowRequest(w http.ResponseWriter, r *http.Request){
	vars := mux.Vars(r)
	id1, e1 := strconv.ParseUint(vars["followerId"],10,32)
	id2, e2 := strconv.ParseUint(vars["profileId"],10,32)
	if e1!=nil || e2!=nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	connection, ok := handler.ConnectionService.FollowRequest(uint(id1), uint(id2))
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