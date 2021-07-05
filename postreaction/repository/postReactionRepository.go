package repository

import (
	"context"
	"fmt"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"nistagram/postreaction/model"
	"nistagram/util"
)

const reactionDbName = "postReactionDB"
const reactionsCollectionName = "reactions"
const reportsCollectionName = "reports"
const commentsCollectionName = "comments"

const profileIDColumn = "profileid"
const postIDColumn = "postid"
const reactionTypeColumn = "reactiontype"

var emptyContext = context.TODO()

type PostReactionRepository struct {
	Client *mongo.Client
}

func (repo *PostReactionRepository) ReactOnPost(ctx context.Context, reaction *model.Reaction) (*model.ReactionType, error) {
	span := util.Tracer.StartSpanFromContext(ctx, "ReactOnPost-repository")
	defer util.Tracer.FinishSpan(span)

	reactionsCollection := repo.getCollection(reactionsCollectionName)
	filter := bson.D{{profileIDColumn, reaction.ProfileID}, {postIDColumn, reaction.PostID}}
	var existingReaction model.Reaction
	exists := reactionsCollection.FindOne(emptyContext, filter).Decode(&existingReaction)
	if exists != nil {
		_, err := reactionsCollection.InsertOne(emptyContext, reaction)
		if err != nil {
			util.Tracer.LogError(span, err)
			return nil, err
		}
		return nil, nil
	}
	update := bson.D{
		{"$set", bson.D{
			{reactionTypeColumn, reaction.ReactionType},
		}},
	}
	result, _ := reactionsCollection.UpdateOne(emptyContext, filter, update)
	if result.MatchedCount == 0 {
		util.Tracer.LogError(span, fmt.Errorf("reaction not updated"))
		return nil, mongo.ErrNoDocuments
	}
	return &existingReaction.ReactionType, nil
}

func (repo *PostReactionRepository) DeleteReaction(ctx context.Context, postID string, loggedUserID uint) (*model.Reaction, error) {
	span := util.Tracer.StartSpanFromContext(ctx, "DeleteReaction-repository")
	defer util.Tracer.FinishSpan(span)

	reactionsCollection := repo.getCollection(reactionsCollectionName)
	filter := bson.D{{profileIDColumn, loggedUserID}, {postIDColumn, postID}}
	var existingReaction model.Reaction
	err := reactionsCollection.FindOne(emptyContext, filter).Decode(&existingReaction)
	if err != nil {
		util.Tracer.LogError(span, err)
		return nil, err
	}
	_, err = reactionsCollection.DeleteOne(emptyContext, existingReaction)
	if err != nil {
		util.Tracer.LogError(span, err)
		return nil, err
	}
	return &existingReaction, err
}

func (repo *PostReactionRepository) CommentPost(ctx context.Context, comment *model.Comment) error {
	span := util.Tracer.StartSpanFromContext(ctx, "CommentPost-repository")
	defer util.Tracer.FinishSpan(span)

	commentsCollection := repo.getCollection(commentsCollectionName)
	_, err := commentsCollection.InsertOne(emptyContext, comment)
	if err != nil {
		util.Tracer.LogError(span, err)
	}
	return err
}

func (repo *PostReactionRepository) ReportPost(ctx context.Context, report *model.Report) error {
	span := util.Tracer.StartSpanFromContext(ctx, "ReportPost-repository")
	defer util.Tracer.FinishSpan(span)
	reportsCollection := repo.getCollection(reportsCollectionName)
	_, err := reportsCollection.InsertOne(emptyContext, report)
	if err != nil {
		util.Tracer.LogError(span, err)
	}
	return err
}

func (repo *PostReactionRepository) GetProfileReactions(ctx context.Context, reactionType model.ReactionType, profileID uint) ([]model.Reaction, error) {
	span := util.Tracer.StartSpanFromContext(ctx, "GetProfileReactions-repository")
	defer util.Tracer.FinishSpan(span)

	reactionsCollection := repo.getCollection(reactionsCollectionName)
	filter := bson.D{{profileIDColumn, profileID}, {reactionTypeColumn, reactionType}}
	reactionCursor, err := reactionsCollection.Find(emptyContext, filter)
	if err != nil {
		util.Tracer.LogError(span, err)
		return nil, err
	}
	var reactions []model.Reaction
	for reactionCursor.Next(emptyContext) {
		var result model.Reaction
		err = reactionCursor.Decode(&result)
		if err != nil {
			util.Tracer.LogError(span, err)
			return nil, err
		}
		fmt.Println(result.ProfileID, result.ReactionType, result.PostID)
		reactions = append(reactions, result)
	}
	return reactions, nil
}

