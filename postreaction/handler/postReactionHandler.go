package handler

import (
	"encoding/json"
	"fmt"
	"net/http"
	"nistagram/postreaction/dto"
	"nistagram/postreaction/model"
	"nistagram/postreaction/service"
	"nistagram/util"
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
	loggedUserID := util.GetLoggedUserIDFromToken(r)
	if loggedUserID == 0 {
		fmt.Println("NO TOKEN")
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	err = handler.PostReactionService.ReactOnPost(reactionDTO.PostID, loggedUserID, reactionType)
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
