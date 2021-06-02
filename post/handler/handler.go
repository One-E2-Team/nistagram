package handler

import (
	"encoding/json"
	"net/http"
	"nistagram/post/dto"
	"nistagram/post/service"
)

type Handler struct {
	PostService *service.PostService
}

func (handler *Handler) Create(w http.ResponseWriter, r *http.Request){
	var postDto dto.PostDto
	err := json.NewDecoder(r.Body).Decode(&postDto)

	if err != nil{
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	err = handler.PostService.CreatePost(postDto)

	if err != nil{
		w.WriteHeader(http.StatusBadRequest)
		return
	}else {
		w.WriteHeader(http.StatusCreated)
	}
	w.Write([]byte("{\"success\":\"ok\"}"))
	w.Header().Set("Content-Type", "application/json")

}
