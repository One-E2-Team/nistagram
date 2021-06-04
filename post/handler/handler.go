package handler

import (
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"net/http"
	"nistagram/post/dto"
	"nistagram/post/model"
	"nistagram/post/service"
)

type Handler struct {
	PostService *service.PostService
}

func (handler *Handler) Create(w http.ResponseWriter, r *http.Request){
	(w).Header().Set("Access-Control-Allow-Origin", "*")

	var postDto dto.PostDto
	err := json.NewDecoder(r.Body).Decode(&postDto)

	postType := model.GetPostType(postDto.PostType)

	fmt.Println(postDto)

	if err != nil || postType == model.NONE{
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	err = handler.PostService.CreatePost(postType,postDto)

	if err != nil{
		w.WriteHeader(http.StatusBadRequest)
		return
	}else {
		w.WriteHeader(http.StatusCreated)
	}
	w.Write([]byte("{\"success\":\"ok\"}"))
	w.Header().Set("Content-Type", "application/json")

}

func (handler Handler) GetPost(w http.ResponseWriter, r *http.Request){
	params := mux.Vars(r)
	postType := model.GetPostType(params["postType"])

	id, err := primitive.ObjectIDFromHex(params["id"])
	if err != nil || postType == model.NONE {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	var result model.Post
	result, err = handler.PostService.ReadPost(id, postType)
	if  err == mongo.ErrNoDocuments  {
		w.WriteHeader(http.StatusNotFound)
		return
	}
	if err != nil{
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(&result)
}

func (handler Handler) DeletePost(w http.ResponseWriter, r *http.Request){
	params := mux.Vars(r)
	postType := model.GetPostType(params["postType"])
	id, err := primitive.ObjectIDFromHex(params["id"])
	if err != nil || postType == model.NONE {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	err = handler.PostService.DeletePost(id, postType)
	if  err == mongo.ErrNoDocuments  {
		w.WriteHeader(http.StatusNotFound)
		return
	}
	if err != nil{
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json")
	w.Write([]byte("{\"success\":\"ok\"}"))
}

func (handler Handler) UpdatePost(w http.ResponseWriter, r *http.Request){
	params := mux.Vars(r)
	postType := model.GetPostType(params["postType"])
	id, err := primitive.ObjectIDFromHex(params["id"])
	var postDto dto.PostDto
	err = json.NewDecoder(r.Body).Decode(&postDto)
	if err != nil || postType == model.NONE{
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	err = handler.PostService.UpadtePost(id,postType,postDto)
	if  err == mongo.ErrNoDocuments  {
		w.WriteHeader(http.StatusNotFound)
		return
	}
	if err != nil{
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json")
	w.Write([]byte("{\"success\":\"ok\"}"))
}




