package model

type Reaction struct {
	ReactionType ReactionType `json:"reactionType"`
	PostID       uint         `json:"postId"`
	ProfileID    uint         `json:"profileId"`
}
