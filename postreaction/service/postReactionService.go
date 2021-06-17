package service

import (
	"nistagram/postreaction/model"
	"nistagram/postreaction/repository"
	"time"
)

type PostReactionService struct {
	PostReactionRepository *repository.PostReactionRepository
}

func (service *PostReactionService) ReactOnPost(postID string, loggedUserID uint, reactionType model.ReactionType) error {
	reaction := model.Reaction{ReactionType: reactionType, PostID: postID, ProfileID: loggedUserID}
	return service.PostReactionRepository.ReactOnPost(&reaction)
}

func (service *PostReactionService) ReportPost(postID string, reason string) error {
	report := model.Report{PostID: postID, Time: time.Now(), Reason: reason}
	return service.PostReactionRepository.ReportPost(&report)
}
