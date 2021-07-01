package dto

import (
	"time"
)

type Media struct {
	FilePath string `json:"filePath"`
	WebSite  string `json:"webSite"`
}

type PostDTO struct {
	PostType           int  `json:"postType"`
	Medias             []Media   `json:"medias"`
	PublishDate        time.Time `json:"publishDate"`
	Description        string    `json:"description"`
	IsHighlighted      bool      `json:"isHighlighted"`
	IsCloseFriendsOnly bool      `json:"isCloseFriendsOnly"`
	Location           string    `json:"location"`
	HashTags           string    `json:"hashTags"`
	TaggedUsers        []string  `json:"taggedUsers"`
	IsPrivate          bool      `json:"isPrivate"`
	IsDeleted          bool      `json:"isDeleted"`
}
