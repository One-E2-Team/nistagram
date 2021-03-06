package handler

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"net/http"
	"nistagram/connection/dto"
	"nistagram/connection/model"
	"nistagram/util"
	"strconv"
)

func (handler *Handler) FollowRequest(w http.ResponseWriter, r *http.Request) {
	span := util.Tracer.StartSpanFromRequest("FollowRequest-handler", r)
	defer util.Tracer.FinishSpan(span)
	util.Tracer.LogFields(span, "handler", fmt.Sprintf("handling %s\n", r.URL.Path))
	ctx := util.Tracer.ContextWithSpan(context.Background(), span)
	method := "nistagram/connection/handler.FollowRequest"
	vars := mux.Vars(r)
	id2, e2 := strconv.ParseUint(vars["profileId"], 10, 32)
	if e2 != nil {
		util.Tracer.LogError(span, e2)
		w.WriteHeader(http.StatusBadRequest)
		util.Logging(util.ERROR, method, "", "Failed Follow request (possible auto-approve): "+util.GetLoggingStringFromID(util.GetLoggedUserIDFromToken(r))+"->"+util.Uint2String(uint(id2)), "connection")
		return
	}
	id := util.GetLoggedUserIDFromToken(r)
	if id == 0 {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}
	connection, ok := handler.ConnectionService.FollowRequest(ctx, id, uint(id2))
	if ok && connection != nil {
		w.WriteHeader(http.StatusOK)
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(connection)
		util.Logging(util.INFO, method, "", "Follow request (possible auto-approve): "+util.Uint2String(id)+"->"+util.Uint2String(uint(id2)), "connection")
	} else {
		util.Tracer.LogError(span, fmt.Errorf("error in connection service"))
		w.WriteHeader(http.StatusInternalServerError)
		util.Logging(util.ERROR, method, "", "Service call failed:"+util.GetLoggingStringFromID(util.GetLoggedUserIDFromToken(r))+"->"+util.Uint2String(uint(id2)), "connection")
	}
}

func (handler *Handler) UnfollowProfile(w http.ResponseWriter, r *http.Request) {
	span := util.Tracer.StartSpanFromRequest("UnfollowProfile-handler", r)
	defer util.Tracer.FinishSpan(span)
	util.Tracer.LogFields(span, "handler", fmt.Sprintf("handling %s\n", r.URL.Path))
	ctx := util.Tracer.ContextWithSpan(context.Background(), span)
	vars := mux.Vars(r)
	id1, e1 := strconv.ParseUint(vars["profileId"], 10, 32)
	if e1 != nil {
		util.Tracer.LogError(span, e1)
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	id := util.GetLoggedUserIDFromToken(r)
	if id == 0 {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}
	connection, ok := handler.ConnectionService.Unfollow(ctx, id, uint(id1))
	if ok {
		w.WriteHeader(http.StatusOK)
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte("{\"status\":\"ok\"}"))
		json.NewEncoder(w).Encode(*connection)
	} else {
		util.Tracer.LogError(span, fmt.Errorf("error in connection service"))
		w.WriteHeader(http.StatusInternalServerError)
	}
}

func (handler *Handler) FollowApprove(w http.ResponseWriter, r *http.Request) {
	span := util.Tracer.StartSpanFromRequest("FollowApprove-handler", r)
	defer util.Tracer.FinishSpan(span)
	util.Tracer.LogFields(span, "handler", fmt.Sprintf("handling %s\n", r.URL.Path))
	ctx := util.Tracer.ContextWithSpan(context.Background(), span)
	method := "nistagram/connection/handler.FollowApprove"
	vars := mux.Vars(r)
	id1, e1 := strconv.ParseUint(vars["profileId"], 10, 32)
	if e1 != nil {
		util.Tracer.LogError(span, e1)
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
	_, ok := handler.ConnectionService.ApproveConnection(ctx, uint(id1), id)
	if ok {
		w.WriteHeader(http.StatusOK)
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte("{\"status\":\"ok\"}"))
		util.Logging(util.INFO, method, "", "Approved follow request: "+util.Uint2String(id)+"->"+util.Uint2String(uint(id1)), "connection")
	} else {
		util.Tracer.LogError(span, fmt.Errorf("error in connection service"))
		util.Logging(util.ERROR, method, "", "Service Error", "connection")
		w.WriteHeader(http.StatusInternalServerError)
	}
}

