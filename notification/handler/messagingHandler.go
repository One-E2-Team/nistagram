package handler

import (
	"encoding/json"
	"fmt"
	"nistagram/notification/model"
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