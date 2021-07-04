package repository

import (
	"context"
	"fmt"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"nistagram/post/dto"
	"nistagram/post/model"
	"nistagram/util"
	"strings"
	"time"
)

const postsCollectionName = "posts"
const postDbName = "postdb"

type PostRepository struct {
	Client *mongo.Client
}

func (repo *PostRepository) GetProfilesPosts(ctx context.Context, followingProfiles []util.FollowingProfileDTO, targetUsername string) ([]model.Post, error) {
	span := util.Tracer.StartSpanFromContext(ctx, "GetProfilesPosts-repository")
	defer util.Tracer.FinishSpan(span)

	collection := repo.getCollection()

	filter := bson.D{{"isdeleted", false}, {"publisherusername", targetUsername}}

	cursor, err := collection.Find(context.TODO(), filter)
	if err != nil {
		util.Tracer.LogError(span, err)
		return nil, err
	}

	var posts []model.Post

	for cursor.Next(context.TODO()) {
		var result model.Post
		err = cursor.Decode(&result)
		if result.IsPrivate == false || util.IsFollowed(followingProfiles, result.PublisherId) {
			if result.PostType == model.GetPostType("story") {
				if result.IsCloseFriendsOnly && !util.IsCloseFriend(followingProfiles, result.PublisherId){
					continue
				}
				duration, err := time.ParseDuration("24h")
				if err != nil {
					util.Tracer.LogError(span, err)
					return nil, err
				}
				if result.IsHighlighted || time.Now().Before(result.PublishDate.Add(duration)) {
					posts = append(posts, result)
				}
			} else {
				posts = append(posts, result)
			}
		}
	}

	return posts, nil
}

func (repo *PostRepository) GetPublic(ctx context.Context, blockedRelationships []uint) ([]model.Post, error) {
	span := util.Tracer.StartSpanFromContext(ctx, "GetPublic-repository")
	defer util.Tracer.FinishSpan(span)

	collection := repo.getCollection()

	filter := bson.D{{"isdeleted", false}, {"isprivate", false}}

	cursor, err := collection.Find(context.TODO(), filter)
	if err != nil {
		util.Tracer.LogError(span, err)
		return nil, err
	}

	var posts []model.Post

	for cursor.Next(context.TODO()) {
		var result model.Post
		err = cursor.Decode(&result)
		if !util.Contains(blockedRelationships, result.PublisherId) {
			if result.PostType == model.GetPostType("story") {
				duration, err := time.ParseDuration("24h")
				if err != nil {
					util.Tracer.LogError(span, err)
					return nil, err
				}
				if time.Now().Before(result.PublishDate.Add(duration)) {
					posts = append(posts, result)
				}
			} else {
				posts = append(posts, result)
			}
		}
	}

	return posts, nil
}

func (repo *PostRepository) GetPublicPostByLocation(ctx context.Context, location string, blockedRelationships []uint) ([]model.Post, error) {
	span := util.Tracer.StartSpanFromContext(ctx, "GetPublicPostByLocation-repository")
	defer util.Tracer.FinishSpan(span)
	collection := repo.getCollection()

	filter := bson.D{{"isdeleted", false}, {"isprivate", false}}

	cursor, err := collection.Find(context.TODO(), filter)
	if err != nil {
		util.Tracer.LogError(span, err)
		return nil, err
	}

	var posts []model.Post

	for cursor.Next(context.TODO()) {
		var result model.Post
		err = cursor.Decode(&result)
		if !util.Contains(blockedRelationships, result.PublisherId) {
			if strings.Contains(result.Location, location) {
				if result.PostType == model.GetPostType("story") {
					duration, err := time.ParseDuration("24h")
					if err != nil {
						util.Tracer.LogError(span, err)
						return nil, err
					}
					if time.Now().Before(result.PublishDate.Add(duration)) {
						posts = append(posts, result)
					}
				} else {
					posts = append(posts, result)
				}
			}
		}
	}

	return posts, nil
}

