package handler

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/google/uuid"
	"github.com/gorilla/mux"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"golang.org/x/crypto/openpgp/errors"
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

func (handler Handler) GetPublic(w http.ResponseWriter, r *http.Request) {
	span := util.Tracer.StartSpanFromRequest("GetPublic-handler", r)
	defer util.Tracer.FinishSpan(span)

	ctx := util.Tracer.ContextWithSpan(context.Background(), span)
	result, err := handler.PostService.GetPublic(ctx, util.GetLoggedUserIDFromToken(r))
	if err != nil {
		util.Tracer.LogError(span, err)
		fmt.Println(err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	if js, err := json.Marshal(result); err != nil {
		util.Tracer.LogError(span, err)
		w.WriteHeader(http.StatusInternalServerError)
	} else {
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write(js)
	}
	w.Header().Set("Content-Type", "application/json")
}

func (handler Handler) GetMyPosts(w http.ResponseWriter, r *http.Request) {
	span := util.Tracer.StartSpanFromRequest("GetMyPosts-handler", r)
	defer util.Tracer.FinishSpan(span)

	ctx := util.Tracer.ContextWithSpan(context.Background(), span)

	w.Header().Set("Content-Type", "application/json")

	loggedUserID := util.GetLoggedUserIDFromToken(r)
	if loggedUserID == 0 {
		util.Tracer.LogError(span, fmt.Errorf("user is not logged in"))
		http.Error(w, "user is not logged in", http.StatusForbidden)
		return
	}

	result, err := handler.PostService.GetMyPosts(ctx, loggedUserID)
	if err != nil {
		util.Tracer.LogError(span, err)
		fmt.Println(err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	js, err := json.Marshal(result)
	if err != nil {
		util.Tracer.LogError(span, err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
	_, _ = w.Write(js)
}

func (handler Handler) GetPostsForHomePage(w http.ResponseWriter, r *http.Request) {
	span := util.Tracer.StartSpanFromRequest("GetPostsForHomePage-handler", r)
	defer util.Tracer.FinishSpan(span)

	ctx := util.Tracer.ContextWithSpan(context.Background(), span)

	w.Header().Set("Content-Type", "application/json")

	loggedUserID := util.GetLoggedUserIDFromToken(r)
	if loggedUserID == 0 {
		util.Tracer.LogError(span, fmt.Errorf("user is not logged in"))
		http.Error(w, "user is not logged in", http.StatusForbidden)
		return
	}

	followingProfiles, err := getFollowingProfiles(ctx, loggedUserID)
	if err != nil {
		util.Tracer.LogError(span, err)
		fmt.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	result, err := handler.PostService.GetPostsForHomePage(ctx, followingProfiles, loggedUserID)
	if err != nil {
		util.Tracer.LogError(span, err)
		fmt.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	if js, err := json.Marshal(result); err != nil {
		util.Tracer.LogError(span, err)
		w.WriteHeader(http.StatusInternalServerError)
	} else {
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write(js)
	}

}

func (handler Handler) GetProfilesPosts(w http.ResponseWriter, r *http.Request) {
	span := util.Tracer.StartSpanFromRequest("GetProfilesPosts-handler", r)
	defer util.Tracer.FinishSpan(span)

	ctx := util.Tracer.ContextWithSpan(context.Background(), span)
	params := mux.Vars(r)
	targetUsername := template.HTMLEscapeString(params["username"])
	loggedUserId := util.GetLoggedUserIDFromToken(r)
	var followingProfiles []util.FollowingProfileDTO
	var err error
	if loggedUserId != 0 {
		followingProfiles, err = getFollowingProfiles(ctx, loggedUserId)
		if err != nil {
			util.Tracer.LogError(span, err)
			fmt.Println(err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
	}

	result, err := handler.PostService.GetProfilesPosts(ctx, followingProfiles, targetUsername, loggedUserId)
	if err != nil {
		util.Tracer.LogError(span, err)
		fmt.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	if js, err := json.Marshal(result); err != nil {
		util.Tracer.LogError(span, err)
		w.WriteHeader(http.StatusInternalServerError)
	} else {
		_, _ = w.Write(js)
	}

	w.Header().Set("Content-Type", "application/json")
}

func (handler *Handler) SearchPublicByLocation(w http.ResponseWriter, r *http.Request) {
	span := util.Tracer.StartSpanFromRequest("SearchPublicByLocation-handler", r)
	defer util.Tracer.FinishSpan(span)

	ctx := util.Tracer.ContextWithSpan(context.Background(), span)
	params := mux.Vars(r)
	location := template.HTMLEscapeString(params["value"])

	result, err := handler.PostService.GetPublicPostByLocation(ctx, location, util.GetLoggedUserIDFromToken(r))
	if err != nil {
		util.Tracer.LogError(span, err)
		fmt.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	if err = json.NewEncoder(w).Encode(&result); err != nil {
		util.Tracer.LogError(span, err)
		w.WriteHeader(http.StatusInternalServerError)
	} else {
		w.WriteHeader(http.StatusOK)
	}
	w.Header().Set("Content-Type", "application/json")
}

func (handler *Handler) SearchPublicByHashTag(w http.ResponseWriter, r *http.Request) {
	span := util.Tracer.StartSpanFromRequest("SearchPublicByHashTag-handler", r)
	defer util.Tracer.FinishSpan(span)

	ctx := util.Tracer.ContextWithSpan(context.Background(), span)
	params := mux.Vars(r)
	hashTag := template.HTMLEscapeString(params["value"])

	result, err := handler.PostService.GetPublicPostByHashTag(ctx, hashTag, util.GetLoggedUserIDFromToken(r))
	if err != nil {
		util.Tracer.LogError(span, err)
		fmt.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	if err = json.NewEncoder(w).Encode(&result); err != nil {
		util.Tracer.LogError(span, err)
		w.WriteHeader(http.StatusInternalServerError)
	} else {
		w.WriteHeader(http.StatusOK)
	}
	w.Header().Set("Content-Type", "application/json")
}

func (handler *Handler) Create(w http.ResponseWriter, r *http.Request) {
	span := util.Tracer.StartSpanFromRequest("Create-handler", r)
	defer util.Tracer.FinishSpan(span)

	ctx := util.Tracer.ContextWithSpan(context.Background(), span)
	profileId := util.GetLoggedUserIDFromToken(r)
	methodPath := "nistagram/post/handler.Create"
	if err := r.ParseMultipartForm(0); err != nil {
		util.Tracer.LogError(span, err)
		util.Logging(util.ERROR, methodPath, "", err.Error(), "post")
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	var postDto dto.PostDto
	data := r.MultipartForm.Value["data"]

	if err := json.Unmarshal([]byte(data[0]), &postDto); err != nil {
		util.Tracer.LogError(span, err)
		util.Logging(util.ERROR, methodPath, "", err.Error(), "post")
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	var files []*multipart.FileHeader
	for i := 0; ; i++ {
		if file := r.MultipartForm.File["file"+strconv.Itoa(i)]; len(file) > 0 {
			files = append(files, file[0])
		} else {
			break
		}
	}

	mediaNames := generateMediaNames(files)

	fmt.Println(postDto)

	switch err := handler.createPost(ctx, profileId, postDto, mediaNames); err {
	case nil:
		w.WriteHeader(http.StatusCreated)
		util.Logging(util.SUCCESS, methodPath, util.GetIPAddress(r), "Success in creating post, "+util.GetLoggingStringFromID(profileId), "post")
	case errors.InvalidArgumentError("input value"):
		util.Tracer.LogError(span, err)
		util.Logging(util.ERROR, methodPath, util.GetIPAddress(r), "Wrong post type", "post")
		w.WriteHeader(http.StatusBadRequest)
		return
	default:
		util.Tracer.LogError(span, err)
		fmt.Println(err)
		util.Logging(util.ERROR, methodPath, util.GetIPAddress(r), "Wrong post type", "post")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	if err := insertIntoFiles(ctx, files, mediaNames); err != nil {
		util.Tracer.LogError(span, err)
		util.Logging(util.ERROR, methodPath, "", err.Error(), "post")
		w.WriteHeader(http.StatusInternalServerError)
	}

}

func (handler *Handler) GetPost(w http.ResponseWriter, r *http.Request) {
	span := util.Tracer.StartSpanFromRequest("GetPost-handler", r)
	defer util.Tracer.FinishSpan(span)

	ctx := util.Tracer.ContextWithSpan(context.Background(), span)
	params := mux.Vars(r)

	id, err := primitive.ObjectIDFromHex(params["id"])
	if err != nil {
		util.Tracer.LogError(span, err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	switch result, err := handler.PostService.ReadPost(ctx, id); err {
	case mongo.ErrNoDocuments:
		util.Tracer.LogError(span, err)
		w.WriteHeader(http.StatusNotFound)
	case nil:
		w.WriteHeader(http.StatusOK)
		err = json.NewEncoder(w).Encode(&result)
		if err != nil {
			util.Tracer.LogError(span, err)
			return
		}
	default:
		util.Tracer.LogError(span, err)
		w.WriteHeader(http.StatusInternalServerError)
	}
	w.Header().Set("Content-Type", "application/json")
}

func (handler *Handler) DeletePost(w http.ResponseWriter, r *http.Request) {
	span := util.Tracer.StartSpanFromRequest("DeletePost-handler", r)
	defer util.Tracer.FinishSpan(span)

	ctx := util.Tracer.ContextWithSpan(context.Background(), span)
	params := mux.Vars(r)
	id, err := primitive.ObjectIDFromHex(params["id"])
	if err != nil {
		util.Tracer.LogError(span, err)
		fmt.Println(err)
		w.WriteHeader(http.StatusBadRequest)
		_, _ = w.Write([]byte("{\"message\":\"error\"}"))
		return
	}

	err = handler.PostService.DeletePost(ctx, id)
	if err != nil {
		util.Tracer.LogError(span, err)
		fmt.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		_, _ = w.Write([]byte("{\"message\":\"error\"}"))
		return
	}
	w.WriteHeader(http.StatusOK)
	_, _ = w.Write([]byte("{\"message\":\"ok\"}"))
	w.Header().Set("Content-Type", "application/json")
}

func (handler *Handler) UpdatePost(w http.ResponseWriter, r *http.Request) {
	span := util.Tracer.StartSpanFromRequest("UpdatePost-handler", r)
	defer util.Tracer.FinishSpan(span)

	ctx := util.Tracer.ContextWithSpan(context.Background(), span)
	params := mux.Vars(r)
	id, err := primitive.ObjectIDFromHex(params["id"])
	if err != nil {
		util.Tracer.LogError(span, err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	var postDto dto.PostDto
	if err = json.NewDecoder(r.Body).Decode(&postDto); err != nil {
		util.Tracer.LogError(span, err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	postDto = safePostDto(postDto)

	switch err = handler.PostService.UpdatePost(ctx, id, postDto); err {
	case mongo.ErrNoDocuments:
		util.Tracer.LogError(span, err)
		w.WriteHeader(http.StatusNotFound)
	case nil:
		w.WriteHeader(http.StatusOK)
	default:
		util.Tracer.LogError(span, err)
		w.WriteHeader(http.StatusInternalServerError)
	}
}

func (handler *Handler) DeleteUserPosts(w http.ResponseWriter, r *http.Request) {
	span := util.Tracer.StartSpanFromRequest("DeleteUserPosts-handler", r)
	defer util.Tracer.FinishSpan(span)

	ctx := util.Tracer.ContextWithSpan(context.Background(), span)
	params := mux.Vars(r)
	profileId := util.String2Uint(params["id"])
	switch err := handler.PostService.DeleteUserPosts(ctx, profileId); err {
	case mongo.ErrNoDocuments:
		util.Tracer.LogError(span, err)
		w.WriteHeader(http.StatusNotFound)
	case nil:
		w.WriteHeader(http.StatusOK)
	default:
		util.Tracer.LogError(span, err)
		w.WriteHeader(http.StatusInternalServerError)
	}
}

func (handler *Handler) ChangeUsername(w http.ResponseWriter, r *http.Request) {
	span := util.Tracer.StartSpanFromRequest("ChangeUsername-handler", r)
	defer util.Tracer.FinishSpan(span)

	ctx := util.Tracer.ContextWithSpan(context.Background(), span)
	params := mux.Vars(r)
	publisherId := util.String2Uint(params["loggedUserId"])

	type data struct {
		Username string `json:"username"`
	}
	var input data
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		util.Tracer.LogError(span, err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	input.Username = template.HTMLEscapeString(input.Username)

	switch err := handler.PostService.ChangeUsername(ctx, publisherId, input.Username); err {
	case mongo.ErrNoDocuments:
		util.Tracer.LogError(span, err)
		w.WriteHeader(http.StatusNotFound)
	case nil:
		w.WriteHeader(http.StatusOK)
	default:
		util.Tracer.LogError(span, err)
		w.WriteHeader(http.StatusInternalServerError)
	}
}

func (handler *Handler) ChangePrivacy(w http.ResponseWriter, r *http.Request) {
	span := util.Tracer.StartSpanFromRequest("ChangePrivacy-handler", r)
	defer util.Tracer.FinishSpan(span)

	ctx := util.Tracer.ContextWithSpan(context.Background(), span)
	type data struct {
		IsPrivate bool `json:"IsPrivate"`
	}
	params := mux.Vars(r)
	publisherId := util.String2Uint(params["loggedUserId"])

	var input data
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		util.Tracer.LogError(span, err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	switch err := handler.PostService.ChangePrivacy(ctx, publisherId, input.IsPrivate); err {
	case mongo.ErrNoDocuments:
		util.Tracer.LogError(span, err)
		w.WriteHeader(http.StatusNotFound)
	case nil:
		w.WriteHeader(http.StatusOK)
	default:
		util.Tracer.LogError(span, err)
		w.WriteHeader(http.StatusInternalServerError)
	}
}

func (handler *Handler) GetPosts(w http.ResponseWriter, r *http.Request) {
	span := util.Tracer.StartSpanFromRequest("GetPosts-handler", r)
	defer util.Tracer.FinishSpan(span)

	ctx := util.Tracer.ContextWithSpan(context.Background(), span)
	type data struct {
		Ids []string `json:"ids"`
	}
	var input data
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		util.Tracer.LogError(span, err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	var posts []model.Post
	for _, value := range input.Ids {
		postID, err := primitive.ObjectIDFromHex(value)
		if err != nil {
			util.Tracer.LogError(span, err)
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		switch result, err := handler.PostService.ReadPost(ctx, postID); err {
		case mongo.ErrNoDocuments:
			util.Tracer.LogError(span, fmt.Errorf("post not found"))
			continue //escaping deleted posts
		case nil:
			posts = append(posts, result)
		default:
			util.Tracer.LogError(span, err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
	}
	if js, err := json.Marshal(posts); err != nil {
		util.Tracer.LogError(span, err)
		w.WriteHeader(http.StatusInternalServerError)
	} else {
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write(js)
	}
}

func (handler *Handler) MakeCampaign(w http.ResponseWriter, r *http.Request) {
	span := util.Tracer.StartSpanFromRequest("MakeCampaign-handler", r)
	defer util.Tracer.FinishSpan(span)

	ctx := util.Tracer.ContextWithSpan(context.Background(), span)
	params := mux.Vars(r)
	err := handler.PostService.MakeCampaign(ctx, params["id"], util.String2Uint(params["agentID"]))
	if err != nil {
		util.Tracer.LogError(span, err)
		fmt.Println(err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	w.WriteHeader(http.StatusOK)
	_, _ = w.Write([]byte("{\"message\":\"ok\"}"))
	w.Header().Set("Content-Type", "application/json")
}

func (handler Handler) GetMediaById(w http.ResponseWriter, r *http.Request) {
	span := util.Tracer.StartSpanFromRequest("GetMediaById-handler", r)
	defer util.Tracer.FinishSpan(span)

	ctx := util.Tracer.ContextWithSpan(context.Background(), span)

	params := mux.Vars(r)
	mediaId := template.HTMLEscapeString(params["id"])

	result, err := handler.PostService.GetMediaById(ctx, mediaId)
	if err != nil {
		util.Tracer.LogError(span, err)
		fmt.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	err = json.NewEncoder(w).Encode(result)
	if err != nil {
		util.Tracer.LogError(span, err)
		fmt.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json")
}

func (handler *Handler) GetPostById(w http.ResponseWriter, r *http.Request) {
	span := util.Tracer.StartSpanFromRequest("GetPostById-handler", r)
	defer util.Tracer.FinishSpan(span)

	ctx := util.Tracer.ContextWithSpan(context.Background(), span)

	params := mux.Vars(r)
	id, err := primitive.ObjectIDFromHex(params["id"])
	if err != nil {
		util.Tracer.LogError(span, err)
		fmt.Println(err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	loggedUserId := util.GetLoggedUserIDFromToken(r)
	followingProfiles,err := getFollowingProfiles(ctx, loggedUserId)
	if err != nil{
		util.Tracer.LogError(span, err)
		fmt.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	result, err := handler.PostService.ReadPost(ctx, id)
	if err != nil {
		util.Tracer.LogError(span, err)
	}

	if result.IsPrivate {
		if !util.IsFollowed(followingProfiles, result.PublisherId) {
			fmt.Println("Publisher is not followed by logged user!")
			w.WriteHeader(http.StatusBadRequest)
			return
		}
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(result)
	w.Header().Set("Content-Type", "application/json")
}

func safePostDto(postDto dto.PostDto) dto.PostDto {
	postDto.Description = template.HTMLEscapeString(postDto.Description)
	postDto.HashTags = template.HTMLEscapeString(postDto.HashTags)
	postDto.Location = template.HTMLEscapeString(postDto.Location)
	return postDto
}

func getFollowingProfiles(ctx context.Context, loggedUserId uint) ([]util.FollowingProfileDTO, error) {
	span := util.Tracer.StartSpanFromContext(ctx, "getFollowingProfiles-handler")
	defer util.Tracer.FinishSpan(span)

	nextCtx := util.Tracer.ContextWithSpan(ctx, span)

	connHost, connPort := util.GetConnectionHostAndPort()
	resp, err := util.CrossServiceRequest(nextCtx, http.MethodGet,
		util.GetCrossServiceProtocol()+"://"+connHost+":"+connPort+"/connection/following/show/"+util.Uint2String(loggedUserId),
		nil, map[string]string{})

	if err != nil {
		util.Tracer.LogError(span, err)
		return nil, err
	}

	var followingProfiles []util.FollowingProfileDTO

	err = json.NewDecoder(resp.Body).Decode(&followingProfiles)
	if err != nil {
		util.Tracer.LogError(span, err)
		return nil, err
	}

	return followingProfiles, err
}

func getProfileByProfileId(profileId uint) (*http.Response, error) {
	profileHost, profilePort := util.GetProfileHostAndPort()
	resp, err := util.CrossServiceRequest(context.Background(), http.MethodGet,
		util.GetCrossServiceProtocol()+"://"+profileHost+":"+profilePort+"/get-by-id/"+util.Uint2String(profileId),
		nil, map[string]string{})
	return resp, err
}

func (handler *Handler) createPost(ctx context.Context, profileId uint, postDto dto.PostDto, mediaNames []string) error {
	span := util.Tracer.StartSpanFromContext(ctx, "createPost-handler")
	defer util.Tracer.FinishSpan(span)

	nextCtx := util.Tracer.ContextWithSpan(ctx, span)

	postType := model.GetPostType(postDto.PostType)
	if postType == model.NONE {
		util.Tracer.LogError(span, fmt.Errorf("invalid post type value"))
		return errors.InvalidArgumentError("input value")
	}

	var profile dto.ProfileDto
	if resp, err := getProfileByProfileId(profileId); err != nil {
		util.Tracer.LogError(span, err)
		return err
	} else {
		body, _ := io.ReadAll(resp.Body)
		defer func(Body io.ReadCloser) {
			_ = Body.Close()
		}(resp.Body)
		if err := json.Unmarshal(body, &profile); err != nil {
			util.Tracer.LogError(span, err)
			return err
		}
	}

	profile.ProfileId = profileId
	if err := handler.PostService.CreatePost(nextCtx, postType, postDto, mediaNames, profile); err != nil {
		util.Tracer.LogError(span, err)
		return err
	}
	return nil
}

func generateMediaNames(files []*multipart.FileHeader) []string {
	var mediaNames []string
	for i := 0; i < len(files); i++ {
		uid := uuid.NewString()
		fn := strings.Split(files[i].Filename, ".")
		mediaNames = append(mediaNames, uid+"."+fn[1])
	}
	return mediaNames
}

func insertIntoFiles(ctx context.Context, files []*multipart.FileHeader, mediaNames []string) error {
	span := util.Tracer.StartSpanFromContext(ctx, "insertIntoFiles-handler")
	defer util.Tracer.FinishSpan(span)

	for i := 0; i < len(files); i++ {
		file, err := files[i].Open()
		if err != nil {
			util.Tracer.LogError(span, err)
			return err
		}

		f, err := os.OpenFile("../../nistagramstaticdata/data/"+mediaNames[i], os.O_WRONLY|os.O_CREATE, 0666)
		if err != nil {
			util.Tracer.LogError(span, err)
			return err
		}
		if _, err = io.Copy(f, file); err != nil {
			util.Tracer.LogError(span, err)
			return err
		}
		if err = f.Close(); err != nil {
			util.Tracer.LogError(span, err)
			return err
		}
		if err = file.Close(); err != nil {
			util.Tracer.LogError(span, err)
			return err
		}
	}
	return nil
}
