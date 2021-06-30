package handler

import (
	"fmt"
	"io"
	"net/http"
	"nistagram/agent/service"
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
	//TODO: extract this in method
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		fmt.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	defer func(Body io.ReadCloser) {
		_ = Body.Close()
	}(resp.Body)

	w.WriteHeader(http.StatusOK)
	_, _ = w.Write(body)
	w.Header().Set("Content-Type", "application/json")
}
