package handler

import (
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"github.com/google/uuid"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"html/template"
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

func (handler *Handler) GetAll(w http.ResponseWriter, _ *http.Request) {
	result := handler.PostService.GetAll()
	//json.NewEncoder(w).Encode(&result)

	js, err := json.Marshal(result)
	if err != nil {
		_, _ = w.Write([]byte("{\"success\":\"error\"}"))
	}
	_, _ = w.Write(js)
	w.Header().Set("Content-Type", "application/json")
}

func (handler Handler) GetPublic(w http.ResponseWriter, _ *http.Request) {
	result := handler.PostService.GetPublic()

	js, err := json.Marshal(result)
	if err != nil {
		_, _ = w.Write([]byte("{\"success\":\"error\"}"))
	}
	_, _ = w.Write(js)
	w.Header().Set("Content-Type", "application/json")
}

func (handler Handler) GetMyPosts(w http.ResponseWriter, r *http.Request) {

	loggedUserId := util.GetLoggedUserIDFromToken(r)
	if loggedUserId == 0 {
		fmt.Println("User is not logged in..")
		_, _ = w.Write([]byte("{\"success\":\"error\"}"))
		return
	}

	result := handler.PostService.GetMyPosts(loggedUserId)

	js, err := json.Marshal(result)
	if err != nil {
		_, _ = w.Write([]byte("{\"success\":\"error\"}"))
	}
	_, err = w.Write(js)
	if err != nil {
		return
	}

	//w.Write([]byte("{\"success\":\"ok\"}"))
	w.Header().Set("Content-Type", "application/json")
}

func (handler Handler) GetPostsForHomePage(w http.ResponseWriter, r *http.Request) {

	loggedUserId := util.GetLoggedUserIDFromToken(r)
	if loggedUserId == 0 {
		fmt.Println("User is not logged in..")
		_, _ = w.Write([]byte("{\"success\":\"error\"}"))
		return
	}

	connHost, connPort := util.GetConnectionHostAndPort()
	resp, err := http.Get(util.CrossServiceProtocol + "://" + connHost + ":" + connPort + "/connection/following/show/" + util.Uint2String(loggedUserId))

	if err != nil {
		fmt.Println(err)
	}
	var followingProfiles []uint
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println("Body: ", body)
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {

		}
	}(resp.Body)
	err = json.Unmarshal(body, &followingProfiles)

	if err != nil {
		fmt.Println(err)
	}

	result := handler.PostService.GetPostsForHomePage(followingProfiles)

	js, err := json.Marshal(result)
	if err != nil {
		_, _ = w.Write([]byte("{\"success\":\"error\"}"))
	}
	_, _ = w.Write(js)
	w.Header().Set("Content-Type", "application/json")
}

func (handler Handler) GetProfilesPosts(w http.ResponseWriter, r *http.Request) {

	params := mux.Vars(r)
	targetUsername := template.HTMLEscapeString(params["username"])

	var followingProfiles []uint

	loggedUserId := util.GetLoggedUserIDFromToken(r)
	if loggedUserId != 0 {
		connHost, connPort := util.GetConnectionHostAndPort()
		resp, err := http.Get(util.CrossServiceProtocol + "://" + connHost + ":" + connPort + "/connection/following/show/" + util.Uint2String(loggedUserId))

		if err != nil {
			fmt.Println(err)
		}

		body, err := io.ReadAll(resp.Body)
		if err != nil {
			fmt.Println(err)
		}
		fmt.Println("Body: ", body)
		defer func(Body io.ReadCloser) {
			err := Body.Close()
			if err != nil {

			}
		}(resp.Body)
		err = json.Unmarshal(body, &followingProfiles)

		if err != nil {
			fmt.Println(err)
		}
	}

	result := handler.PostService.GetProfilesPosts(followingProfiles, targetUsername)

	js, err := json.Marshal(result)
	if err != nil {
		_, _ = w.Write([]byte("{\"success\":\"error\"}"))
	}
	_, _ = w.Write(js)
	w.Header().Set("Content-Type", "application/json")
}

