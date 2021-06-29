package model

import "go.mongodb.org/mongo-driver/bson/primitive"

type Media struct {
	ID    primitive.ObjectID 	 `bson:"_id" json:"id,omitempty"`
	FilePath string `json:"filePath"`
	WebSite  string `json:"webSite"`
}
