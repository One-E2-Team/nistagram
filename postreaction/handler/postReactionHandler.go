package handler

import (
	"context"
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
	span := util.Tracer.StartSpanFromRequest("ReactOnPost-handler", r)
	defer util.Tracer.FinishSpan(span)

	var reactionDTO dto.ReactionDTO
	err := json.NewDecoder(r.Body).Decode(&reactionDTO)
	if err != nil {
		util.Tracer.LogError(span, err)
		fmt.Println(err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	reactionType := model.GetReactionType(reactionDTO.ReactionType)
	if reactionType == model.NONE {
		util.Tracer.LogError(span, fmt.Errorf("bad reaction type in request"))
		fmt.Println("Bad reaction type in request!")
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	ctx := util.Tracer.ContextWithSpan(context.Background(), span)
	err = handler.PostReactionService.ReactOnPost(ctx, reactionDTO, util.GetLoggedUserIDFromToken(r), reactionType)
	if err != nil {
		util.Tracer.LogError(span, err)
		fmt.Println(err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	w.WriteHeader(http.StatusOK)
	_, _ = w.Write([]byte("{\"success\":\"ok\"}"))
	w.Header().Set("Content-Type", "application/json")
}

func (handler *PostReactionHandler) DeleteReaction(w http.ResponseWriter, r *http.Request) {
	span := util.Tracer.StartSpanFromRequest("DeleteReaction-handler", r)
	defer util.Tracer.FinishSpan(span)

	var reactionDTO dto.ReactionDTO
	err := json.NewDecoder(r.Body).Decode(&reactionDTO)
	if err != nil {
		util.Tracer.LogError(span, err)
		fmt.Println(err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	ctx := util.Tracer.ContextWithSpan(context.Background(), span)
	err = handler.PostReactionService.DeleteReaction(ctx, reactionDTO, util.GetLoggedUserIDFromToken(r))
	if err != nil {
		util.Tracer.LogError(span, err)
		fmt.Println(err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	w.WriteHeader(http.StatusOK)
	_, _ = w.Write([]byte("{\"success\":\"ok\"}"))
	w.Header().Set("Content-Type", "application/json")
}

func (handler *PostReactionHandler) CommentPost(w http.ResponseWriter, r *http.Request) {
	span := util.Tracer.StartSpanFromRequest("CommentPost-handler", r)
	defer util.Tracer.FinishSpan(span)
	var commentDTO dto.CommentDTO
	err := json.NewDecoder(r.Body).Decode(&commentDTO)
	if err != nil {
		util.Tracer.LogError(span, err)
		fmt.Println(err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	ctx := util.Tracer.ContextWithSpan(context.Background(), span)
	err = handler.PostReactionService.CommentPost(ctx, commentDTO, util.GetLoggedUserIDFromToken(r))
	if err != nil {
		util.Tracer.LogError(span, err)
		fmt.Println(err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	w.WriteHeader(http.StatusOK)
	_, _ = w.Write([]byte("{\"success\":\"ok\"}"))
	w.Header().Set("Content-Type", "application/json")
}

func (handler *PostReactionHandler) ReportPost(w http.ResponseWriter, r *http.Request) {
	span := util.Tracer.StartSpanFromRequest("ReportPost-handler", r)
	defer util.Tracer.FinishSpan(span)
	var reportDTO dto.ReportDTO
	err := json.NewDecoder(r.Body).Decode(&reportDTO)
	if err != nil {
		util.Tracer.LogError(span, err)
		fmt.Println(err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	ctx := util.Tracer.ContextWithSpan(context.Background(), span)
	err = handler.PostReactionService.ReportPost(ctx, reportDTO.PostID, reportDTO.Reason)
	if err != nil {
		util.Tracer.LogError(span, err)
		fmt.Println(err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	w.WriteHeader(http.StatusOK)
	_, _ = w.Write([]byte("{\"success\":\"ok\"}"))
	w.Header().Set("Content-Type", "application/json")
}

func (handler *PostReactionHandler) GetMyReactions(w http.ResponseWriter, r *http.Request) {
	span := util.Tracer.StartSpanFromRequest("GetMyReactions-handler", r)
	defer util.Tracer.FinishSpan(span)
	vars := mux.Vars(r)
	reactionType := model.GetReactionType(strings.ToLower(vars["type"]))
	if reactionType == model.NONE {
		util.Tracer.LogError(span, fmt.Errorf("bad reaction type in request"))
		fmt.Println("Bad reaction type in request!")
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	loggedUserID := util.GetLoggedUserIDFromToken(r)
	ctx := util.Tracer.ContextWithSpan(context.Background(), span)
	posts, err := handler.PostReactionService.GetMyReactions(ctx, reactionType, loggedUserID)
	if err != nil {
		util.Tracer.LogError(span, err)
		fmt.Println(err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	js, err := json.Marshal(posts)
	if err != nil {
		util.Tracer.LogError(span, err)
		fmt.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
	_, _ = w.Write(js)
}

func (handler *PostReactionHandler) GetReactionTypes(w http.ResponseWriter, r *http.Request) {
	span := util.Tracer.StartSpanFromRequest("GetReactionTypes-handler", r)
	defer util.Tracer.FinishSpan(span)
	params := mux.Vars(r)
	profileID := util.String2Uint(params["profileID"])
	type data struct {
		Ids []string `json:"ids"`
	}
	var input data
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		util.Tracer.LogError(span, err)
		fmt.Println(err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	ctx := util.Tracer.ContextWithSpan(context.Background(), span)
	ret := handler.PostReactionService.GetReactionTypes(ctx, profileID, input.Ids)
	js, err := json.Marshal(ret)
	if err != nil {
		util.Tracer.LogError(span, err)
		fmt.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
	_, _ = w.Write(js)
}

func (handler *PostReactionHandler) GetAllReports(w http.ResponseWriter, r *http.Request) {
	span := util.Tracer.StartSpanFromRequest("GetAllReports-handler", r)
	defer util.Tracer.FinishSpan(span)

	ctx := util.Tracer.ContextWithSpan(context.Background(), span)
	reports, err := handler.PostReactionService.GetAllReports(ctx)
	if err != nil {
		util.Tracer.LogError(span, err)
		fmt.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		_, _ = w.Write([]byte("{\"message\":\"error\"}"))
		return
	}
	js, err := json.Marshal(reports)
	if err != nil {
		util.Tracer.LogError(span, err)
		fmt.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		_, _ = w.Write([]byte("{\"message\":\"error\"}"))
		return
	}
	w.WriteHeader(http.StatusOK)
	_, _ = w.Write(js)
	w.Header().Set("Content-Type", "application/json")
}

func (handler *PostReactionHandler) DeletePostsReports(w http.ResponseWriter, r *http.Request) {
	span := util.Tracer.StartSpanFromRequest("DeletePostsReports-handler", r)
	defer util.Tracer.FinishSpan(span)

	vars := mux.Vars(r)
	postId := vars["postId"]
	ctx := util.Tracer.ContextWithSpan(context.Background(), span)
	err := handler.PostReactionService.DeletePostsReports(ctx, postId)
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

func (handler *PostReactionHandler) GetAllReactions(w http.ResponseWriter, r *http.Request) {
	span := util.Tracer.StartSpanFromRequest("GetAllReactions-handler", r)
	defer util.Tracer.FinishSpan(span)
	vars := mux.Vars(r)
	ctx := util.Tracer.ContextWithSpan(context.Background(), span)
	likes, dislikes, err := handler.PostReactionService.GetAllReactions(ctx, vars["postID"])
	if err != nil {
		util.Tracer.LogError(span, err)
		fmt.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	type ReturnDTO struct {
		Likes    []string `json:"likes"`
		Dislikes []string `json:"dislikes"`
	}
	ret := ReturnDTO{Likes: likes, Dislikes: dislikes}
	js, err := json.Marshal(ret)
	if err != nil {
		util.Tracer.LogError(span, err)
		fmt.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		_, _ = w.Write([]byte("{\"message\":\"error\"}"))
		return
	}
	w.WriteHeader(http.StatusOK)
	_, _ = w.Write(js)
	w.Header().Set("Content-Type", "application/json")
}

func (handler *PostReactionHandler) GetAllComments(w http.ResponseWriter, r *http.Request) {
	span := util.Tracer.StartSpanFromRequest("GetAllComments-handler", r)
	defer util.Tracer.FinishSpan(span)

	vars := mux.Vars(r)
	ctx := util.Tracer.ContextWithSpan(context.Background(), span)
	comments, err := handler.PostReactionService.GetAllComments(ctx, vars["postID"])
	if err != nil {
		util.Tracer.LogError(span, err)
		fmt.Println(err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	js, err := json.Marshal(comments)
	if err != nil {
		util.Tracer.LogError(span, err)
		fmt.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
	_, _ = w.Write(js)
	w.Header().Set("Content-Type", "application/json")
}