func (handler *Handler) SearchPublicByLocation(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	location := template.HTMLEscapeString(params["value"])

	var result []model.Post
	result = handler.PostService.GetPublicPostByLocation(location)

	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json")
	err := json.NewEncoder(w).Encode(&result)
	if err != nil {
		return
	}
}

func (handler *Handler) SearchPublicByHashTag(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	hashTag := template.HTMLEscapeString(params["value"])

	var result []model.Post
	result = handler.PostService.GetPublicPostByHashTag(hashTag)

	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json")
	err := json.NewEncoder(w).Encode(&result)
	if err != nil {
		return
	}
}

func (handler *Handler) Create(w http.ResponseWriter, r *http.Request) {
	profileId := util.GetLoggedUserIDFromToken(r)
	methodPath := "nistagram/post/handler.Create"
	err := r.ParseMultipartForm(0)
	if err != nil {
		util.Logging(util.ERROR, methodPath, "", err.Error(), "post")
		_, _ = w.Write([]byte("{\"success\":\"error\"}"))
		return
	}

	var postDto dto.PostDto
	data := r.MultipartForm.Value["data"]
	err = json.Unmarshal([]byte(data[0]), &postDto)
	if err != nil {
		util.Logging(util.ERROR, methodPath, "", err.Error(), "post")
		_, _ = w.Write([]byte("{\"success\":\"error\"}"))
		return
	}

	postType := model.GetPostType(postDto.PostType)
	if postType == model.NONE {
		util.Logging(util.ERROR, methodPath, util.GetIPAddress(r), "Wrong post type", "post")
		_, _ = w.Write([]byte("{\"success\":\"error\"}"))
		return
	}

	i := 0
	var files []*multipart.FileHeader
	for {
		file := r.MultipartForm.File["file"+strconv.Itoa(i)]
		if len(file) > 0 {
			files = append(files, file[0])
			i++
		} else {
			break
		}
	}

	var mediaNames []string
	for i := 0; i < len(files); i++ {
		file, err := files[i].Open()
		if err != nil {
			util.Logging(util.ERROR, methodPath, "", err.Error(), "post")
			_, _ = w.Write([]byte("{\"success\":\"error\"}"))
			return
		}
		uid := uuid.NewString()

		fn := strings.Split(files[i].Filename, ".")
		mediaNames = append(mediaNames, uid + "."+fn[1])
		f, err := os.OpenFile("../../nistagramstaticdata/data/"+uid + "."+fn[1], os.O_WRONLY|os.O_CREATE, 0666)
		if err != nil {
			util.Logging(util.ERROR, methodPath, "", err.Error(), "post")
			_, _ = w.Write([]byte("{\"success\":\"error\"}"))
			return
		}
		_, err = io.Copy(f, file)
		if err != nil {
			return
		}
		err = f.Close()
		if err != nil {
			return
		}
		err = file.Close()
		if err != nil {
			return
		}
	}

	profileHost, profilePort := util.GetProfileHostAndPort()
	resp, err := http.Get(util.CrossServiceProtocol + "://" + profileHost + ":" + profilePort + "/get-by-id/" + util.Uint2String(profileId))
	if err != nil {
		util.Logging(util.ERROR, methodPath, "", err.Error(), "post")
		_, _ = w.Write([]byte("{\"success\":\"error\"}"))
		return
	}
	var profile dto.ProfileDto
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		util.Logging(util.ERROR, methodPath, "", err.Error(), "post")
		_, _ = w.Write([]byte("{\"success\":\"error\"}"))
		return
	}

	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {

		}
	}(resp.Body)
	err = json.Unmarshal(body, &profile)
	if err != nil {
		util.Logging(util.ERROR, methodPath, "", err.Error(), "post")
		_, _ = w.Write([]byte("{\"success\":\"error\"}"))
		return
	}
	profile.ProfileId = profileId
	err = handler.PostService.CreatePost(postType, postDto, mediaNames, profile)
	if err != nil {
		util.Logging(util.ERROR, methodPath, "", err.Error(), "post")
		_, _ = w.Write([]byte("{\"success\":\"error\"}"))
		return
	} else {
		w.WriteHeader(http.StatusCreated)
	}
	util.Logging(util.SUCCESS, methodPath, util.GetIPAddress(r), "Success in creating post, "+util.GetLoggingStringFromID(profileId), "post")
	_, _ = w.Write([]byte("{\"success\":\"ok\"}"))
	w.Header().Set("Content-Type", "application/json")

}

