package dto

type PostDto struct {
	Description        string   `json:"description"`
	IsHighlighted      bool     `json:"isHighlighted"`
	IsCloseFriendsOnly bool     `json:"isCloseFriendsOnly"`
	Location           string   `json:"location"`
	HashTags           string   `json:"hashTags"`
	PostType           string   `json:"postType"`
	Links              []string `json:"links"`
}
