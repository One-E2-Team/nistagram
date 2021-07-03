package handler

import (
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"net/http"
	"nistagram/notification/model"
	"nistagram/util"
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