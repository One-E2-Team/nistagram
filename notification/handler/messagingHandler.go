package handler

import "fmt"

func (handler *Handler) RndHandler(senderId uint, object []byte) ([]byte, bool) {
	fmt.Println(senderId, object)
	return []byte("RandomIndeed"), true
}