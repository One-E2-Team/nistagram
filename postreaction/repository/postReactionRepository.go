package repository

import (
	"context"
	"fmt"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"nistagram/postreaction/model"
)

const reactionsCollectionName = "reactions"
const reportsCollectionName = "reports"

type PostReactionRepository struct {
	Client *mongo.Client
}

func (repo *PostReactionRepository) ReactOnPost(reaction *model.Reaction) error {
	reactionsCollection := repo.getCollection(reactionsCollectionName)
	filter := bson.D{{"profileid", reaction.ProfileID}, {"postid", reaction.PostID}}
	var existingReaction model.Reaction
	exists := reactionsCollection.FindOne(context.TODO(), filter).Decode(&existingReaction)
	if exists != nil {
		_, err := reactionsCollection.InsertOne(context.TODO(), reaction)
		return err
	}
	update := bson.D{
		{"$set", bson.D{
			{"reactiontype", reaction.ReactionType},
		}},
	}
	result, _ := reactionsCollection.UpdateOne(context.TODO(), filter, update)
	if result.MatchedCount == 0 {
		return mongo.ErrNoDocuments
	}
	return nil
}

func (repo *PostReactionRepository) ReportPost(report *model.Report) error {
	reportsCollection := repo.getCollection(reportsCollectionName)
	_, err := reportsCollection.InsertOne(context.TODO(), report)
	return err
}

func (repo *PostReactionRepository) GetProfileReactions(reactionType model.ReactionType, profileID uint) ([]model.Reaction, error) {
	reactionsCollection := repo.getCollection(reactionsCollectionName)
	filter := bson.D{{"profileid", profileID}, {"reactiontype", reactionType}}
	reactionCursor, err := reactionsCollection.Find(context.TODO(), filter)
	if err != nil {
		return nil, err
	}
	var reactions []model.Reaction
	for reactionCursor.Next(context.TODO()) {
		var result model.Reaction
		err = reactionCursor.Decode(&result)
		if err != nil {
			return nil, err
		}
		fmt.Println(result.ProfileID, result.ReactionType, result.PostID)
		reactions = append(reactions, result)
	}
	return reactions, nil
}

func (repo *PostReactionRepository) getCollection(name string) *mongo.Collection {
	return repo.Client.Database("postReactionDB").Collection(name)
}
