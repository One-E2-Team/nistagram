package model

import "time"

type Message struct{
	SenderId		uint		`json:"senderId"`
	ReceiverId		uint		`json:"receiverId"`
	Text			string		`json:"text"`
	Timestamp		time.Time	`json:"timestamp"`
	MediaPath		string		`json:"mediaPath"`
	OneOf			bool		`json:"oneOf"`
	Seen			bool		`json:"seen"`
	PostId			string		`json:"postId"`
}