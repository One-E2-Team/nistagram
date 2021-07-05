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

func (handler *Handler) AddProfile(w http.ResponseWriter, r *http.Request) {
	span := util.Tracer.StartSpanFromRequest("AddProfile-handler", r)
	defer util.Tracer.FinishSpan(span)
	util.Tracer.LogFields(span, "handler", fmt.Sprintf("handling %s\n", r.URL.Path))
	ctx := util.Tracer.ContextWithSpan(context.Background(), span)
	method := "nistagram/connection/handler.AddProfile"
	vars := mux.Vars(r)
	id, _ := strconv.ParseUint(vars["id"], 10, 32)
	profile, ok := handler.ConnectionService.AddOrUpdateProfile(ctx, model.ProfileVertex{ProfileID: uint(id)})
	if ok {
		w.WriteHeader(http.StatusOK)
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(profile)
		util.Logging(util.INFO, method, "", "Added user: "+util.Uint2String(uint(id)), "connection")
	} else {
		util.Tracer.LogError(span, fmt.Errorf("error in connection service"))
		w.WriteHeader(http.StatusInternalServerError)
		util.Logging(util.INFO, method, "", "Failed to add user: "+util.Uint2String(uint(id)), "connection")
	}
}

func (handler *Handler) DeleteProfile(w http.ResponseWriter, r *http.Request) {
	span := util.Tracer.StartSpanFromRequest("DeleteProfile-handler", r)
	defer util.Tracer.FinishSpan(span)
	util.Tracer.LogFields(span, "handler", fmt.Sprintf("handling %s\n", r.URL.Path))
	ctx := util.Tracer.ContextWithSpan(context.Background(), span)
	vars := mux.Vars(r)
	id, _ := strconv.ParseUint(vars["id"], 10, 32)
	profile, ok := handler.ConnectionService.AddOrUpdateProfile(ctx, model.ProfileVertex{ProfileID: uint(id), Deleted: true})
	if ok {
		w.WriteHeader(http.StatusOK)
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(profile)
	} else {
		util.Tracer.LogError(span, fmt.Errorf("error in connection service"))
		w.WriteHeader(http.StatusInternalServerError)
	}
}

func (handler *Handler) ReActivateProfile(w http.ResponseWriter, r *http.Request) {
	span := util.Tracer.StartSpanFromRequest("AddProfile-handler", r)
	defer util.Tracer.FinishSpan(span)
	util.Tracer.LogFields(span, "handler", fmt.Sprintf("handling %s\n", r.URL.Path))
	ctx := util.Tracer.ContextWithSpan(context.Background(), span)
	vars := mux.Vars(r)
	id, _ := strconv.ParseUint(vars["id"], 10, 32)
	profile, ok := handler.ConnectionService.AddOrUpdateProfile(ctx, model.ProfileVertex{ProfileID: uint(id), Deleted: false})
	if ok {
		w.WriteHeader(http.StatusOK)
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(profile)
	} else {
		util.Tracer.LogError(span, fmt.Errorf("error in connection service"))
		w.WriteHeader(http.StatusInternalServerError)
	}
}


func (handler *Handler) RecommendProfiles(writer http.ResponseWriter, request *http.Request) {
	span := util.Tracer.StartSpanFromRequest("RecommendProfiles-handler", request)
	defer util.Tracer.FinishSpan(span)
	util.Tracer.LogFields(span, "handler", fmt.Sprintf("handling %s\n", request.URL.Path))
	ctx := util.Tracer.ContextWithSpan(context.Background(), span)
	profileId := util.GetLoggedUserIDFromToken(request) // TODO ctx
	if profileId == 0 {
		util.Tracer.LogError(span, fmt.Errorf("coukd not get id from token"))
		writer.Write([]byte("[{\"status\":\"error\"}]"))
		writer.WriteHeader(http.StatusBadRequest)
		return
	}

	var users *[]dto.ProfileRecommendationDTO
	users, ok := handler.ConnectionService.GetRecommendations(ctx, profileId)

	if !ok {
		util.Tracer.LogError(span, fmt.Errorf("error in recomendation service"))
		writer.Write([]byte("[{\"status\":\"error\"}]"))
		writer.WriteHeader(http.StatusInternalServerError)
		return
	}

	writer.WriteHeader(http.StatusOK)
	writer.Header().Set("Content-Type", "application/json")
	json.NewEncoder(writer).Encode(users)
}