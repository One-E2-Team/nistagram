package handler

import (
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"net/http"
	"nistagram/connection/dto"
	"nistagram/connection/model"
	"nistagram/util"
)

func (handler *Handler) ToggleNotifyMessageProfile(writer http.ResponseWriter, request *http.Request) {
	followerId := util.GetLoggedUserIDFromToken(request)
	if followerId == 0 {
		writer.Write([]byte("{\"status\":\"error\"}"))
		writer.WriteHeader(http.StatusBadRequest)
		return
	}
	vars := mux.Vars(request)
	profileId := util.String2Uint(vars["profileId"])

	message, ok := handler.ConnectionService.ToggleNotifyMessage(followerId, profileId)

	if !ok {
		writer.Write([]byte("{\"status\":\"error\"}"))
		writer.WriteHeader(http.StatusInternalServerError)
		return
	}

	writer.WriteHeader(http.StatusOK)
	writer.Header().Set("Content-Type", "application/json")
	json.NewEncoder(writer).Encode(*message)
}

func (handler *Handler) GetMessageRelationship(writer http.ResponseWriter, request *http.Request) {
	followerId := util.GetLoggedUserIDFromToken(request)
	if followerId == 0 {
		writer.Write([]byte("{\"status\":\"error\"}"))
		writer.WriteHeader(http.StatusBadRequest)
		return
	}
	vars := mux.Vars(request)
	profileId := util.String2Uint(vars["profileId"])

	message := handler.ConnectionService.GetMessage(followerId, profileId)

	writer.WriteHeader(http.StatusOK)
	writer.Header().Set("Content-Type", "application/json")
	json.NewEncoder(writer).Encode(message)
}

func (handler *Handler) GetMessageRelationships(writer http.ResponseWriter, request *http.Request) {
	followerId := util.GetLoggedUserIDFromToken(request)
	if followerId == 0 {
		writer.Write([]byte("{\"status\":\"error\"}"))
		writer.WriteHeader(http.StatusBadRequest)
		return
	}

	type data struct {
		Ids []string `json:"ids"`
	}

	var input data
	if err := json.NewDecoder(request.Body).Decode(&input); err != nil {
		fmt.Println(err)
		writer.WriteHeader(http.StatusBadRequest)
		return
	}

	var messages []model.MessageEdge

	for _, id := range input.Ids{
		message := handler.ConnectionService.GetMessage(followerId, util.String2Uint(id))
		messages = append(messages, *message)
	}

	writer.WriteHeader(http.StatusOK)
	writer.Header().Set("Content-Type", "application/json")
	json.NewEncoder(writer).Encode(messages)
}

func (handler *Handler) MessageRequest(writer http.ResponseWriter, request *http.Request) {
	followerId := util.GetLoggedUserIDFromToken(request)
	if followerId == 0 {
		writer.Write([]byte("{\"status\":\"error\"}"))
		writer.WriteHeader(http.StatusBadRequest)
		return
	}
	vars := mux.Vars(request)
	profileId := util.String2Uint(vars["profileId"])

	message, ok := handler.ConnectionService.MessageRequest(followerId, profileId)

	if !ok || message == nil {
		writer.Write([]byte("{\"status\":\"error\"}"))
		writer.WriteHeader(http.StatusInternalServerError)
		return
	}

	writer.WriteHeader(http.StatusOK)
	writer.Header().Set("Content-Type", "application/json")
	json.NewEncoder(writer).Encode(*message)
}

func (handler *Handler) DeclineMessageRequest(writer http.ResponseWriter, request *http.Request) {
	profileId := util.GetLoggedUserIDFromToken(request)
	if profileId == 0 {
		writer.Write([]byte("{\"status\":\"error\"}"))
		writer.WriteHeader(http.StatusBadRequest)
		return
	}
	vars := mux.Vars(request)
	followerId := util.String2Uint(vars["profileId"])

	message, ok := handler.ConnectionService.DeclineMessageRequest(followerId, profileId)

	if !ok || message == nil {
		writer.WriteHeader(http.StatusInternalServerError)
		writer.Write([]byte("{\"status\":\"error\"}"))
		return
	}

	writer.WriteHeader(http.StatusOK)
	writer.Header().Set("Content-Type", "application/json")
	json.NewEncoder(writer).Encode(*message)
}

func (handler *Handler) MessageConnect(writer http.ResponseWriter, request *http.Request) {
	profileId := util.GetLoggedUserIDFromToken(request)
	if profileId == 0 {
		writer.Write([]byte("{\"status\":\"error\"}"))
		writer.WriteHeader(http.StatusBadRequest)
		return
	}
	vars := mux.Vars(request)
	followerId := util.String2Uint(vars["profileId"])
	fmt.Println("mess conn")
	message, ok := handler.ConnectionService.MessageConnect(followerId, profileId)

	if !ok || message == nil {
		writer.Write([]byte("{\"status\":\"error\"}"))
		writer.WriteHeader(http.StatusInternalServerError)
		return
	}

	writer.WriteHeader(http.StatusOK)
	writer.Header().Set("Content-Type", "application/json")
	json.NewEncoder(writer).Encode(*message)
}


func (handler *Handler) GetAllMessageRequests(writer http.ResponseWriter, request *http.Request) {
	id := util.GetLoggedUserIDFromToken(request)
	if id == 0 {
		writer.WriteHeader(http.StatusBadRequest)
		return
	}
	var userDtos *[]dto.UserDTO = handler.ConnectionService.GetAllMessageRequests(id)
	if userDtos == nil {
		writer.WriteHeader(http.StatusInternalServerError)
		return
	}
	writer.WriteHeader(http.StatusOK)
	writer.Header().Set("Content-Type", "application/json")
	json.NewEncoder(writer).Encode(*userDtos)
}