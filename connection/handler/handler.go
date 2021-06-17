package handler

import (
	"encoding/json"
	"net/http"
	"nistagram/connection/dto"
	"nistagram/connection/model"
	"nistagram/connection/service"
	"nistagram/util"
	"strconv"

	"github.com/gorilla/mux"
)

type Handler struct {
	ConnectionService *service.ConnectionService
}

func (handler *Handler) AddProfile(w http.ResponseWriter, r *http.Request) {
	method := "nistagram/connection/handler.AddProfile"
	vars := mux.Vars(r)
	id, _ := strconv.ParseUint(vars["id"], 10, 32)
	profile, ok := handler.ConnectionService.AddProfile(uint(id))
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

func (handler *Handler) GetConnection(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id1, e1 := strconv.ParseUint(vars["followerId"], 10, 32)
	id2, e2 := strconv.ParseUint(vars["profileId"], 10, 32)
	if e1 != nil || e2 != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	connection := handler.ConnectionService.GetConnection(uint(id1), uint(id2))
	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(connection)
}

func (handler *Handler) GetConnectionPublic(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id1, e1 := strconv.ParseUint(vars["profileId"], 10, 32)
	if e1 != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	id := util.GetLoggedUserIDFromToken(r)
	if id == 0 {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}
	connection := handler.ConnectionService.GetConnection(id, uint(id1))
	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(connection)
}

func (handler *Handler) FollowRequest(w http.ResponseWriter, r *http.Request) {
	method := "nistagram/connection/handler.FollowRequest"
	vars := mux.Vars(r)
	id2, e2 := strconv.ParseUint(vars["profileId"], 10, 32)
	if e2 != nil {
		w.WriteHeader(http.StatusBadRequest)
		util.Logging(util.ERROR, method, "", "Failed Follow request (possible auto-approve): "+util.GetLoggingStringFromID(util.GetLoggedUserIDFromToken(r))+"->"+util.Uint2String(uint(id2)), "connection")
		return
	}
	id := util.GetLoggedUserIDFromToken(r)
	if id == 0 {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}
	connection, ok := handler.ConnectionService.FollowRequest(id, uint(id2))
	if ok && connection != nil {
		w.WriteHeader(http.StatusOK)
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(connection)
		util.Logging(util.INFO, method, "", "Follow request (possible auto-approve): "+util.Uint2String(id)+"->"+util.Uint2String(uint(id2)), "connection")
	} else {
		w.WriteHeader(http.StatusInternalServerError)
		util.Logging(util.ERROR, method, "", "Service call failed:"+util.GetLoggingStringFromID(util.GetLoggedUserIDFromToken(r))+"->"+util.Uint2String(uint(id2)), "connection")
	}
}

func (handler *Handler) FollowApprove(w http.ResponseWriter, r *http.Request) {
	method := "nistagram/connection/handler.FollowApprove"
	vars := mux.Vars(r)
	id1, e1 := strconv.ParseUint(vars["profileId"], 10, 32)
	if e1 != nil {
		w.WriteHeader(http.StatusBadRequest)
		util.Logging(util.ERROR, method, "", e1.Error(), "connection")
		return
	}
	id := util.GetLoggedUserIDFromToken(r)
	if id == 0 {
		w.WriteHeader(http.StatusUnauthorized)
		util.Logging(util.ERROR, method, "", "Unauthorized", "connection")
		return
	}
	_, ok := handler.ConnectionService.ApproveConnection(uint(id1), id)
	if ok {
		w.WriteHeader(http.StatusOK)
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte("{\"status\":\"ok\"}"))
		util.Logging(util.INFO, method, "", "Approved follow request: "+util.Uint2String(id)+"->"+util.Uint2String(uint(id1)), "connection")
	} else {
		util.Logging(util.ERROR, method, "", "Service Error", "connection")
		w.WriteHeader(http.StatusInternalServerError)
	}
}

