package model

import "time"

type Report struct {
	PostID string    `json:"postId"`
	Time   time.Time `json:"time"`
	Reason string    `json:"reason"`
}