func (repo *PostRepository) GetPublicPostByHashTag(ctx context.Context, hashTag string, blockedRelationships []uint) ([]model.Post, error) {
	span := util.Tracer.StartSpanFromContext(ctx, "GetPublicPostByHashTag-repository")
	defer util.Tracer.FinishSpan(span)
	collection := repo.getCollection()

	filter := bson.D{{"isdeleted", false}, {"isprivate", false}}

	cursor, err := collection.Find(context.TODO(), filter)
	if err != nil {
		util.Tracer.LogError(span, err)
		return nil, err
	}

	var posts []model.Post

	for cursor.Next(context.TODO()) {
		var result model.Post
		err = cursor.Decode(&result)
		if !util.Contains(blockedRelationships, result.PublisherId) {
			if strings.Contains(result.HashTags, hashTag) {
				if result.PostType == model.GetPostType("story") {
					duration, err := time.ParseDuration("24h")
					if err != nil {
						util.Tracer.LogError(span, err)
						return nil, err
					}
					if time.Now().Before(result.PublishDate.Add(duration)) {
						posts = append(posts, result)
					}
				} else {
					posts = append(posts, result)
				}
			}
		}
	}

	return posts, nil
}

func (repo *PostRepository) GetMyPosts(ctx context.Context, loggedUserId uint) ([]model.Post, error) {
	span := util.Tracer.StartSpanFromContext(ctx, "GetMyPosts-repository")
	defer util.Tracer.FinishSpan(span)

	collection := repo.getCollection()

	filter := bson.D{{"isdeleted", false}, {"publisherid", loggedUserId}}

	cursor, err := collection.Find(context.TODO(), filter)
	if err != nil {
		util.Tracer.LogError(span, err)
		return nil, err
	}

	var posts []model.Post

	for cursor.Next(context.TODO()) {
		var result model.Post
		err = cursor.Decode(&result)
		if err != nil{
			util.Tracer.LogError(span, err)
		}
		posts = append(posts, result)
	}

	return posts, nil
}

func (repo *PostRepository) GetPostsForHomePage(ctx context.Context, followingProfiles []util.FollowingProfileDTO) ([]model.Post, error) {
	span := util.Tracer.StartSpanFromContext(ctx, "GetPostsForHomePage-repository")
	defer util.Tracer.FinishSpan(span)

	collection := repo.getCollection()

	cursor, err := collection.Find(context.TODO(), bson.D{})
	if err != nil {
		util.Tracer.LogError(span, err)
		return nil, err
	}

	var posts []model.Post

	for cursor.Next(context.TODO()) {
		var result model.Post
		err = cursor.Decode(&result)
		if util.IsFollowed(followingProfiles, result.PublisherId) && !result.IsDeleted {
			if result.PostType == model.GetPostType("story") {
				if result.IsCloseFriendsOnly && !util.IsCloseFriend(followingProfiles, result.PublisherId) {
					continue
				}

				duration, err := time.ParseDuration("24h")
				if err != nil {
					util.Tracer.LogError(span, err)
					return nil, err
				}
				if time.Now().Before(result.PublishDate.Add(duration)) {
					posts = append(posts, result)
				}
			} else {
				posts = append(posts, result)
			}
		}
	}

	return posts, nil
}

func (repo *PostRepository) Create(ctx context.Context, post *model.Post) error {
	span := util.Tracer.StartSpanFromContext(ctx, "Create-repository")
	defer util.Tracer.FinishSpan(span)
	collection := repo.getCollection()
	_, err := collection.InsertOne(context.TODO(), post)
	return err
}

func (repo *PostRepository) Read(ctx context.Context, id primitive.ObjectID) (model.Post, error) {
	span := util.Tracer.StartSpanFromContext(ctx, "Read-repository")
	defer util.Tracer.FinishSpan(span)
	collection := repo.getCollection()
	filter := bson.D{{"_id", id}, {"isdeleted", false}}
	var result model.Post
	err := collection.FindOne(context.TODO(), filter).Decode(&result)
	return result, err
}

func (repo *PostRepository) Delete(ctx context.Context, id primitive.ObjectID) error {
	span := util.Tracer.StartSpanFromContext(ctx, "Delete-repository")
	defer util.Tracer.FinishSpan(span)
	collection := repo.getCollection()
	filter := bson.D{{"_id", id}}
	update := bson.D{
		{"$set", bson.D{
			{"isdeleted", true},
		}},
	}
	result, _ := collection.UpdateOne(context.TODO(), filter, update)
	if result.MatchedCount == 0 {
		util.Tracer.LogError(span, fmt.Errorf("post not found"))
		return mongo.ErrNoDocuments
	}
	return nil
}

