package repository

import (
	"context"
	"errors"
	"go.mongodb.org/mongo-driver/mongo"
	"nistagram/postReaction/model"
)

type PostReactionRepository struct {
	Client *mongo.Client
}

func (repo *PostReactionRepository) ReactOnPost(reaction *model.Reaction) error {
	collection, err := repo.getCollection(reaction.ReactionType)
	if err != nil {
		return err
	}
	_, err = collection.InsertOne(context.TODO(), reaction)
	return err
}

func (repo *PostReactionRepository) getCollection(reactionType model.ReactionType) (*mongo.Collection, error) {
	if reactionType == model.LIKE || reactionType == model.DISLIKE {
		return repo.Client.Database("postReactionDB").Collection("reactions"), nil
	}
	return nil, errors.New("COLLECTION_DOES_NOT_EXIST")
}