func (handler *Handler) GetFollowedProfilesNotMuted(w http.ResponseWriter, r *http.Request) {
	span := util.Tracer.StartSpanFromRequest("GetFollowedProfilesNotMuted-handler", r)
	defer util.Tracer.FinishSpan(span)
	util.Tracer.LogFields(span, "handler", fmt.Sprintf("handling %s\n", r.URL.Path))
	ctx := util.Tracer.ContextWithSpan(context.Background(), span)
	vars := mux.Vars(r)
	id, _ := strconv.ParseUint(vars["id"], 10, 32)
	conn := model.ConnectionEdge{
		PrimaryProfile:    uint(id),
		SecondaryProfile:  0,
		Muted:             false,
		CloseFriend:       false,
		NotifyPost:        false,
		NotifyStory:       false,
		NotifyComment:     false,
		ConnectionRequest: false,
		Approved:          true,
	}
	profiles := handler.ConnectionService.GetConnectedProfiles(ctx, conn, true, false)
	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(*profiles)
}

func (handler *Handler) GetFollowedProfiles(w http.ResponseWriter, r *http.Request) {
	span := util.Tracer.StartSpanFromRequest("GetFollowedProfiles-handler", r)
	defer util.Tracer.FinishSpan(span)
	util.Tracer.LogFields(span, "handler", fmt.Sprintf("handling %s\n", r.URL.Path))
	ctx := util.Tracer.ContextWithSpan(context.Background(), span)
	vars := mux.Vars(r)
	id, _ := strconv.ParseUint(vars["id"], 10, 32)
	conn := model.ConnectionEdge{
		PrimaryProfile:    uint(id),
		SecondaryProfile:  0,
		Muted:             false,
		CloseFriend:       false,
		NotifyPost:        false,
		NotifyStory:       false,
		NotifyComment:     false,
		ConnectionRequest: false,
		Approved:          true,
	}
	profiles := handler.ConnectionService.GetProfilesInFollowRelationship(ctx, conn, false, false, true)
	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(*profiles)
}

func (handler *Handler) GetMyFollowedProfiles(w http.ResponseWriter, r *http.Request) {
	span := util.Tracer.StartSpanFromRequest("GetMyFollowedProfiles-handler", r)
	defer util.Tracer.FinishSpan(span)
	util.Tracer.LogFields(span, "handler", fmt.Sprintf("handling %s\n", r.URL.Path))
	ctx := util.Tracer.ContextWithSpan(context.Background(), span)
	id := util.GetLoggedUserIDFromToken(r)
	conn := model.ConnectionEdge{
		PrimaryProfile:    id,
		SecondaryProfile:  0,
		Muted:             false,
		CloseFriend:       false,
		NotifyPost:        false,
		NotifyStory:       false,
		NotifyComment:     false,
		ConnectionRequest: false,
		Approved:          true,
	}
	profiles := handler.ConnectionService.GetProfilesInFollowRelationship(ctx, conn, false, false, true)
	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(*profiles)
}

