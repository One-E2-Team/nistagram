package repository

import (
	"context"
	"errors"
	"fmt"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"nistagram/post/dto"
	"nistagram/post/model"
	"strings"
	"time"
)

type PostRepository struct {
	Client *mongo.Client
}

func (repo *PostRepository) GetAll() []model.Post {
	postCollection, err := repo.getCollection(model.GetPostType("post"))
	if err != nil {
		fmt.Println("Error: can't get post collection!")
	}
	storyCollection, err := repo.getCollection(model.GetPostType("story"))
	if err != nil {
		fmt.Println("Error: can't get story collection!")
	}
	postCursor, err := postCollection.Find(context.TODO(), bson.D{})
	if err != nil {
		fmt.Println("Error: can't find posts!")
	}
	storyCursor, err := storyCollection.Find(context.TODO(), bson.D{})
	if err != nil {
		fmt.Println("Error: can't find stories!")
	}
	var posts []model.Post

	for postCursor.Next(context.TODO()) {
		var result model.Post
		err = postCursor.Decode(&result)
		posts = append(posts, result)
	}

	for storyCursor.Next(context.TODO()) {
		var result model.Post
		err = storyCursor.Decode(&result)
		posts = append(posts, result)
	}

	return posts
}

func (repo *PostRepository) GetProfilesPosts(followingProfiles []uint, targetUsername string) []model.Post {
	postCollection, err := repo.getCollection(model.GetPostType("post"))
	if err != nil {
		fmt.Println("Error: can't get post collection!")
	}
	storyCollection, err := repo.getCollection(model.GetPostType("story"))
	if err != nil {
		fmt.Println("Error: can't get story collection!")
	}
	postCursor, err := postCollection.Find(context.TODO(), bson.D{})
	if err != nil {
		fmt.Println("Error: can't find posts!")
	}
	storyCursor, err := storyCollection.Find(context.TODO(), bson.D{})
	if err != nil {
		fmt.Println("Error: can't find stories!")
	}
	var posts []model.Post

	for postCursor.Next(context.TODO()) {
		var result model.Post
		err = postCursor.Decode(&result)
		if result.PublisherUsername == targetUsername &&
			(result.IsPrivate == false || contains(followingProfiles, result.PublisherId)) {
			posts = append(posts, result)
		}
	}

	for storyCursor.Next(context.TODO()) {
		var result model.Post
		err = storyCursor.Decode(&result)
		if result.PublisherUsername == targetUsername &&
			(result.IsPrivate == false || contains(followingProfiles, result.PublisherId)) {
			duration, err := time.ParseDuration("24h")
			if err != nil {
				fmt.Println(err)
			}
			if result.IsHighlighted || time.Now().Before(result.PublishDate.Add(duration)) {
				posts = append(posts, result)
			}
		}
	}

	return posts
}

func (repo *PostRepository) GetPublic() []model.Post {
	postCollection, err := repo.getCollection(model.GetPostType("post"))
	if err != nil {
		fmt.Println("Error: can't get post collection!")
	}
	storyCollection, err := repo.getCollection(model.GetPostType("story"))
	if err != nil {
		fmt.Println("Error: can't get story collection!")
	}
	postCursor, err := postCollection.Find(context.TODO(), bson.D{})
	if err != nil {
		fmt.Println("Error: can't find posts!")
	}
	storyCursor, err := storyCollection.Find(context.TODO(), bson.D{})
	if err != nil {
		fmt.Println("Error: can't find stories!")
	}
	var posts []model.Post

	for postCursor.Next(context.TODO()) {
		var result model.Post
		err = postCursor.Decode(&result)
		if result.IsPrivate == false {
			posts = append(posts, result)
		}
	}

	for storyCursor.Next(context.TODO()) {
		var result model.Post
		err = storyCursor.Decode(&result)
		if result.IsPrivate == false {
			duration, err := time.ParseDuration("24h")
			if err != nil {
				fmt.Println(err)
			}
			if time.Now().Before(result.PublishDate.Add(duration)) {
				posts = append(posts, result)
			}
		}
	}

	return posts
}

