package handler

import (
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	uuid "github.com/nu7hatch/gouuid"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"io"
	"mime/multipart"
	"net/http"
	"nistagram/post/dto"
	"nistagram/post/model"
	"nistagram/post/service"
	"nistagram/util"
	"os"
	"strconv"
	"strings"
)

type Handler struct {
	PostService *service.PostService
}

func (handler *Handler) GetAll(w http.ResponseWriter, r *http.Request){
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

func (handler Handler) GetPublic(w http.ResponseWriter, r *http.Request){
	result := handler.PostService.GetPublic()

	js, err := json.Marshal(result)
	if err != nil{
		w.Write([]byte("{\"success\":\"error\"}"))
	}
	w.Write(js)

	//w.Write([]byte("{\"success\":\"ok\"}"))
	w.Header().Set("Content-Type", "application/json")
}

func (handler *Handler) Create(w http.ResponseWriter, r *http.Request){

	fmt.Println("In function create..")

	err := r.ParseMultipartForm(0)

	if err != nil{
		fmt.Println(err)
		w.Write([]byte("{\"success\":\"error\"}"))
		return
	}

	i := 0
	var files []*multipart.FileHeader
	for {
		file := r.MultipartForm.File["file"+strconv.Itoa(i)]
		if len(file) > 0 {
			files = append(files, file[0])
			i++
		}else{
			break
		}
	}

	var mediaNames []string
	for i:=0;i<len(files);i++{
		file,err := files[i].Open()
		if err != nil{
			fmt.Println(err)
			w.Write([]byte("{\"success\":\"error\"}"))
			return
		}
		uid, err := uuid.NewV4()
		if err != nil{
			fmt.Println(err)
		}
		fn := strings.Split(files[i].Filename, ".")
		mediaNames = append(mediaNames, uid.String() + "." + fn[1])
		f, err := os.OpenFile("../../nistagramstaticdata/data/" + uid.String() + "." + fn[1], os.O_WRONLY|os.O_CREATE, 0666)
		if err != nil{
			fmt.Println(err)
			w.Write([]byte("{\"success\":\"error\"}"))
			return
		}
		io.Copy(f, file)
		f.Close()
		file.Close()
	}

	var postDto dto.PostDto
	data := r.MultipartForm.Value["data"]
	fmt.Println(data)
	err = json.Unmarshal([]byte(data[0]), &postDto)
	if err != nil{
		fmt.Println(err)
		w.Write([]byte("{\"success\":\"error\"}"))
		return
	}

	postType := model.GetPostType(postDto.PostType)

	if postType == model.NONE{
		w.Write([]byte("{\"success\":\"error\"}"))
		return
	}

	err = handler.PostService.CreatePost(postType,postDto, mediaNames )

	if err != nil{
		fmt.Println(err)
		w.Write([]byte("{\"success\":\"error\"}"))
		return
	}else {
		w.WriteHeader(http.StatusCreated)
	}

	w.Write([]byte("{\"success\":\"ok\"}"))
	w.Header().Set("Content-Type", "application/json")

}

func (handler *Handler) GetPost(w http.ResponseWriter, r *http.Request){
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

func (handler *Handler) DeletePost(w http.ResponseWriter, r *http.Request){
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

func (handler *Handler) UpdatePost(w http.ResponseWriter, r *http.Request){
	params := mux.Vars(r)
	postType := model.GetPostType(params["postType"])
	id, err := primitive.ObjectIDFromHex(params["id"])
	var postDto dto.PostDto
	err = json.NewDecoder(r.Body).Decode(&postDto)
	if err != nil || postType == model.NONE{
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	err = handler.PostService.UpdatePost(id,postType,postDto)
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

func (handler *Handler) DeleteUserPosts (w http.ResponseWriter, r *http.Request){
	switch err := handler.PostService.DeleteUserPosts(util.GetLoggedUserIDFromToken(r)) ; err{
	case mongo.ErrNoDocuments:
		w.WriteHeader(http.StatusNotFound)
	case nil :
		w.WriteHeader(http.StatusOK)
	default:
		w.WriteHeader(http.StatusInternalServerError)
	}
}

func (handler *Handler) ChangeUsername(w http.ResponseWriter, r *http.Request) {
	type data struct { Username string `json:"username"` }
	var input data
	err := json.NewDecoder(r.Body).Decode(&input)

	if err != nil{
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	switch err =  handler.PostService.ChangeUsername(util.GetLoggedUserIDFromToken(r),input.Username) ; err{
	case mongo.ErrNoDocuments:
		w.WriteHeader(http.StatusNotFound)
	case nil :
		w.WriteHeader(http.StatusOK)
	default:
		w.WriteHeader(http.StatusInternalServerError)
	}
}

func (handler *Handler) ChangePrivacy (w http.ResponseWriter, r *http.Request) {
	type data struct { IsPrivate bool `json:"IsPrivate"` }
	var input data
	err := json.NewDecoder(r.Body).Decode(&input)

	if err != nil{
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	switch err := handler.PostService.ChangePrivacy(123, input.IsPrivate) ; err {
	case mongo.ErrNoDocuments:
		w.WriteHeader(http.StatusNotFound)
	case nil :
		w.WriteHeader(http.StatusOK)
	default:
		w.WriteHeader(http.StatusInternalServerError)
	}
}




