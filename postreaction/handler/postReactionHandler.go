package handler

import (
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"net/http"
	"nistagram/postreaction/dto"
	"nistagram/postreaction/model"
	"nistagram/postreaction/service"
	"nistagram/util"
	"strings"
)

type PostReactionHandler struct {
	PostReactionService *service.PostReactionService
}

func (handler *PostReactionHandler) ReactOnPost(w http.ResponseWriter, r *http.Request) {
	var reactionDTO dto.ReactionDTO
	err := json.NewDecoder(r.Body).Decode(&reactionDTO)
	if err != nil {
		fmt.Println(err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	reactionType := model.GetReactionType(reactionDTO.ReactionType)
	if reactionType == model.NONE {
		fmt.Println("Bad reaction type in request!")
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	err = handler.PostReactionService.ReactOnPost(reactionDTO.PostID, util.GetLoggedUserIDFromToken(r), reactionType)
	if err != nil {
		fmt.Println(err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	w.WriteHeader(http.StatusOK)
	_, _ = w.Write([]byte("{\"success\":\"ok\"}"))
	w.Header().Set("Content-Type", "application/json")
}

func (handler *PostReactionHandler) DeleteReaction(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	postID := vars["postID"]
	err := handler.PostReactionService.DeleteReaction(postID, util.GetLoggedUserIDFromToken(r))
	if err != nil {
		fmt.Println(err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	w.WriteHeader(http.StatusOK)
	_, _ = w.Write([]byte("{\"success\":\"ok\"}"))
	w.Header().Set("Content-Type", "application/json")
}

func (handler *PostReactionHandler) CommentPost(w http.ResponseWriter, r *http.Request) {
	var commentDTO dto.CommentDTO
	err := json.NewDecoder(r.Body).Decode(&commentDTO)
	if err != nil {
		fmt.Println(err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	err = handler.PostReactionService.CommentPost(commentDTO, util.GetLoggedUserIDFromToken(r))
	if err != nil {
		fmt.Println(err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	w.WriteHeader(http.StatusOK)
	_, _ = w.Write([]byte("{\"success\":\"ok\"}"))
	w.Header().Set("Content-Type", "application/json")
}

func (handler *PostReactionHandler) ReportPost(w http.ResponseWriter, r *http.Request) {
	var reportDTO dto.ReportDTO
	err := json.NewDecoder(r.Body).Decode(&reportDTO)
	if err != nil {
		fmt.Println(err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	err = handler.PostReactionService.ReportPost(reportDTO.PostID, reportDTO.Reason)
	if err != nil {
		fmt.Println(err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	w.WriteHeader(http.StatusOK)
	_, _ = w.Write([]byte("{\"success\":\"ok\"}"))
	w.Header().Set("Content-Type", "application/json")
}

func (handler *PostReactionHandler) GetMyReactions(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	reactionType := model.GetReactionType(strings.ToLower(vars["type"]))
	if reactionType == model.NONE {
		fmt.Println("Bad reaction type in request!")
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	loggedUserID := util.GetLoggedUserIDFromToken(r)
	posts, err := handler.PostReactionService.GetMyReactions(reactionType, loggedUserID)
	if err != nil {
		fmt.Println(err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	js, err := json.Marshal(posts)
	if err != nil {
		fmt.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
	_, _ = w.Write(js)
}

func (handler *PostReactionHandler) GetReactionTypes(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	profileID := util.String2Uint(params["profileID"])
	type data struct {
		Ids []string `json:"ids"`
	}
	var input data
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		fmt.Println(err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	ret := handler.PostReactionService.GetReactionTypes(profileID, input.Ids)
	js, err := json.Marshal(ret)
	if err != nil {
		fmt.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
	_, _ = w.Write(js)
}

func (handler *PostReactionHandler) GetAllReports(w http.ResponseWriter, r *http.Request) {
	reports, err := handler.PostReactionService.GetAllReports()
	if err != nil {
		fmt.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	js, err := json.Marshal(reports)
	if err != nil {
		fmt.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
	_, _ = w.Write(js)
	w.Header().Set("Content-Type", "application/json")
}