func (handler *Handler) GetPost(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	postType := model.GetPostType(params["postType"])

	id, err := primitive.ObjectIDFromHex(params["id"])
	if err != nil || postType == model.NONE {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	var result model.Post
	result, err = handler.PostService.ReadPost(id, postType)
	if err == mongo.ErrNoDocuments {
		w.WriteHeader(http.StatusNotFound)
		return
	}
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json")
	err = json.NewEncoder(w).Encode(&result)
	if err != nil {
		return
	}
}

func (handler *Handler) DeletePost(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	postType := model.GetPostType(params["postType"])
	id, err := primitive.ObjectIDFromHex(params["id"])
	if err != nil || postType == model.NONE {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	err = handler.PostService.DeletePost(id, postType)
	if err == mongo.ErrNoDocuments {
		w.WriteHeader(http.StatusNotFound)
		return
	}
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json")
	_, _ = w.Write([]byte("{\"success\":\"ok\"}"))
}

func (handler *Handler) UpdatePost(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	postType := model.GetPostType(params["postType"])
	id, err := primitive.ObjectIDFromHex(params["id"])
	var postDto dto.PostDto
	err = json.NewDecoder(r.Body).Decode(&postDto)

	postDto = safePostDto(postDto)
	if err != nil || postType == model.NONE {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	err = handler.PostService.UpdatePost(id, postType, postDto)
	if err == mongo.ErrNoDocuments {
		w.WriteHeader(http.StatusNotFound)
		return
	}
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json")
	_, _ = w.Write([]byte("{\"success\":\"ok\"}"))
}

func (handler *Handler) DeleteUserPosts(w http.ResponseWriter, r *http.Request) {
	switch err := handler.PostService.DeleteUserPosts(util.GetLoggedUserIDFromToken(r)); err {
	case mongo.ErrNoDocuments:
		w.WriteHeader(http.StatusNotFound)
	case nil:
		w.WriteHeader(http.StatusOK)
	default:
		w.WriteHeader(http.StatusInternalServerError)
	}
}

func (handler *Handler) ChangeUsername(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	publisherId := util.String2Uint(params["loggedUserId"])

	type data struct {
		Username string `json:"username"`
	}
	var input data
	err := json.NewDecoder(r.Body).Decode(&input)

	input.Username = template.HTMLEscapeString(input.Username)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	switch err = handler.PostService.ChangeUsername(publisherId, input.Username); err {
	case mongo.ErrNoDocuments:
		w.WriteHeader(http.StatusNotFound)
	case nil:
		w.WriteHeader(http.StatusOK)
	default:
		w.WriteHeader(http.StatusInternalServerError)
	}
}

func (handler *Handler) ChangePrivacy(w http.ResponseWriter, r *http.Request) {
	type data struct {
		IsPrivate bool `json:"IsPrivate"`
	}
	params := mux.Vars(r)
	publisherId := util.String2Uint(params["loggedUserId"])

	var input data
	err := json.NewDecoder(r.Body).Decode(&input)

	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	switch err := handler.PostService.ChangePrivacy(publisherId, input.IsPrivate); err {
	case mongo.ErrNoDocuments:
		w.WriteHeader(http.StatusNotFound)
	case nil:
		w.WriteHeader(http.StatusOK)
	default:
		w.WriteHeader(http.StatusInternalServerError)
	}
}

func safePostDto(postDto dto.PostDto) dto.PostDto {
	postDto.Description = template.HTMLEscapeString(postDto.Description)
	taggedUsers := postDto.TaggedUsers
	for i := 0; i < len(postDto.TaggedUsers); i++ {
		taggedUsers[i] = template.HTMLEscapeString(postDto.TaggedUsers[i])
	}
	postDto.TaggedUsers = taggedUsers
	postDto.HashTags = template.HTMLEscapeString(postDto.HashTags)
	postDto.Location = template.HTMLEscapeString(postDto.Location)
	return postDto
}
