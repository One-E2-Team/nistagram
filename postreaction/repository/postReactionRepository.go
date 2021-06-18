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

const profileIDColumn = "profileid"
const postIDColumn = "postid"
const reactionTypeColumn = "reactiontype"

type PostReactionRepository struct {
	Client *mongo.Client
}

func (repo *PostReactionRepository) ReactOnPost(reaction *model.Reaction) error {
	reactionsCollection := repo.getCollection(reactionsCollectionName)
	filter := bson.D{{profileIDColumn, reaction.ProfileID}, {postIDColumn, reaction.PostID}}
	var existingReaction model.Reaction
	exists := reactionsCollection.FindOne(context.TODO(), filter).Decode(&existingReaction)
	if exists != nil {
		_, err := reactionsCollection.InsertOne(context.TODO(), reaction)
		return err
	}
	update := bson.D{
		{"$set", bson.D{
			{reactionTypeColumn, reaction.ReactionType},
		}},
	}
	result, _ := reactionsCollection.UpdateOne(context.TODO(), filter, update)
	if result.MatchedCount == 0 {
		return mongo.ErrNoDocuments
	}
	return nil
}

func (repo *PostReactionRepository) DeleteReaction(postID string, loggedUserID uint) error {
	reactionsCollection := repo.getCollection(reactionsCollectionName)
	filter := bson.D{{profileIDColumn, loggedUserID}, {postIDColumn, postID}}
	var existingReaction model.Reaction
	err := reactionsCollection.FindOne(context.TODO(), filter).Decode(&existingReaction)
	if err != nil {
		return err
	}
	_, err = reactionsCollection.DeleteOne(context.TODO(), existingReaction)
	return err
}

func (repo *PostReactionRepository) ReportPost(report *model.Report) error {
	reportsCollection := repo.getCollection(reportsCollectionName)
	_, err := reportsCollection.InsertOne(context.TODO(), report)
	return err
}

func (repo *PostReactionRepository) GetProfileReactions(reactionType model.ReactionType, profileID uint) ([]model.Reaction, error) {
	reactionsCollection := repo.getCollection(reactionsCollectionName)
	filter := bson.D{{profileIDColumn, profileID}, {reactionTypeColumn, reactionType}}
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

func (repo *PostReactionRepository) GetReactionType(profileID uint, postID string) string {
	reactionsCollection := repo.getCollection(reactionsCollectionName)
	filter := bson.D{{profileIDColumn, profileID}, {postIDColumn, postID}}
	var existingReaction model.Reaction
	err := reactionsCollection.FindOne(context.TODO(), filter).Decode(&existingReaction)
	if err != nil {
		return "none"
	}
	return model.GetReactionTypeString(existingReaction.ReactionType)
}

func (repo *PostReactionRepository) getCollection(name string) *mongo.Collection {
	return repo.Client.Database("postReactionDB").Collection(name)
}
