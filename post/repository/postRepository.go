package repository

import (
	"context"
	"errors"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"nistagram/post/dto"
	"nistagram/post/model"
)

type PostRepository struct {
	Client *mongo.Client
}

func (repo *PostRepository) Create(post *model.Post) error {
	collection,err := repo.getCollection(post.PostType)
	if err != nil { return err 	}
	_, err = collection.InsertOne(context.TODO(), post)
	return err
}

func (repo *PostRepository) Read(id primitive.ObjectID, postType model.PostType ) (model.Post, error) {
	collection, err:= repo.getCollection(postType)
	if err != nil { return model.Post{},err }
	filter := bson.D{{"_id", id}}
	var result model.Post
	err = collection.FindOne(context.TODO(), filter).Decode(&result)
	return result,err
}

func (repo *PostRepository) Delete(id primitive.ObjectID, postType model.PostType)  error {
	collection, err:= repo.getCollection(postType)
	if err != nil { return err }
	filter := bson.D{{"_id", id}}
	update := bson.D{
		{"$set", bson.D{
			{"isDeleted", true},
		}},
	}
	result, _ := collection.UpdateOne(context.TODO(), filter,update)
	if result.MatchedCount == 0 {return mongo.ErrNoDocuments}
	return nil
}

func (repo *PostRepository) Update(id primitive.ObjectID,postType model.PostType, post dto.PostDto) error {

	collection, err:= repo.getCollection(postType)
	if err != nil { return err }

	filter := bson.D{{"_id", id}}
	update := bson.D{
		{"$set", bson.D{
			{"isHighlighted",post.IsHighlighted},
			{"isCloseFriendsOnly", post.IsHighlighted},
		}},
	}

	result, _ := collection.UpdateOne(context.TODO(), filter,update)
	if result.MatchedCount == 0 {return mongo.ErrNoDocuments}
	return nil
}

func (repo *PostRepository) getCollection(postType model.PostType) (*mongo.Collection,error) {
	if postType == model.POST {
		return repo.Client.Database("postdb").Collection("posts") , nil
	}else if postType == model.STORY{
		return repo.Client.Database("postdb").Collection("stories"), nil
	}
	return nil, errors.New("collection doesn't exist!")
}









