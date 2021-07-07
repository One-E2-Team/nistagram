package handler

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/google/uuid"
	"github.com/gorilla/mux"
	"io"
	"mime/multipart"
	"net/http"
	"nistagram/notification/model"
	"nistagram/util"
	"os"
	"strings"
)

func (handler *Handler) RndHandler(ctx context.Context, senderId uint, object []byte) ([]byte, bool) {
	span := util.Tracer.StartSpanFromContext(ctx, "RndHandler-service")
	defer util.Tracer.FinishSpan(span)
	util.Tracer.LogFields(span, "service", fmt.Sprintf("servicing id %v\n", senderId))
	fmt.Println(senderId, object)
	return []byte("RandomIndeed"), true
}

func (handler *Handler) SendMessage(ctx context.Context, senderId uint, object []byte) ([]byte, bool) {
	span := util.Tracer.StartSpanFromContext(ctx, "SendMessage-service")
	defer util.Tracer.FinishSpan(span)
	util.Tracer.LogFields(span, "service", fmt.Sprintf("servicing id %v\n", senderId))
	nextCtx := util.Tracer.ContextWithSpan(ctx, span)
	fmt.Println(senderId, object)
	var message model.Message
	err := json.Unmarshal(object, &message)
	if err != nil{
		util.Tracer.LogError(span, err)
		fmt.Println(err)
	}else{
		handler.Service.CreateMessage(nextCtx, &message)
		TrySendMessage(nextCtx, message.ReceiverId, "message", object)
	}
	return []byte("Sent"), true
}

func (handler *Handler) DeleteMessage(w http.ResponseWriter, r *http.Request) {
	span := util.Tracer.StartSpanFromRequest("DeleteMessage-handler", r)
	defer util.Tracer.FinishSpan(span)
	util.Tracer.LogFields(span, "handler", fmt.Sprintf("handling %s\n", r.URL.Path))
	ctx := util.Tracer.ContextWithSpan(context.Background(), span)
	vars := mux.Vars(r)
	messageId := vars["id"]
	loggedUserId := util.GetLoggedUserIDFromToken(r)

	err := handler.Service.DeleteMessage(ctx, loggedUserId, messageId)
	if err != nil{
		util.Tracer.LogError(span, err)
		fmt.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		_, _ = w.Write([]byte("{\"message\":\"error\"}"))
		return
	}

	w.WriteHeader(http.StatusOK)
	_, _ = w.Write([]byte("{\"message\":\"ok\"}"))
	w.Header().Set("Content-Type", "application/json")
}

