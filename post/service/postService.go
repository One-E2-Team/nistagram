package service

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/go-redis/redis/v8"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"io"
	"net/http"
	"nistagram/post/dto"
	"nistagram/post/model"
	"nistagram/post/repository"
	"nistagram/profile/saga"
	"nistagram/util"
	"os"
	"strings"
	"time"
)

type PostService struct {
	PostRepository *repository.PostRepository
}

func (service *PostService) GetPublic(loggedUserID uint) ([]dto.ResponsePostDTO, error) {
	posts, err := service.PostRepository.GetPublic()
	if err != nil{
		return nil, err
	}
	responseDTO, err := getReactionsForPosts(posts, loggedUserID)
	return responseDTO, err
}

func (service *PostService) GetProfilesPosts(followingProfiles []uint, targetUsername string, loggedUserID uint) ([]dto.ResponsePostDTO, error) {
	posts, err := service.PostRepository.GetProfilesPosts(followingProfiles, targetUsername)
	if err != nil{
		return nil, err
	}
	responseDTO, err := getReactionsForPosts(posts, loggedUserID)
	return responseDTO, err
}

func (service *PostService) GetPublicPostByLocation(location string, loggedUserID uint) ([]dto.ResponsePostDTO, error) {
	posts, err := service.PostRepository.GetPublicPostByLocation(location)
	if err != nil{
		return nil, err
	}
	responseDTO, err := getReactionsForPosts(posts, loggedUserID)
	return responseDTO, err
}

func (service *PostService) GetPublicPostByHashTag(hashTag string, loggedUserID uint) ([]dto.ResponsePostDTO, error) {
	posts, err := service.PostRepository.GetPublicPostByHashTag(hashTag)
	if err != nil{
		return nil, err
	}
	responseDTO, err := getReactionsForPosts(posts, loggedUserID)
	return responseDTO, err
}

func (service *PostService) GetMyPosts(loggedUserID uint) ([]dto.ResponsePostDTO, error) {
	posts, err := service.PostRepository.GetMyPosts(loggedUserID)
	if err != nil{
		return nil, err
	}
	responseDTO, err := getReactionsForPosts(posts, loggedUserID)
	return responseDTO, err
}

func (service *PostService) GetPostsForHomePage(followingProfiles []uint, loggedUserID uint) ([]dto.ResponsePostDTO, error) {
	posts, err := service.PostRepository.GetPostsForHomePage(followingProfiles)
	if err != nil{
		return nil, err
	}
	responseDTO, err := getReactionsForPosts(posts, loggedUserID)
	return responseDTO, err
}

func (service *PostService) CreatePost(postType model.PostType, post dto.PostDto, mediaNames []string, profile dto.ProfileDto) error {
	var medias []model.Media
	for i := 0; i < len(mediaNames); i++ {
		m := model.Media{FilePath: mediaNames[i], WebSite: ""}
		medias = append(medias, m)
	}

	if strings.Contains(post.Description, "@") {
		err := canUsersBeTagged(post.Description, profile.ProfileId)
		if err != nil {
			return err
		}
	}

	newPost := model.Post{ID: primitive.NewObjectID(), PublisherId: profile.ProfileId, PublisherUsername: profile.Username,
		PostType: postType, Medias: medias, PublishDate: time.Now(),
		Description: post.Description, IsHighlighted: post.IsHighlighted, IsCampaign: false,
		IsCloseFriendsOnly: post.IsCloseFriendsOnly, Location: post.Location,
		HashTags: post.HashTags, IsPrivate: profile.ProfileSettings.IsPrivate, IsDeleted: false}

	return service.PostRepository.Create(&newPost)
}

func (service *PostService) ReadPost(id primitive.ObjectID) (model.Post, error) {
	return service.PostRepository.Read(id)
}

func (service *PostService) DeletePost(id primitive.ObjectID) error {
	err := service.PostRepository.Delete(id)
	if err != nil{
		return err
	}

	err = deletePostsReports(id)
	return err
}

