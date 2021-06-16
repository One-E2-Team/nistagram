package repository

import (
	"context"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"nistagram/postReaction/model"
)

type PostReactionRepository struct {
	Client *mongo.Client
}

func (repo *PostReactionRepository) ReactOnPost(reaction *model.Reaction) error {
	collection := repo.getCollection("reactions")
	filter := bson.D{{"profileid", reaction.ProfileID}, {"postid", reaction.PostID}}
	var existingReaction model.Reaction
	exists := collection.FindOne(context.TODO(), filter).Decode(&existingReaction)
	if exists != nil {
		_, err := collection.InsertOne(context.TODO(), reaction)
		return err
	}
	update := bson.D{
		{"$set", bson.D{
			{"reactiontype", reaction.ReactionType},
		}},
	}
	result, _ := collection.UpdateOne(context.TODO(), filter, update)
	if result.MatchedCount == 0 {
		return mongo.ErrNoDocuments
	}
	return nil
}

func (repo *PostReactionRepository) ReportPost(report *model.Report) error {
	collection := repo.getCollection("reports")
	_, err := collection.InsertOne(context.TODO(), report)
	return err
}

func (repo *PostReactionRepository) getCollection(name string) *mongo.Collection {
	return repo.Client.Database("postReactionDB").Collection(name)
}
