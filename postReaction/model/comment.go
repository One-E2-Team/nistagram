package model

import "time"

type Comment struct {
	PostID    uint      `json:"postId"`
	ProfileID uint      `json:"profileId"`
	Content   string    `json:"content"`
	Time      time.Time `json:"time"`
}
