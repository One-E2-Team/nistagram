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
)

func (handler *Handler) ToggleNotifyMessageProfile(writer http.ResponseWriter, request *http.Request) {
	span := util.Tracer.StartSpanFromRequest("ToggleNotifyMessageProfile-handler", request)
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

	message, ok := handler.ConnectionService.ToggleNotifyMessage(ctx, followerId, profileId)

	if !ok {
		util.Tracer.LogError(span, fmt.Errorf("error in connection service"))
		writer.Write([]byte("{\"status\":\"error\"}"))
		writer.WriteHeader(http.StatusInternalServerError)
		return
	}

	writer.WriteHeader(http.StatusOK)
	writer.Header().Set("Content-Type", "application/json")
	json.NewEncoder(writer).Encode(*message)
}

func (handler *Handler) GetMessageRelationship(writer http.ResponseWriter, request *http.Request) {
	span := util.Tracer.StartSpanFromRequest("GetMessageRelationship-handler", request)
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

	message := handler.ConnectionService.GetMessage(ctx, followerId, profileId)

	writer.WriteHeader(http.StatusOK)
	writer.Header().Set("Content-Type", "application/json")
	json.NewEncoder(writer).Encode(message)
}

func (handler *Handler) GetMessageRelationships(writer http.ResponseWriter, request *http.Request) {
	span := util.Tracer.StartSpanFromRequest("GetMessageRelationships-handler", request)
	defer util.Tracer.FinishSpan(span)
	util.Tracer.LogFields(span, "handler", fmt.Sprintf("handling %s\n", request.URL.Path))
	ctx := util.Tracer.ContextWithSpan(context.Background(), span)

	type data struct {
		FollowerId string `json:"followerId"`
		Ids []string `json:"ids"`
	}

	var input data
	if err := json.NewDecoder(request.Body).Decode(&input); err != nil {
		util.Tracer.LogError(span, err)
		fmt.Println(err)
		writer.WriteHeader(http.StatusBadRequest)
		return
	}

	var messages []model.MessageEdge

	for _, id := range input.Ids{
		message := handler.ConnectionService.GetMessage(ctx, util.String2Uint(input.FollowerId), util.String2Uint(id))
		if message == nil{
			message = handler.ConnectionService.GetMessage(ctx, util.String2Uint(id), util.String2Uint(input.FollowerId))
		}
		messages = append(messages, *message)
	}

	fmt.Println("Messages len: ", len(messages))
	fmt.Println("Messages: ", messages)

	writer.WriteHeader(http.StatusOK)
	writer.Header().Set("Content-Type", "application/json")
	json.NewEncoder(writer).Encode(messages)
}

func (handler *Handler) MessageRequest(writer http.ResponseWriter, request *http.Request) {
	span := util.Tracer.StartSpanFromRequest("MessageRequest-handler", request)
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

	message, ok := handler.ConnectionService.MessageRequest(ctx, followerId, profileId)

	if !ok || message == nil {
		util.Tracer.LogError(span, fmt.Errorf("error in connection service"))
		writer.Write([]byte("{\"status\":\"error\"}"))
		writer.WriteHeader(http.StatusInternalServerError)
		return
	}

	writer.WriteHeader(http.StatusOK)
	writer.Header().Set("Content-Type", "application/json")
	json.NewEncoder(writer).Encode(*message)
}

func (handler *Handler) DeclineMessageRequest(writer http.ResponseWriter, request *http.Request) {
	span := util.Tracer.StartSpanFromRequest("DeclineMessageRequest-handler", request)
	defer util.Tracer.FinishSpan(span)
	util.Tracer.LogFields(span, "handler", fmt.Sprintf("handling %s\n", request.URL.Path))
	ctx := util.Tracer.ContextWithSpan(context.Background(), span)
	profileId := util.GetLoggedUserIDFromToken(request)
	if profileId == 0 {
		writer.Write([]byte("{\"status\":\"error\"}"))
		writer.WriteHeader(http.StatusBadRequest)
		return
	}
	vars := mux.Vars(request)
	followerId := util.String2Uint(vars["profileId"])

	message, ok := handler.ConnectionService.DeclineMessageRequest(ctx, followerId, profileId)

	if !ok || message == nil {
		util.Tracer.LogError(span, fmt.Errorf("error in connection service"))
		writer.WriteHeader(http.StatusInternalServerError)
		writer.Write([]byte("{\"status\":\"error\"}"))
		return
	}

	writer.WriteHeader(http.StatusOK)
	writer.Header().Set("Content-Type", "application/json")
	json.NewEncoder(writer).Encode(*message)
}

func (handler *Handler) MessageConnect(writer http.ResponseWriter, request *http.Request) {
	span := util.Tracer.StartSpanFromRequest("MessageConnect-handler", request)
	defer util.Tracer.FinishSpan(span)
	util.Tracer.LogFields(span, "handler", fmt.Sprintf("handling %s\n", request.URL.Path))
	ctx := util.Tracer.ContextWithSpan(context.Background(), span)
	profileId := util.GetLoggedUserIDFromToken(request)
	if profileId == 0 {
		writer.Write([]byte("{\"status\":\"error\"}"))
		writer.WriteHeader(http.StatusBadRequest)
		return
	}
	vars := mux.Vars(request)
	followerId := util.String2Uint(vars["profileId"])
	fmt.Println("mess conn")
	message, ok := handler.ConnectionService.MessageConnect(ctx, followerId, profileId)

	if !ok || message == nil {
		util.Tracer.LogError(span, fmt.Errorf("error in connection service"))
		writer.Write([]byte("{\"status\":\"error\"}"))
		writer.WriteHeader(http.StatusInternalServerError)
		return
	}

	writer.WriteHeader(http.StatusOK)
	writer.Header().Set("Content-Type", "application/json")
	json.NewEncoder(writer).Encode(*message)
}


func (handler *Handler) GetAllMessageRequests(writer http.ResponseWriter, request *http.Request) {
	span := util.Tracer.StartSpanFromRequest("GetAllMessageRequests-handler", request)
	defer util.Tracer.FinishSpan(span)
	util.Tracer.LogFields(span, "handler", fmt.Sprintf("handling %s\n", request.URL.Path))
	ctx := util.Tracer.ContextWithSpan(context.Background(), span)
	id := util.GetLoggedUserIDFromToken(request)
	if id == 0 {
		writer.WriteHeader(http.StatusBadRequest)
		return
	}
	var userDtos *[]dto.UserDTO = handler.ConnectionService.GetAllMessageRequests(ctx, id)
	if userDtos == nil {
		util.Tracer.LogError(span, fmt.Errorf("error in connection service"))
		writer.WriteHeader(http.StatusInternalServerError)
		return
	}
	writer.WriteHeader(http.StatusOK)
	writer.Header().Set("Content-Type", "application/json")
	json.NewEncoder(writer).Encode(*userDtos)
}