func (repo *PostRepository) GetPublicPostByLocation(location string) []model.Post {
	postCollection, err := repo.getCollection(model.GetPostType("post"))
	if err != nil {
		fmt.Println("Error: can't get post collection!")
	}
	storyCollection, err := repo.getCollection(model.GetPostType("story"))
	if err != nil {
		fmt.Println("Error: can't get story collection!")
	}
	postCursor, err := postCollection.Find(context.TODO(), bson.D{})
	if err != nil {
		fmt.Println("Error: can't find posts!")
	}
	storyCursor, err := storyCollection.Find(context.TODO(), bson.D{})
	if err != nil {
		fmt.Println("Error: can't find stories!")
	}
	var posts []model.Post

	for postCursor.Next(context.TODO()) {
		var result model.Post
		err = postCursor.Decode(&result)
		if result.IsPrivate == false && strings.Contains(result.Location, location) {
			posts = append(posts, result)
		}
	}

	for storyCursor.Next(context.TODO()) {
		var result model.Post
		err = storyCursor.Decode(&result)
		if result.IsPrivate == false && strings.Contains(result.Location, location) {
			duration, err := time.ParseDuration("24h")
			if err != nil {
				fmt.Println(err)
			}
			if time.Now().Before(result.PublishDate.Add(duration)) {
				posts = append(posts, result)
			}
		}
	}

	return posts
}

func (repo *PostRepository) GetPublicPostByHashTag(hashTag string) []model.Post {
	postCollection, err := repo.getCollection(model.GetPostType("post"))
	if err != nil {
		fmt.Println("Error: can't get post collection!")
	}
	storyCollection, err := repo.getCollection(model.GetPostType("story"))
	if err != nil {
		fmt.Println("Error: can't get story collection!")
	}
	postCursor, err := postCollection.Find(context.TODO(), bson.D{})
	if err != nil {
		fmt.Println("Error: can't find posts!")
	}
	storyCursor, err := storyCollection.Find(context.TODO(), bson.D{})
	if err != nil {
		fmt.Println("Error: can't find stories!")
	}
	var posts []model.Post

	for postCursor.Next(context.TODO()) {
		var result model.Post
		err = postCursor.Decode(&result)
		if result.IsPrivate == false && strings.Contains(result.HashTags, hashTag) {
			posts = append(posts, result)
		}
	}

	for storyCursor.Next(context.TODO()) {
		var result model.Post
		err = storyCursor.Decode(&result)
		if result.IsPrivate == false && strings.Contains(result.HashTags, hashTag) {
			duration, err := time.ParseDuration("24h")
			if err != nil {
				fmt.Println(err)
			}
			if time.Now().Before(result.PublishDate.Add(duration)) {
				posts = append(posts, result)
			}
		}
	}

	return posts
}

func (repo *PostRepository) GetMyPosts(loggedUserId uint) []model.Post {
	postCollection, err := repo.getCollection(model.GetPostType("post"))
	if err != nil {
		fmt.Println("Error: can't get post collection!")
	}
	storyCollection, err := repo.getCollection(model.GetPostType("story"))
	if err != nil {
		fmt.Println("Error: can't get story collection!")
	}
	postCursor, err := postCollection.Find(context.TODO(), bson.D{})
	if err != nil {
		fmt.Println("Error: can't find posts!")
	}
	storyCursor, err := storyCollection.Find(context.TODO(), bson.D{})
	if err != nil {
		fmt.Println("Error: can't find stories!")
	}
	var posts []model.Post

	for postCursor.Next(context.TODO()) {
		var result model.Post
		err = postCursor.Decode(&result)
		if result.PublisherId == loggedUserId {
			posts = append(posts, result)
		}
	}

	for storyCursor.Next(context.TODO()) {
		var result model.Post
		err = storyCursor.Decode(&result)
		if result.PublisherId == loggedUserId {
			posts = append(posts, result)
		}
	}

	return posts
}