func (handler *Handler) GetFollowerProfiles(w http.ResponseWriter, r *http.Request) {
	span := util.Tracer.StartSpanFromRequest("GetFollowerProfiles-handler", r)
	defer util.Tracer.FinishSpan(span)
	util.Tracer.LogFields(span, "handler", fmt.Sprintf("handling %s\n", r.URL.Path))
	ctx := util.Tracer.ContextWithSpan(context.Background(), span)
	vars := mux.Vars(r)
	id, _ := strconv.ParseUint(vars["id"], 10, 32)
	conn := model.ConnectionEdge{
		PrimaryProfile:    0,
		SecondaryProfile:  uint(id),
		Muted:             false,
		CloseFriend:       false,
		NotifyPost:        false,
		NotifyStory:       false,
		NotifyComment:     false,
		ConnectionRequest: false,
		Approved:          true,
	}
	profiles := handler.ConnectionService.GetProfilesInFollowRelationship(ctx, conn, false, false, false)
	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(*profiles)
}

func (handler *Handler) GetMyFollowerProfiles(w http.ResponseWriter, r *http.Request) {
	span := util.Tracer.StartSpanFromRequest("GetMyFollowerProfiles-handler", r)
	defer util.Tracer.FinishSpan(span)
	util.Tracer.LogFields(span, "handler", fmt.Sprintf("handling %s\n", r.URL.Path))
	ctx := util.Tracer.ContextWithSpan(context.Background(), span)
	id := util.GetLoggedUserIDFromToken(r)
	conn := model.ConnectionEdge{
		PrimaryProfile:    0,
		SecondaryProfile:  id,
		Muted:             false,
		CloseFriend:       false,
		NotifyPost:        false,
		NotifyStory:       false,
		NotifyComment:     false,
		ConnectionRequest: false,
		Approved:          true,
	}
	profiles := handler.ConnectionService.GetProfilesInFollowRelationship(ctx, conn, false, false, false)
	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(*profiles)
}

func (handler *Handler) DeclineFollowRequest(writer http.ResponseWriter, request *http.Request) {
	span := util.Tracer.StartSpanFromRequest("DeclineFollowRequest-handler", request)
	defer util.Tracer.FinishSpan(span)
	util.Tracer.LogFields(span, "handler", fmt.Sprintf("handling %s\n", request.URL.Path))
	ctx := util.Tracer.ContextWithSpan(context.Background(), span)
	method := "nistagram/connection/handler.DeclineFollowRequest"
	vars := mux.Vars(request)
	followerId, e1 := strconv.ParseUint(vars["profileId"], 10, 32)
	if e1 != nil {
		util.Tracer.LogError(span, e1)
		writer.WriteHeader(http.StatusInternalServerError)
		util.Logging(util.INFO, method, "", "FAIL on Decline follow request:"+util.Uint2String(uint(followerId))+"->"+util.Uint2String(util.GetLoggedUserIDFromToken(request)), "connection")
	}
	id := util.GetLoggedUserIDFromToken(request)
	if id == 0 {
		writer.WriteHeader(http.StatusUnauthorized)
	}
	_, ok := handler.ConnectionService.DeleteConnection(ctx, uint(followerId), id)
	if !ok {
		util.Tracer.LogError(span, fmt.Errorf("error in connection service"))
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
	span := util.Tracer.StartSpanFromRequest("GetAllFollowRequests-handler", request)
	defer util.Tracer.FinishSpan(span)
	util.Tracer.LogFields(span, "handler", fmt.Sprintf("handling %s\n", request.URL.Path))
	ctx := util.Tracer.ContextWithSpan(context.Background(), span)
	id := util.GetLoggedUserIDFromToken(request)
	if id == 0 {
		writer.WriteHeader(http.StatusBadRequest)
		return
	}
	var userDtos *[]dto.UserDTO = handler.ConnectionService.GetAllFollowRequests(ctx, id)
	if userDtos == nil {
		util.Tracer.LogError(span, fmt.Errorf("error in connection service"))
		writer.WriteHeader(http.StatusInternalServerError)
		return
	}
	writer.WriteHeader(http.StatusOK)
	writer.Header().Set("Content-Type", "application/json")
	json.NewEncoder(writer).Encode(*userDtos)
}