func (handler *Handler) GetAllMesssages(w http.ResponseWriter, r *http.Request) {
	span := util.Tracer.StartSpanFromRequest("GetAllMesssages-handler", r)
	defer util.Tracer.FinishSpan(span)
	util.Tracer.LogFields(span, "handler", fmt.Sprintf("handling %s\n", r.URL.Path))
	ctx := util.Tracer.ContextWithSpan(context.Background(), span)
	vars := mux.Vars(r)
	profileId := util.String2Uint(vars["id"])
	fmt.Println("Profile id: ", profileId)
	loggedUserId := util.GetLoggedUserIDFromToken(r)

	fmt.Println("Logged user id: ", loggedUserId)

	result, err := handler.Service.GetAllMessages(ctx, loggedUserId, profileId)
	if err != nil{
		util.Tracer.LogError(span, err)
		fmt.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		_, _ = w.Write([]byte("{\"message\":\"error\"}"))
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(result)
	w.Header().Set("Content-Type", "application/json")
}

func (handler *Handler) SaveFile(w http.ResponseWriter, r *http.Request) {
	span := util.Tracer.StartSpanFromRequest("SaveFile-handler", r)
	defer util.Tracer.FinishSpan(span)
	util.Tracer.LogFields(span, "handler", fmt.Sprintf("handling %s\n", r.URL.Path))
	w.Header().Set("Content-Type", "application/json")

	if err := r.ParseMultipartForm(0); err != nil {
		util.Tracer.LogError(span, err)
		fmt.Println(err)
		_, _ = w.Write([]byte("{\"message\":\"error\"}"))
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	file, fileHeader, err := r.FormFile("file")
	if err != nil{
		fmt.Println(err)
		_, _ = w.Write([]byte("{\"message\":\"error\"}"))
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	uid := uuid.NewString()

	fileSplitted := strings.Split(fileHeader.Filename, ".")
	fileName := uid + "." + fileSplitted[1]
	f, err := os.OpenFile("../../nistagramstaticdata/data/"+fileName, os.O_WRONLY|os.O_CREATE, 0666)
	defer func(f *os.File) {
		_ = f.Close()
	}(f)
	defer func(picture multipart.File) {
		_ = picture.Close()
	}(file)
	if err != nil {
		fmt.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		_, _ = w.Write([]byte("{\"message\":\"error\"}"))
		return
	}
	_, err = io.Copy(f, file)
	if err != nil {
		fmt.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		_, _ = w.Write([]byte("{\"message\":\"error\"}"))
		return
	}

	type data struct{
		FileName 	string		`json:"fileName"`
	}
	ret := data{FileName: fileName}
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(ret)
}

func (handler *Handler) Seen(w http.ResponseWriter, r *http.Request) {
	span := util.Tracer.StartSpanFromRequest("Seen-handler", r)
	defer util.Tracer.FinishSpan(span)
	util.Tracer.LogFields(span, "handler", fmt.Sprintf("handling %s\n", r.URL.Path))
	ctx := util.Tracer.ContextWithSpan(context.Background(), span)
	vars := mux.Vars(r)
	messageId := vars["id"]

	err := handler.Service.Seen(ctx, messageId)
	if err != nil{
		fmt.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		_, _ = w.Write([]byte("{\"message\":\"error\"}"))
		return
	}

	w.WriteHeader(http.StatusOK)
	_, _ = w.Write([]byte("{\"message\":\"ok\"}"))
	w.Header().Set("Content-Type", "application/json")
}

func (handler *Handler) GetNotifications(w http.ResponseWriter, r *http.Request) {
	span := util.Tracer.StartSpanFromRequest("GetNotifications-handler", r)
	defer util.Tracer.FinishSpan(span)
	util.Tracer.LogFields(span, "handler", fmt.Sprintf("handling %s\n", r.URL.Path))
	ctx := util.Tracer.ContextWithSpan(context.Background(), span)

	loggedUserId := util.GetLoggedUserIDFromToken(r)
	fmt.Println("Logged user id: ", loggedUserId)

	result, err := handler.Service.GetNotifications(ctx, loggedUserId)
	if err != nil{
		util.Tracer.LogError(span, err)
		fmt.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		_, _ = w.Write([]byte("{\"message\":\"error\"}"))
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(result)
	w.Header().Set("Content-Type", "application/json")
}

func (handler *Handler) SeenMessage(w http.ResponseWriter, r *http.Request) {
	span := util.Tracer.StartSpanFromRequest("SeenMessage-handler", r)
	defer util.Tracer.FinishSpan(span)
	util.Tracer.LogFields(span, "handler", fmt.Sprintf("handling %s\n", r.URL.Path))
	ctx := util.Tracer.ContextWithSpan(context.Background(), span)
	vars := mux.Vars(r)
	senderId := util.String2Uint(vars["senderId"])
	loggedUseId := util.GetLoggedUserIDFromToken(r)

	err := handler.Service.SeenMessage(ctx, loggedUseId, senderId)
	if err != nil{
		fmt.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		_, _ = w.Write([]byte("{\"message\":\"error\"}"))
		return
	}

	w.WriteHeader(http.StatusOK)
	_, _ = w.Write([]byte("{\"message\":\"ok\"}"))
	w.Header().Set("Content-Type", "application/json")
}