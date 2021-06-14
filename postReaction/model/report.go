package model

import "time"

type Report struct {
	PostID uint      `json:"postId"`
	Time   time.Time `json:"time"`
	Reason string    `json:"reason"`
}
