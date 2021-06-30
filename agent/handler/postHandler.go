package handler

import (
	"fmt"
	"net/http"
	"nistagram/agent/service"
	"nistagram/agent/util"
)

type PostHandler struct {
	PostService *service.PostService
}

func (handler *PostHandler) GetMyPosts(w http.ResponseWriter, r *http.Request) {
	resp, err := handler.PostService.GetMyPosts()
	if err != nil {
		fmt.Println(err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	w.WriteHeader(http.StatusOK)
	_, _ = w.Write(util.GetResponseJSON(*resp))
	w.Header().Set("Content-Type", "application/json")
}
