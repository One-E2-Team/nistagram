package model

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
	"time"
)

type Message struct{
	ID    primitive.ObjectID 	 `bson:"_id" json:"id,omitempty"`
	SenderId		uint		`json:"senderId"`
	ReceiverId		uint		`json:"receiverId"`
	Text			string		`json:"text"`
	Timestamp		time.Time	`json:"timestamp"`
	MediaPath		string		`json:"mediaPath"`
	OneOf			bool		`json:"oneOf"`
	Seen			bool		`json:"seen"`
	PostId			string		`json:"postId"`
}