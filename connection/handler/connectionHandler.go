package handler

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"net/http"
	"nistagram/connection/model"
	"nistagram/util"
	"strconv"
)

func (handler *Handler) GetConnection(w http.ResponseWriter, r *http.Request) {
	span := util.Tracer.StartSpanFromRequest("GetConnection-handler", r)
	defer util.Tracer.FinishSpan(span)
	util.Tracer.LogFields(span, "handler", fmt.Sprintf("handling %s\n", r.URL.Path))
	ctx := util.Tracer.ContextWithSpan(context.Background(), span)
	vars := mux.Vars(r)
	id1, e1 := strconv.ParseUint(vars["followerId"], 10, 32)
	id2, e2 := strconv.ParseUint(vars["profileId"], 10, 32)
	if e1 != nil || e2 != nil {
		util.Tracer.LogError(span, fmt.Errorf("%s \n %s", e1, e2))
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	connection := handler.ConnectionService.GetConnection(ctx, uint(id1), uint(id2))
	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(connection)
}

func (handler *Handler) GetConnectionPublic(w http.ResponseWriter, r *http.Request) {
	span := util.Tracer.StartSpanFromRequest("GetConnectionPublic-handler", r)
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
	connection := handler.ConnectionService.GetConnection(ctx, id, uint(id1))
	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(connection)
}

func (handler *Handler) UpdateConnection(w http.ResponseWriter, r *http.Request) {
	span := util.Tracer.StartSpanFromRequest("UpdateConnection-handler", r)
	defer util.Tracer.FinishSpan(span)
	util.Tracer.LogFields(span, "handler", fmt.Sprintf("handling %s\n", r.URL.Path))
	ctx := util.Tracer.ContextWithSpan(context.Background(), span)
	var dto model.ConnectionEdge
	err := json.NewDecoder(r.Body).Decode(&dto)
	if err != nil {
		util.Tracer.LogError(span, err)
		w.WriteHeader(http.StatusInternalServerError)
	} else {
		id := util.GetLoggedUserIDFromToken(r)
		ret, ok := handler.ConnectionService.UpdateConnection(ctx, id, dto)
		if !ok {
			w.WriteHeader(http.StatusInternalServerError)
			util.Tracer.LogError(span, fmt.Errorf("error in connection service"))
		} else {
			w.WriteHeader(http.StatusOK)
			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(*ret)

		}
	}
}

func (handler *Handler) ToggleMuteProfile(writer http.ResponseWriter, request *http.Request) {
	span := util.Tracer.StartSpanFromRequest("ToggleMuteProfile-handler", request)
	defer util.Tracer.FinishSpan(span)
	util.Tracer.LogFields(span, "handler", fmt.Sprintf("handling %s\n", request.URL.Path))
	ctx := util.Tracer.ContextWithSpan(context.Background(), span)
	followerId := util.GetLoggedUserIDFromToken(request)
	if followerId == 0 {
		writer.Write([]byte("{\"status\":\"error\"}"))
		writer.WriteHeader(http.StatusBadRequest)
		return
	}
	vars := mux.Vars(request)
	profileId := util.String2Uint(vars["profileId"])

	connection, ok := handler.ConnectionService.ToggleMuted(ctx, followerId, profileId)

	if !ok {
		util.Tracer.LogError(span, fmt.Errorf("error in connection service"))
		writer.Write([]byte("{\"status\":\"error\"}"))
		writer.WriteHeader(http.StatusInternalServerError)
		return
	}

	writer.WriteHeader(http.StatusOK)
	writer.Header().Set("Content-Type", "application/json")
	json.NewEncoder(writer).Encode(*connection)
}

func (handler *Handler) ToggleCloseFriendProfile(writer http.ResponseWriter, request *http.Request) {
	span := util.Tracer.StartSpanFromRequest("ToggleCloseFriendProfile-handler", request)
	defer util.Tracer.FinishSpan(span)
	util.Tracer.LogFields(span, "handler", fmt.Sprintf("handling %s\n", request.URL.Path))
	ctx := util.Tracer.ContextWithSpan(context.Background(), span)
	followerId := util.GetLoggedUserIDFromToken(request)
	if followerId == 0 {
		writer.Write([]byte("{\"status\":\"error\"}"))
		writer.WriteHeader(http.StatusBadRequest)
		return
	}
	vars := mux.Vars(request)
	profileId := util.String2Uint(vars["profileId"])

	connection, ok := handler.ConnectionService.ToggleCloseFriend(ctx, followerId, profileId)

	if !ok {
		util.Tracer.LogError(span, fmt.Errorf("error in connection service"))
		writer.Write([]byte("{\"status\":\"error\"}"))
		writer.WriteHeader(http.StatusInternalServerError)
		return
	}

	writer.WriteHeader(http.StatusOK)
	writer.Header().Set("Content-Type", "application/json")
	json.NewEncoder(writer).Encode(*connection)
}

func (handler *Handler) ToggleNotifyPostProfile(writer http.ResponseWriter, request *http.Request) {
	span := util.Tracer.StartSpanFromRequest("ToggleNotifyPostProfile-handler", request)
	defer util.Tracer.FinishSpan(span)
	util.Tracer.LogFields(span, "handler", fmt.Sprintf("handling %s\n", request.URL.Path))
	ctx := util.Tracer.ContextWithSpan(context.Background(), span)
	followerId := util.GetLoggedUserIDFromToken(request)
	if followerId == 0 {
		writer.Write([]byte("{\"status\":\"error\"}"))
		writer.WriteHeader(http.StatusBadRequest)
		return
	}
	vars := mux.Vars(request)
	profileId := util.String2Uint(vars["profileId"])

	connection, ok := handler.ConnectionService.ToggleNotifyPost(ctx, followerId, profileId)

	if !ok {
		util.Tracer.LogError(span, fmt.Errorf("error in connection service"))
		writer.Write([]byte("{\"status\":\"error\"}"))
		writer.WriteHeader(http.StatusInternalServerError)
		return
	}

	writer.WriteHeader(http.StatusOK)
	writer.Header().Set("Content-Type", "application/json")
	json.NewEncoder(writer).Encode(*connection)
}

func (handler *Handler) ToggleNotifyStoryProfile(writer http.ResponseWriter, request *http.Request) {
	span := util.Tracer.StartSpanFromRequest("ToggleNotifyStoryProfile-handler", request)
	defer util.Tracer.FinishSpan(span)
	util.Tracer.LogFields(span, "handler", fmt.Sprintf("handling %s\n", request.URL.Path))
	ctx := util.Tracer.ContextWithSpan(context.Background(), span)
	followerId := util.GetLoggedUserIDFromToken(request)
	if followerId == 0 {
		writer.Write([]byte("{\"status\":\"error\"}"))
		writer.WriteHeader(http.StatusBadRequest)
		return
	}
	vars := mux.Vars(request)
	profileId := util.String2Uint(vars["profileId"])

	connection, ok := handler.ConnectionService.ToggleNotifyStory(ctx, followerId, profileId)

	if !ok {
		util.Tracer.LogError(span, fmt.Errorf("error in connection service"))
		writer.Write([]byte("{\"status\":\"error\"}"))
		writer.WriteHeader(http.StatusInternalServerError)
		return
	}

	writer.WriteHeader(http.StatusOK)
	writer.Header().Set("Content-Type", "application/json")
	json.NewEncoder(writer).Encode(*connection)
}

func (handler *Handler) ToggleNotifyCommentProfile(writer http.ResponseWriter, request *http.Request) {
	span := util.Tracer.StartSpanFromRequest("ToggleNotifyCommentProfile-handler", request)
	defer util.Tracer.FinishSpan(span)
	util.Tracer.LogFields(span, "handler", fmt.Sprintf("handling %s\n", request.URL.Path))
	ctx := util.Tracer.ContextWithSpan(context.Background(), span)
	followerId := util.GetLoggedUserIDFromToken(request)
	if followerId == 0 {
		writer.Write([]byte("{\"status\":\"error\"}"))
		writer.WriteHeader(http.StatusBadRequest)
		return
	}
	vars := mux.Vars(request)
	profileId := util.String2Uint(vars["profileId"])

	connection, ok := handler.ConnectionService.ToggleNotifyComment(ctx, followerId, profileId)

	if !ok {
		util.Tracer.LogError(span, fmt.Errorf("error in connection service"))
		writer.Write([]byte("{\"status\":\"error\"}"))
		writer.WriteHeader(http.StatusInternalServerError)
		return
	}

	writer.WriteHeader(http.StatusOK)
	writer.Header().Set("Content-Type", "application/json")
	json.NewEncoder(writer).Encode(*connection)
}
