package model

import "strings"

type ReactionType int

const (
	LIKE ReactionType = iota
	DISLIKE
	LIKE_RESET
	DISLIKE_RESET
	NONE
)

func GetReactionTypeString(reactionType ReactionType) string {
	switch reactionType {
	case LIKE:
		return "like"
	case DISLIKE:
		return "dislike"
	case LIKE_RESET:
		return "like_reset"
	case DISLIKE_RESET:
		return "dislike_reset"
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
	if strings.ToLower(reactionType) == "like_reset" {
		return LIKE_RESET
	}
	if strings.ToLower(reactionType) == "dislike_reset" {
		return DISLIKE_RESET
	}
	return NONE
}