func (repo *PostReactionRepository) GetReactionType(ctx context.Context, profileID uint, postID string) string {
	span := util.Tracer.StartSpanFromContext(ctx, "GetReactionType-repository")
	defer util.Tracer.FinishSpan(span)

	reactionsCollection := repo.getCollection(reactionsCollectionName)
	filter := bson.D{{profileIDColumn, profileID}, {postIDColumn, postID}}
	var existingReaction model.Reaction
	err := reactionsCollection.FindOne(emptyContext, filter).Decode(&existingReaction)
	if err != nil {
		util.Tracer.LogError(span, err)
		return "none"
	}
	return model.GetReactionTypeString(existingReaction.ReactionType)
}

func (repo *PostReactionRepository) GetAllReports(ctx context.Context) ([]model.Report, error) {
	span := util.Tracer.StartSpanFromContext(ctx, "GetAllReports-repository")
	defer util.Tracer.FinishSpan(span)

	var reports []model.Report
	reportsCollection := repo.getCollection(reportsCollectionName)
	filter := bson.D{{"isdeleted", false}}
	cursor, err := reportsCollection.Find(context.TODO(), filter)
	if err != nil {
		util.Tracer.LogError(span, err)
		return nil, err
	}
	for cursor.Next(context.TODO()) {
		var result model.Report
		err = cursor.Decode(&result)
		if err != nil {
			util.Tracer.LogError(span, err)
			return nil, err
		}
		reports = append(reports, result)
	}
	return reports, nil
}

func (repo *PostReactionRepository) GetReportsByPostId(ctx context.Context, postId string) ([]model.Report, error) {
	span := util.Tracer.StartSpanFromContext(ctx, "GetReportsByPostId-repository")
	defer util.Tracer.FinishSpan(span)

	var reports []model.Report
	reportsCollection := repo.getCollection(reportsCollectionName)
	filter := bson.D{{"postid", postId}}
	cursor, err := reportsCollection.Find(context.TODO(), filter)
	if err != nil {
		util.Tracer.LogError(span, err)
		return nil, err
	}
	for cursor.Next(context.TODO()) {
		var result model.Report
		err = cursor.Decode(&result)
		if err != nil {
			util.Tracer.LogError(span, err)
			return nil, err
		}
		reports = append(reports, result)
	}
	return reports, nil
}

func (repo *PostReactionRepository) DeleteReport(ctx context.Context, reportId primitive.ObjectID) error {
	span := util.Tracer.StartSpanFromContext(ctx, "DeleteReport-repository")
	defer util.Tracer.FinishSpan(span)

	reportsCollection := repo.getCollection(reportsCollectionName)
	filter := bson.D{{"_id", reportId}}
	update := bson.D{
		{"$set", bson.D{
			{"isdeleted", true},
		}},
	}

	_, err := reportsCollection.UpdateOne(emptyContext, filter, update)
	if err != nil {
		util.Tracer.LogError(span, err)
		return err
	}

	return nil
}

func (repo *PostReactionRepository) GetAllReactions(ctx context.Context, postID string) ([]uint, []uint, error) {
	span := util.Tracer.StartSpanFromContext(ctx, "GetAllReactions-repository")
	defer util.Tracer.FinishSpan(span)

	reactionsCollection := repo.getCollection(reactionsCollectionName)
	filter := bson.D{{postIDColumn, postID}}
	reactionCursor, err := reactionsCollection.Find(emptyContext, filter)
	if err != nil {
		util.Tracer.LogError(span, err)
		return nil, nil, err
	}
	likes := make([]uint, 0)
	dislikes := make([]uint, 0)
	for reactionCursor.Next(emptyContext) {
		var result model.Reaction
		err = reactionCursor.Decode(&result)
		if err != nil {
			util.Tracer.LogError(span, err)
			return nil, nil, err
		}
		if result.ReactionType == model.LIKE {
			likes = append(likes, result.ProfileID)
		} else if result.ReactionType == model.DISLIKE {
			dislikes = append(dislikes, result.ProfileID)
		}
	}
	return likes, dislikes, nil
}

func (repo *PostReactionRepository) GetAllComments(ctx context.Context, postID string) ([]model.Comment, error) {
	span := util.Tracer.StartSpanFromContext(ctx, "GetAllComments-repository")
	defer util.Tracer.FinishSpan(span)

	commentsCollection := repo.getCollection(commentsCollectionName)
	filter := bson.D{{postIDColumn, postID}}
	commentsCursor, err := commentsCollection.Find(emptyContext, filter)
	if err != nil {
		util.Tracer.LogError(span, err)
		return nil, err
	}
	comments := make([]model.Comment, 0)
	for commentsCursor.Next(emptyContext) {
		var result model.Comment
		err = commentsCursor.Decode(&result)
		if err != nil {
			util.Tracer.LogError(span, err)
			return nil, err
		}
		comments = append(comments, result)
	}
	return comments, nil
}

func (repo *PostReactionRepository) getCollection(name string) *mongo.Collection {
	return repo.Client.Database(reactionDbName).Collection(name)
}
