package model

import "strings"

type ReactionType int

const (
	LIKE ReactionType = iota
	DISLIKE
	NONE
)

func GetReactionTypeString(reactionType ReactionType) string {
	switch reactionType {
	case LIKE:
		return "like"
	case DISLIKE:
		return "dislike"
	}
	return "none"
}

func GetReactionType(reactionType string) ReactionType {
	if strings.ToLower(reactionType) == "like" {
		return LIKE
	}
	if strings.ToLower(reactionType) == "dislike" {
		return DISLIKE
	}
	return NONE
}