func (handler *Handler) GetFollowedProfilesNotMuted(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, _ := strconv.ParseUint(vars["id"], 10, 32)
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

func (handler *Handler) GetFollowedProfiles(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, _ := strconv.ParseUint(vars["id"], 10, 32)
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
	method := "nistagram/connection/handler.DeclineFollowRequest"
	vars := mux.Vars(request)
	followerId, e1 := strconv.ParseUint(vars["profileId"], 10, 32)
	if e1 != nil {
		writer.WriteHeader(http.StatusInternalServerError)
		util.Logging(util.INFO, method, "", "FAIL on Decline follow request:"+util.Uint2String(uint(followerId))+"->"+util.Uint2String(util.GetLoggedUserIDFromToken(request)), "connection")
	}
	id := util.GetLoggedUserIDFromToken(request)
	if id == 0 {
		writer.WriteHeader(http.StatusUnauthorized)
	}
	_, ok := handler.ConnectionService.DeleteConnection(uint(followerId), id)
	if !ok {
		writer.WriteHeader(http.StatusInternalServerError)
		util.Logging(util.ERROR, method, "", "Service error on decline follow request:"+util.Uint2String(uint(followerId))+"->"+util.Uint2String(id), "connection")
	} else {
		writer.WriteHeader(http.StatusOK)
		writer.Header().Set("Content-Type", "application/json")
		writer.Write([]byte("{\"status\":\"ok\"}"))
		util.Logging(util.INFO, method, "", "Declined follow request:"+util.Uint2String(uint(followerId))+"->"+util.Uint2String(id), "connection")
	}
}

func (handler *Handler) GetAllFollowRequests(writer http.ResponseWriter, request *http.Request) {
	id := util.GetLoggedUserIDFromToken(request)
	if id == 0 {
		writer.WriteHeader(http.StatusBadRequest)
		return
	}
	var userDtos *[]dto.UserDTO = handler.ConnectionService.GetAllFollowRequests(id)
	if userDtos == nil {
		writer.WriteHeader(http.StatusInternalServerError)
		return
	}
	writer.WriteHeader(http.StatusOK)
	writer.Header().Set("Content-Type", "application/json")
	json.NewEncoder(writer).Encode(*userDtos)
}

func (handler *Handler) BlockProfile(writer http.ResponseWriter, request *http.Request) {
	followerId := util.GetLoggedUserIDFromToken(request)
	if followerId == 0 {
		writer.Write([]byte("{\"status\":\"error\"}"))
		writer.WriteHeader(http.StatusBadRequest)
		return
	}
	vars := mux.Vars(request)
	profileId := util.String2Uint(vars["profileId"])

	_, ok := handler.ConnectionService.ToggleBlock(followerId, profileId)

	if !ok {
		writer.Write([]byte("{\"status\":\"error\"}"))
		writer.WriteHeader(http.StatusInternalServerError)
		return
	}

	writer.WriteHeader(http.StatusOK)
	writer.Write([]byte("{\"status\":\"ok\"}"))
	writer.Header().Set("Content-Type", "application/json")
}

func (handler *Handler) MuteProfile(writer http.ResponseWriter, request *http.Request) {
	followerId := util.GetLoggedUserIDFromToken(request)
	if followerId == 0 {
		writer.Write([]byte("{\"status\":\"error\"}"))
		writer.WriteHeader(http.StatusBadRequest)
		return
	}
	vars := mux.Vars(request)
	profileId := util.String2Uint(vars["profileId"])

	_, ok := handler.ConnectionService.ToggleMuted(followerId, profileId)

	if !ok {
		writer.Write([]byte("{\"status\":\"error\"}"))
		writer.WriteHeader(http.StatusInternalServerError)
		return
	}

	writer.WriteHeader(http.StatusOK)
	writer.Write([]byte("{\"status\":\"ok\"}"))
	writer.Header().Set("Content-Type", "application/json")
}
