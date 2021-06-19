package repository

import (
	"context"
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

func (repo *PostRepository) GetProfilesPosts(followingProfiles []uint, targetUsername string) ([]model.Post, error) {
	collection := repo.getCollection()

	filter := bson.D{{"isdeleted", false}, {"publisherusername", targetUsername}}

	cursor, err := collection.Find(context.TODO(), filter)
	if err != nil {
		return nil, err
	}

	var posts []model.Post

	for cursor.Next(context.TODO()) {
		var result model.Post
		err = cursor.Decode(&result)
		if result.IsPrivate == false || util.Contains(followingProfiles, result.PublisherId) {
			if result.PostType == model.GetPostType("story") {
				duration, err := time.ParseDuration("24h")
				if err != nil {
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

func (repo *PostRepository) GetPublic() ([]model.Post, error) {
	collection := repo.getCollection()

	filter := bson.D{{"isdeleted", false}, {"isprivate", false}}

	cursor, err := collection.Find(context.TODO(), filter)
	if err != nil {
		return nil, err
	}

	var posts []model.Post

	for cursor.Next(context.TODO()) {
		var result model.Post
		err = cursor.Decode(&result)
		if result.PostType == model.GetPostType("story") {
			duration, err := time.ParseDuration("24h")
			if err != nil {
				return nil, err
			}
			if time.Now().Before(result.PublishDate.Add(duration)) {
				posts = append(posts, result)
			}
		} else {
			posts = append(posts, result)
		}
	}

	return posts, nil
}

func (repo *PostRepository) GetPublicPostByLocation(location string) ([]model.Post, error) {
	collection := repo.getCollection()

	filter := bson.D{{"isdeleted", false}, {"isprivate", false}}

	cursor, err := collection.Find(context.TODO(), filter)
	if err != nil {
		return nil, err
	}

	var posts []model.Post

	for cursor.Next(context.TODO()) {
		var result model.Post
		err = cursor.Decode(&result)
		if strings.Contains(result.Location, location) {
			if result.PostType == model.GetPostType("story") {
				duration, err := time.ParseDuration("24h")
				if err != nil {
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

func (repo *PostRepository) GetPublicPostByHashTag(hashTag string) ([]model.Post, error) {
	collection := repo.getCollection()

	filter := bson.D{{"isdeleted", false}, {"isprivate", false}}

	cursor, err := collection.Find(context.TODO(), filter)
	if err != nil {
		return nil, err
	}

	var posts []model.Post

	for cursor.Next(context.TODO()) {
		var result model.Post
		err = cursor.Decode(&result)
		if strings.Contains(result.HashTags, hashTag) {
			if result.PostType == model.GetPostType("story") {
				duration, err := time.ParseDuration("24h")
				if err != nil {
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

func (repo *PostRepository) GetMyPosts(loggedUserId uint) ([]model.Post, error) {
	collection := repo.getCollection()

	filter := bson.D{{"isdeleted", false}, {"publisherid", loggedUserId}}

	cursor, err := collection.Find(context.TODO(), filter)
	if err != nil {
		return nil, err
	}

	var posts []model.Post

	for cursor.Next(context.TODO()) {
		var result model.Post
		err = cursor.Decode(&result)
		posts = append(posts, result)
	}

	return posts, nil
}

func (repo *PostRepository) GetPostsForHomePage(followingProfiles []uint) ([]model.Post, error) {
	collection := repo.getCollection()

	cursor, err := collection.Find(context.TODO(), bson.D{})
	if err != nil {
		return nil, err
	}

	var posts []model.Post

	for cursor.Next(context.TODO()) {
		var result model.Post
		err = cursor.Decode(&result)
		if util.Contains(followingProfiles, result.PublisherId) {
			if result.PostType == model.GetPostType("story") {
				duration, err := time.ParseDuration("24h")
				if err != nil {
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

func (repo *PostRepository) Create(post *model.Post) error {
	collection := repo.getCollection()
	_, err := collection.InsertOne(context.TODO(), post)
	return err
}

func (repo *PostRepository) Read(id primitive.ObjectID) (model.Post, error) {
	collection := repo.getCollection()
	filter := bson.D{{"_id", id}, {"isdeleted", false}}
	var result model.Post
	err := collection.FindOne(context.TODO(), filter).Decode(&result)
	return result, err
}

func (repo *PostRepository) Delete(id primitive.ObjectID) error {
	collection := repo.getCollection()
	filter := bson.D{{"_id", id}}
	update := bson.D{
		{"$set", bson.D{
			{"isDeleted", true},
		}},
	}
	result, _ := collection.UpdateOne(context.TODO(), filter, update)
	if result.MatchedCount == 0 {
		return mongo.ErrNoDocuments
	}
	return nil
}

func (repo *PostRepository) Update(id primitive.ObjectID, post dto.PostDto) error {

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
		return mongo.ErrNoDocuments
	}
	return nil
}

func (repo *PostRepository) DeleteUserPosts(id uint) error {
	filter := bson.D{{"publisherid", id}}
	update := bson.D{
		{"$set", bson.D{
			{"isdeleted", true},
		}},
	}

	return repo.updateMany(filter, update)
}

func (repo *PostRepository) ChangeUsername(id uint, username string) error {
	filter := bson.D{{"publisherid", id}}
	update := bson.D{
		{"$set", bson.D{
			{"publisherusername", username},
		}},
	}

	return repo.updateMany(filter, update)
}

func (repo *PostRepository) ChangePrivacy(id uint, isPrivate bool) error {
	filter := bson.D{{"publisherid", id}}
	update := bson.D{
		{"$set", bson.D{
			{"isprivate", isPrivate},
		}},
	}

	return repo.updateMany(filter, update)
}

func (repo *PostRepository) updateMany(filter bson.D, update bson.D) error {
	collection := repo.getCollection()

	result, _ := collection.UpdateMany(context.TODO(), filter, update)

	if result.MatchedCount == 0 {
		return mongo.ErrNoDocuments
	}
	return nil
}

func (repo *PostRepository) getCollection() *mongo.Collection {
	return repo.Client.Database(postDbName).Collection(postsCollectionName)
}
