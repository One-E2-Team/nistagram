package model

import "strings"

type PostType int

const (
	NONE PostType = iota
	STORY
	POST
)

func GetPostType(postType string) PostType {
	if strings.ToLower(postType) == "post" {
		return POST
	}
	if strings.ToLower(postType) == "story" {
		return STORY
	}
	return NONE
}
