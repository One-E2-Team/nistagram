package repository

import (
	"context"
	"go.mongodb.org/mongo-driver/mongo"
	"nistagram/postReaction/model"
)

type PostReactionRepository struct {
	Client *mongo.Client
}

func (repo *PostReactionRepository) ReactOnPost(reaction *model.Reaction) error {
	collection := repo.getCollection("reactions")
	_, err := collection.InsertOne(context.TODO(), reaction)
	return err
}

func (repo *PostReactionRepository) ReportPost(report *model.Report) error {
	collection := repo.getCollection("reports")
	_, err := collection.InsertOne(context.TODO(), report)
	return err
}

func (repo *PostReactionRepository) getCollection(name string) *mongo.Collection {
	return repo.Client.Database("postReactionDB").Collection(name)
}
