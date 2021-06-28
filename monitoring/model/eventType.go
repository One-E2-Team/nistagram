package model

const (
	LIKE EventType = iota
	DISLIKE
	COMMENT
	VISIT
	LIKE_RESET
	DISLIKE_RESET
)

type EventType int