func (service *PostService) UpdatePost(id primitive.ObjectID, post dto.PostDto) error {
	return service.PostRepository.Update(id, post)
}

func (service *PostService) DeleteUserPosts(profileId uint) error {
	posts, err := service.PostRepository.GetMyPosts(profileId)
	if err != nil {
		return err
	}
	for i := 0; i < len(posts); i++{
		err = service.DeletePost(posts[i].ID)
		if err != nil{
			return err
		}
	}
	return nil
}

func (service *PostService) ChangeUsername(profileId uint, username string) error {
	return service.PostRepository.ChangeUsername(profileId, username)
}

func (service *PostService) ChangePrivacy(profileId uint, isPrivate bool) error {
	return service.PostRepository.ChangePrivacy(profileId, isPrivate)
}

func deletePostsReports(postId primitive.ObjectID) error {
	postReactionHost, postReactionPort := util.GetPostReactionHostAndPort()
	_, err := util.CrossServiceRequest(http.MethodDelete,
		util.CrossServiceProtocol+"://"+postReactionHost+":"+postReactionPort+"/report/"+ util.GetStringIDFromMongoID(postId),
		nil, map[string]string{})
	return err
}

func canUsersBeTagged(description string, publisherId uint) error {
	var followingProfiles []uint

	resp, err := getUserFollowers(publisherId)
	if err != nil {
		return err
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return err
	}
	defer func(Body io.ReadCloser) {
		_ = Body.Close()
	}(resp.Body)

	if err = json.Unmarshal(body, &followingProfiles); err != nil {
		return err
	}

	descriptionParts := strings.Split(description, " ")
	for i := 0; i < len(descriptionParts); i++ {
		if strings.Contains(descriptionParts[i], "@") {
			taggedUsername := descriptionParts[i][1 : len(descriptionParts[i])-1]
			var taggedProfile dto.ProfileDto
			if resp, err := getProfileByUsername(taggedUsername); err != nil {
				return err
			} else {
				body, _ := io.ReadAll(resp.Body)
				defer func(Body io.ReadCloser) {
					_ = Body.Close()
				}(resp.Body)
				if err := json.Unmarshal(body, &taggedProfile); err != nil {
					return err
				}
			}
			if !taggedProfile.ProfileSettings.CanBeTagged {
				return errors.New(taggedProfile.Username + " can't be tagged!")
			}

			if !util.Contains(followingProfiles, taggedProfile.ProfileId) {
				return errors.New(taggedProfile.Username + " is not followed by this profile!")
			}
		}
	}
	return nil
}

func getUserFollowers(loggedUserId uint) (*http.Response, error) {
	connHost, connPort := util.GetConnectionHostAndPort()
	resp, err := util.CrossServiceRequest(http.MethodGet,
		util.CrossServiceProtocol+"://"+connHost+":"+connPort+"/connection/following/show/"+util.Uint2String(loggedUserId),
		nil, map[string]string{})
	return resp, err
}

func getProfileByUsername(username string) (*http.Response, error) {
	profileHost, profilePort := util.GetProfileHostAndPort()
	resp, err := util.CrossServiceRequest(http.MethodGet,
		util.CrossServiceProtocol+"://"+profileHost+":"+profilePort+"/get/"+username,
		nil, map[string]string{})
	return resp, err
}

