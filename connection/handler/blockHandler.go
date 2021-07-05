package handler

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"net/http"
	"nistagram/util"
	"strconv"
)

func (handler *Handler) IsBlocked(w http.ResponseWriter, r *http.Request) {
	span := util.Tracer.StartSpanFromRequest("IsBlocked-handler", r)
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
	ok := handler.ConnectionService.IsBlocked(ctx, id, uint(id1))
	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json")
	type resp struct {
		Blocked bool `json:"blocked"`
	}
	json.NewEncoder(w).Encode(resp{Blocked: ok})
}



func (handler *Handler) ToggleBlockProfile(writer http.ResponseWriter, request *http.Request) {
	span := util.Tracer.StartSpanFromRequest("ToggleBlockProfile-handler", request)
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

	_, ok := handler.ConnectionService.ToggleBlock(ctx, followerId, profileId)

	if !ok {
		util.Tracer.LogError(span, fmt.Errorf("error in connection service"))
		writer.Write([]byte("{\"status\":\"error\"}"))
		writer.WriteHeader(http.StatusInternalServerError)
		return
	}

	writer.WriteHeader(http.StatusOK)
	writer.Write([]byte("{\"status\":\"ok\"}"))
	writer.Header().Set("Content-Type", "application/json")
}

func (handler *Handler) GetBlockingRelationships(writer http.ResponseWriter, request *http.Request) {
	span := util.Tracer.StartSpanFromRequest("GetBlockingRelationships-handler", request)
	defer util.Tracer.FinishSpan(span)
	util.Tracer.LogFields(span, "handler", fmt.Sprintf("handling %s\n", request.URL.Path))
	ctx := util.Tracer.ContextWithSpan(context.Background(), span)
	vars := mux.Vars(request)
	id, _ := strconv.ParseUint(vars["id"], 10, 32)
	profiles := handler.ConnectionService.GetBlockingRelationships(ctx, uint(id))
	writer.WriteHeader(http.StatusOK)
	writer.Header().Set("Content-Type", "application/json")
	if js, err := json.Marshal(*profiles); err != nil {
		util.Tracer.LogError(span, err)
		writer.WriteHeader(http.StatusInternalServerError)
	} else {
		writer.WriteHeader(http.StatusOK)
		_, _ = writer.Write(js)
	}
}

func (handler *Handler) AmBlocked(writer http.ResponseWriter, request *http.Request) {
	span := util.Tracer.StartSpanFromRequest("AmBlocked-handler", request)
	defer util.Tracer.FinishSpan(span)
	util.Tracer.LogFields(span, "handler", fmt.Sprintf("handling %s\n", request.URL.Path))
	ctx := util.Tracer.ContextWithSpan(context.Background(), span)
	vars := mux.Vars(request)
	id1, e1 := strconv.ParseUint(vars["profileId"], 10, 32)
	if e1 != nil {
		util.Tracer.LogError(span, e1)
		writer.WriteHeader(http.StatusBadRequest)
		return
	}
	id := util.GetLoggedUserIDFromToken(request)
	if id == 0 {
		writer.WriteHeader(http.StatusUnauthorized)
		return
	}
	ok := handler.ConnectionService.IsBlocked(ctx, uint(id1), id)
	writer.WriteHeader(http.StatusOK)
	writer.Header().Set("Content-Type", "application/json")
	type resp struct {
		Blocked bool `json:"blocked"`
	}
	json.NewEncoder(writer).Encode(resp{Blocked: ok})
}