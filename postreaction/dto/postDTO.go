package dto

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type PostDTO struct {
	ID    primitive.ObjectID 	 `bson:"_id" json:"id,omitempty"`
	PublisherId        uint      `json:"publisherId"`
	PublisherUsername  string    `json:"publisherUsername"`
	Medias             []Media   `json:"medias"`
	Description        string    `json:"description"`
}