func getReactionsForPosts(posts []model.Post, profileID uint) ([]dto.ResponsePostDTO, error) {
	if len(posts) == 0 {
		return make([]dto.ResponsePostDTO, 0), nil
	}
	if profileID == 0 {
		ret := make([]dto.ResponsePostDTO, len(posts))
		for i, value := range posts {
			ret[i].Reaction = "none"
			ret[i].Post = value
		}
		return ret, nil
	}
	postIDs := make([]string, 0)
	for _, value := range posts {
		postIDs = append(postIDs, util.GetStringIDFromMongoID(value.ID))
	}
	postReactionHost, postReactionPort := util.GetPostReactionHostAndPort()
	postBody, _ := json.Marshal(map[string][]string{
		"ids": postIDs,
	})
	resp, err := util.CrossServiceRequest(http.MethodPost,
		util.CrossServiceProtocol+"://"+postReactionHost+":"+postReactionPort+"/get-reaction-types/"+util.Uint2String(profileID),
		postBody, map[string]string{"Content-Type": "application/json;"})

	if err != nil {
		return nil, err
	}

	if resp.StatusCode != 200 {
		return nil, fmt.Errorf("BAD_POST_ID")
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	defer func(Body io.ReadCloser) {
		_ = Body.Close()
	}(resp.Body)

	var reactions []string
	if err = json.Unmarshal(body, &reactions); err != nil {
		return nil, err
	}
	if len(posts) != len(reactions) {
		return nil, fmt.Errorf("BAD_LISTS")
	}
	ret := make([]dto.ResponsePostDTO, 0)

	for i, value := range posts {
		ret = append(ret, dto.ResponsePostDTO{
			Post:     value,
			Reaction: reactions[i],
		})
	}
	return ret, nil
}

func (service *PostService) ConnectToRedis(){
	var (
		client *redis.Client
		err error
	)
	time.Sleep(5 * time.Second)
	var redisHost, redisPort = "localhost", "6379"          // dev.db environment
	_, ok := os.LookupEnv("DOCKER_ENV_SET_PROD")        // production environment
	if ok {
		redisHost = "message_broker"
		redisPort = "6379"
	} else {
		_, ok := os.LookupEnv("DOCKER_ENV_SET_DEV") // dev front environment
		if ok {
			redisHost = "message_broker"
			redisPort = "6379"
		}
	}
	for {
		client = redis.NewClient(&redis.Options{
			Addr:     redisHost + ":" + redisPort,
			Password: "",
			DB:       0,
		})

		if err := client.Ping(context.TODO()).Err(); err != nil {
			fmt.Println("Cannot connect to redis! Sleeping 10s and then retrying....")
			time.Sleep(10 * time.Second)
		} else {
			fmt.Println("Post connected to redis.")
			break
		}
	}

	pubsub := client.Subscribe(context.TODO(),saga.PostChannel, saga.ReplyChannel)

	if _, err = pubsub.Receive(context.TODO()); err != nil {
		fmt.Println(err)
		return
	}
	defer func() { _ = pubsub.Close() }()
	ch := pubsub.Channel()

	fmt.Println("Starting post saga in go routine..")

	for{
		select{
		case msg := <-ch:
			m := saga.Message{}
			if err = json.Unmarshal([]byte(msg.Payload), &m); err != nil {
				fmt.Println(err)
				continue
			}

			switch msg.Channel {
			case saga.PostChannel:
				if m.Action == saga.ActionStart {
					switch m.Functionality{
					case saga.ChangeProfilesPrivacy:
						err = service.ChangePrivacy(m.Profile.ID, m.Profile.ProfileSettings.IsPrivate)
						if err != nil{
							sendToReplyChannel(client, &m, saga.ActionError, saga.ProfileService, saga.PostService)
						}else{
							sendToReplyChannel(client, &m, saga.ActionDone, saga.ProfileService, saga.PostService)
						}
					}
				}
			}
		}

	}
}

func sendToReplyChannel(client *redis.Client, m *saga.Message, action string, nextService string, senderService string){
	var err error
	m.Action = action
	m.NextService = nextService
	m.SenderService = senderService
	if err = client.Publish(context.TODO(),saga.ReplyChannel, m).Err(); err != nil {
		fmt.Printf("Error publishing done-message to %s channel", saga.ReplyChannel)
	}
	fmt.Printf("Done message published to channel :%s", saga.ReplyChannel)
}
