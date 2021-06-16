package service

import (
	"nistagram/postReaction/model"
	"nistagram/postReaction/repository"
)

type PostReactionService struct {
	PostReactionRepository *repository.PostReactionRepository
}

func (service *PostReactionService) ReactOnPost(postID string, loggedUserID uint, reactionType model.ReactionType) error {
	reaction := model.Reaction{ReactionType: reactionType, PostID: postID, ProfileID: loggedUserID}
	err := service.PostReactionRepository.ReactOnPost(&reaction)
	return err
}
