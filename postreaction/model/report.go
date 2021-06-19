package model

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
	"time"
)

type Report struct {
	ID       primitive.ObjectID  `bson:"_id" json:"id,omitempty"`
	PostID   string              `json:"postId"`
	Time     time.Time           `json:"time"`
	Reason   string              `json:"reason"`
}
