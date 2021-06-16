package model

type Reaction struct {
	ReactionType ReactionType `json:"reactionType"`
	PostID       string       `json:"postId"`
	ProfileID    uint         `json:"profileId"`
}
