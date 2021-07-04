package handler

import (
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

func (handler *Handler) RndHandler(senderId uint, object []byte) ([]byte, bool) {
	fmt.Println(senderId, object)
	return []byte("RandomIndeed"), true
}

func (handler *Handler) SendMessage(senderId uint, object []byte) ([]byte, bool) {
	fmt.Println(senderId, object)
	var message model.Message
	err := json.Unmarshal(object, &message)
	if err != nil{
		fmt.Println(err)
	}else{
		handler.Service.CreateMessage(&message)
		TrySendMessage(message.ReceiverId, "message", object)
	}
	return []byte("Sent"), true
}

func (handler *Handler) DeleteMessage(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	messageId := vars["id"]
	loggedUserId := util.GetLoggedUserIDFromToken(r)

	err := handler.Service.DeleteMessage(loggedUserId, messageId)
	if err != nil{
		fmt.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		_, _ = w.Write([]byte("{\"message\":\"error\"}"))
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode("{\"message\":\"ok\"}")
	w.Header().Set("Content-Type", "application/json")
}

func (handler *Handler) GetAllMesssages(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	profileId := util.String2Uint(vars["id"])
	fmt.Println("Profile id: ", profileId)
	loggedUserId := util.GetLoggedUserIDFromToken(r)

	fmt.Println("Logged user id: ", loggedUserId)

	result, err := handler.Service.GetAllMessages(loggedUserId, profileId)
	if err != nil{
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
	w.Header().Set("Content-Type", "application/json")

	if err := r.ParseMultipartForm(0); err != nil {
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