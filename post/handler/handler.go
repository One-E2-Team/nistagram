package handler

import (
	"encoding/json"
	"github.com/gorilla/mux"
	uuid "github.com/nu7hatch/gouuid"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"io"
	"net/http"
	"nistagram/post/dto"
	"nistagram/post/model"
	"nistagram/post/service"
	"os"
)

type Handler struct {
	PostService *service.PostService
}

func (handler Handler) GetAll(w http.ResponseWriter, r *http.Request){
	(w).Header().Set("Access-Control-Allow-Origin", "*")
	result := handler.PostService.GetAll()
	//json.NewEncoder(w).Encode(&result)

	js, err := json.Marshal(result)
	if err != nil{
		w.Write([]byte("{\"success\":\"error\"}"))
	}
	w.Write(js)

	//w.Write([]byte("{\"success\":\"ok\"}"))
	w.Header().Set("Content-Type", "application/json")
}

func (handler *Handler) Create(w http.ResponseWriter, r *http.Request){
	(w).Header().Set("Access-Control-Allow-Origin", "*")

	err := r.ParseMultipartForm(0)

	if err != nil{
		w.Write([]byte("{\"success\":\"error\"}"))
		return
	}

	files := r.MultipartForm.File["files"]

	var mediaNames []string

	for i:=0;i<len(files);i++{
		file,err := files[i].Open()
		if err != nil{
			w.Write([]byte("{\"success\":\"error\"}"))
			return
		}
		uid, err := uuid.NewV4()
		mediaNames = append(mediaNames, uid.String() + ".jpg")
		f, err := os.OpenFile("../../nistagramstaticdata/data/" + uid.String() + ".jpg", os.O_WRONLY|os.O_CREATE, 0666)
		if err != nil{
			w.Write([]byte("{\"success\":\"error\"}"))
			return
		}
		io.Copy(f, file)
		f.Close()
		file.Close()
	}

	var postDto dto.PostDto
	data := r.FormValue("data")
	err = json.Unmarshal([]byte(data), &postDto)
	if err != nil{
		w.Write([]byte("{\"success\":\"error\"}"))
		return
	}

	postType := model.GetPostType(postDto.PostType)

	if postType == model.NONE{
		w.Write([]byte("{\"success\":\"error\"}"))
		return
	}

	err = handler.PostService.CreatePost(postType,postDto, mediaNames)

	if err != nil{
		w.Write([]byte("{\"success\":\"error\"}"))
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




