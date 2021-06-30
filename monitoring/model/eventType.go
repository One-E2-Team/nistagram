package model

import "strings"

const (
	LIKE EventType = iota
	DISLIKE
	COMMENT
	VISIT
	LIKE_RESET
	DISLIKE_RESET
	NONE
)

type EventType int

func GetEventType(eventType string) EventType {
	if strings.ToLower(eventType) == "like" {
		return LIKE
	}
	if strings.ToLower(eventType) == "dislike" {
		return DISLIKE
	}
	if strings.ToLower(eventType) == "comment" {
		return COMMENT
	}
	if strings.ToLower(eventType) == "visit" {
		return VISIT
	}
	if strings.ToLower(eventType) == "like_reset" {
		return LIKE_RESET
	}
	if strings.ToLower(eventType) == "dislike_reset" {
		return DISLIKE_RESET
	}
	return NONE
}

func (e EventType) ToString() string {
	switch e {
	case LIKE:
		return "LIKE"
	case DISLIKE:
		return "DISLIKE"
	case COMMENT:
		return "COMMENT"
	case VISIT:
		return "VISIT"
	case LIKE_RESET:
		return "LIKE_RESET"
	case DISLIKE_RESET:
		return "DISLIKE_RESET"
	default:
		return "NONE"
	}
}