func (repo *PostRepository) GetPostsForHomePage(followingProfiles []uint) []model.Post {
	postCollection, err := repo.getCollection(model.GetPostType("post"))
	if err != nil {
		fmt.Println("Error: can't get post collection!")
	}
	storyCollection, err := repo.getCollection(model.GetPostType("story"))
	if err != nil {
		fmt.Println("Error: can't get story collection!")
	}
	postCursor, err := postCollection.Find(context.TODO(), bson.D{})
	if err != nil {
		fmt.Println("Error: can't find posts!")
	}
	storyCursor, err := storyCollection.Find(context.TODO(), bson.D{})
	if err != nil {
		fmt.Println("Error: can't find stories!")
	}
	var posts []model.Post

	for postCursor.Next(context.TODO()) {
		var result model.Post
		err = postCursor.Decode(&result)
		if contains(followingProfiles, result.PublisherId) {
			posts = append(posts, result)
		}
	}

	for storyCursor.Next(context.TODO()) {
		var result model.Post
		err = storyCursor.Decode(&result)
		if contains(followingProfiles, result.PublisherId) {
			duration, err := time.ParseDuration("24h")
			if err != nil {
				fmt.Println(err)
			}
			if time.Now().Before(result.PublishDate.Add(duration)) {
				posts = append(posts, result)
			}
		}
	}

	return posts
}

func contains(array []uint, el uint) bool {
	for _, a := range array {
		if a == el {
			return true
		}
	}
	return false
}

func (repo *PostRepository) Create(post *model.Post) error {
	collection, err := repo.getCollection(post.PostType)
	if err != nil {
		return err
	}
	_, err = collection.InsertOne(context.TODO(), post)
	return err
}

func (repo *PostRepository) Read(id primitive.ObjectID, postType model.PostType) (model.Post, error) {
	collection, err := repo.getCollection(postType)
	if err != nil {
		return model.Post{}, err
	}
	filter := bson.D{{"_id", id}}
	var result model.Post
	err = collection.FindOne(context.TODO(), filter).Decode(&result)
	return result, err
}

func (repo *PostRepository) Delete(id primitive.ObjectID, postType model.PostType) error {
	collection, err := repo.getCollection(postType)
	if err != nil {
		return err
	}
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

func (repo *PostRepository) Update(id primitive.ObjectID, postType model.PostType, post dto.PostDto) error {

	collection, err := repo.getCollection(postType)
	if err != nil {
		return err
	}

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

	return repo.updateManyInPostAndStoryCollections(filter, update)

}

func (repo *PostRepository) ChangeUsername(id uint, username string) error {
	filter := bson.D{{"publisherid", id}}
	update := bson.D{
		{"$set", bson.D{
			{"publisherusername", username},
		}},
	}

	return repo.updateManyInPostAndStoryCollections(filter, update)
}

func (repo *PostRepository) ChangePrivacy(id uint, isPrivate bool) error {
	filter := bson.D{{"publisherid", id}}
	update := bson.D{
		{"$set", bson.D{
			{"isprivate", isPrivate},
		}},
	}

	return repo.updateManyInPostAndStoryCollections(filter, update)
}

func (repo *PostRepository) updateManyInPostAndStoryCollections(filter bson.D, update bson.D) error {
	collPosts, _ := repo.getCollection(model.POST)
	collStories, _ := repo.getCollection(model.STORY)

	result1, _ := collPosts.UpdateMany(context.TODO(), filter, update)
	result2, _ := collStories.UpdateMany(context.TODO(), filter, update)

	if result1.MatchedCount == 0 && result2.MatchedCount == 0 {
		return mongo.ErrNoDocuments
	}
	return nil
}

func (repo *PostRepository) getCollection(postType model.PostType) (*mongo.Collection, error) {
	if postType == model.POST {
		return repo.Client.Database("postdb").Collection("posts"), nil
	} else if postType == model.STORY {
		return repo.Client.Database("postdb").Collection("stories"), nil
	}
	return nil, errors.New("collection doesn't exist!")
}
