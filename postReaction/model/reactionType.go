package model

import "strings"

type ReactionType int

const (
	LIKE ReactionType = iota
	DISLIKE
	NONE
)

func GetReactionType(reactionType string) ReactionType {
	if strings.ToLower(reactionType) == "like" {
		return LIKE
	}
	if strings.ToLower(reactionType) == "dislike" {
		return DISLIKE
	}
	return NONE
}
