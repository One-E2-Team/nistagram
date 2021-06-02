package repository

import (
	"context"
	"errors"
	"go.mongodb.org/mongo-driver/mongo"
	"nistagram/post/model"
)

type PostRepository struct {
	Client *mongo.Client
}

func (repo *PostRepository) CreatePost(post *model.Post) error {
	collection:= repo.getCollection(post.PostType)
	if collection == nil {
		return errors.New("collection doesn't exist ")
	}
	_, err := collection.InsertOne(context.TODO(), post)
	return err
}

func (repo *PostRepository) getCollection(postType model.PostType) *mongo.Collection{
	if postType == model.PostType(0) {
		return repo.Client.Database("postdb").Collection("posts")
	}else if postType == model.PostType(1){
		return repo.Client.Database("postdb").Collection("stories")
	}
	return nil
}