func (repo *PostRepository) Update(ctx context.Context, id primitive.ObjectID, post dto.PostDto) error {
	span := util.Tracer.StartSpanFromContext(ctx, "Update-repository")
	defer util.Tracer.FinishSpan(span)
	collection := repo.getCollection()
	filter := bson.D{{"_id", id}}
	update := bson.D{
		{"$set", bson.D{
			{"ishighlighted", post.IsHighlighted},
			{"isclosefriendsonly", post.IsHighlighted},
		}},
	}

	result, _ := collection.UpdateOne(context.TODO(), filter, update)
	if result.MatchedCount == 0 {
		util.Tracer.LogError(span, fmt.Errorf("post not found"))
		return mongo.ErrNoDocuments
	}
	return nil
}

func (repo *PostRepository) DeleteUserPosts(ctx context.Context, id uint) error {
	span := util.Tracer.StartSpanFromContext(ctx, "DeleteUserPosts-repository")
	defer util.Tracer.FinishSpan(span)
	filter := bson.D{{"publisherid", id}}
	update := bson.D{
		{"$set", bson.D{
			{"isdeleted", true},
		}},
	}

	return repo.updateMany(ctx, filter, update)
}

func (repo *PostRepository) ChangeUsername(ctx context.Context, id uint, username string) error {
	span := util.Tracer.StartSpanFromContext(ctx, "ChangeUsername-repository")
	defer util.Tracer.FinishSpan(span)
	filter := bson.D{{"publisherid", id}}
	update := bson.D{
		{"$set", bson.D{
			{"publisherusername", username},
		}},
	}

	return repo.updateMany(ctx, filter, update)
}

func (repo *PostRepository) ChangePrivacy(ctx context.Context, id uint, isPrivate bool) error {
	span := util.Tracer.StartSpanFromContext(ctx, "ChangePrivacy-repository")
	defer util.Tracer.FinishSpan(span)
	filter := bson.D{{"publisherid", id}}
	update := bson.D{
		{"$set", bson.D{
			{"isprivate", isPrivate},
		}},
	}

	return repo.updateMany(ctx, filter, update)
}

func (repo *PostRepository) MakeCampaign(ctx context.Context, id primitive.ObjectID) error {
	span := util.Tracer.StartSpanFromContext(ctx, "MakeCampaign-repository")
	defer util.Tracer.FinishSpan(span)
	collection := repo.getCollection()
	filter := bson.D{{"_id", id}}
	update := bson.D{
		{"$set", bson.D{
			{"iscampaign", true},
		}},
	}

	result, _ := collection.UpdateOne(context.TODO(), filter, update)
	if result.MatchedCount == 0 {
		util.Tracer.LogError(span, fmt.Errorf("post not found"))
		return mongo.ErrNoDocuments
	}
	return nil
}

func (repo *PostRepository) GetMediaById(ctx context.Context, id string) (model.Media, error) {
	span := util.Tracer.StartSpanFromContext(ctx, "GetMediaById-repository")
	defer util.Tracer.FinishSpan(span)
	var retMedia model.Media
	mediaId, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		util.Tracer.LogError(span, err)
		return retMedia, err
	}
	collection := repo.getCollection()

	filter := bson.D{{"medias._id", mediaId}}

	cursor, err := collection.Find(context.TODO(), filter)
	if err != nil {
		util.Tracer.LogError(span, err)
		return retMedia, err
	}

	var post model.Post

	for cursor.Next(context.TODO()) {
		err = cursor.Decode(&post)
	}


	for _, media := range post.Medias{
		if media.ID == mediaId{
			retMedia = media
		}
	}

	return retMedia, nil
}

func (repo *PostRepository) updateMany(ctx context.Context, filter bson.D, update bson.D) error {
	span := util.Tracer.StartSpanFromContext(ctx, "updateMany-repository")
	defer util.Tracer.FinishSpan(span)
	collection := repo.getCollection()

	result, _ := collection.UpdateMany(context.TODO(), filter, update)

	if result.MatchedCount == 0 {
		util.Tracer.LogError(span, fmt.Errorf("post not found"))
		return mongo.ErrNoDocuments
	}
	return nil
}

func (repo *PostRepository) getCollection() *mongo.Collection {
	return repo.Client.Database(postDbName).Collection(postsCollectionName)
}
