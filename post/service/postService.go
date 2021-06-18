package service

import (
	"encoding/json"
	"errors"
	"fmt"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"io"
	"net/http"
	"nistagram/post/dto"
	"nistagram/post/model"
	"nistagram/post/repository"
	"nistagram/util"
	"strings"
	"time"
)

type PostService struct {
	PostRepository *repository.PostRepository
}

func (service *PostService) GetAll() []model.Post {
	return service.PostRepository.GetAll()
}

func (service *PostService) GetPublic(loggedUserID uint) ([]dto.ResponsePostDTO, error) {
	posts := service.PostRepository.GetPublic()
	responseDTO, err := getReactionsForPosts(posts, loggedUserID)
	return responseDTO, err
}

func (service *PostService) GetProfilesPosts(followingProfiles []uint, targetUsername string, loggedUserID uint) ([]dto.ResponsePostDTO, error) {
	posts := service.PostRepository.GetProfilesPosts(followingProfiles, targetUsername)
	responseDTO, err := getReactionsForPosts(posts, loggedUserID)
	return responseDTO, err
}

func (service *PostService) GetPublicPostByLocation(location string, loggedUserID uint) ([]dto.ResponsePostDTO, error) {
	posts := service.PostRepository.GetPublicPostByLocation(location)
	responseDTO, err := getReactionsForPosts(posts, loggedUserID)
	return responseDTO, err
}

func (service *PostService) GetPublicPostByHashTag(hashTag string, loggedUserID uint) ([]dto.ResponsePostDTO, error) {
	posts := service.PostRepository.GetPublicPostByHashTag(hashTag)
	responseDTO, err := getReactionsForPosts(posts, loggedUserID)
	return responseDTO, err
}

func (service *PostService) GetMyPosts(loggedUserID uint) ([]dto.ResponsePostDTO, error) {
	posts := service.PostRepository.GetMyPosts(loggedUserID)
	responseDTO, err := getReactionsForPosts(posts, loggedUserID)
	return responseDTO, err
}

func (service *PostService) GetPostsForHomePage(followingProfiles []uint, loggedUserID uint) ([]dto.ResponsePostDTO, error) {
	posts := service.PostRepository.GetPostsForHomePage(followingProfiles)
	responseDTO, err := getReactionsForPosts(posts, loggedUserID)
	return responseDTO, err
}

func (service *PostService) CreatePost(postType model.PostType, post dto.PostDto, mediaNames []string, profile dto.ProfileDto) error {
	var medias []model.Media
	for i := 0; i < len(mediaNames); i++ {
		m := model.Media{FilePath: mediaNames[i], WebSite: ""}
		medias = append(medias, m)
	}

	if strings.Contains(post.Description, "@"){
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

			if !util.Contains(followingProfiles, taggedProfile.ProfileId){
				return errors.New(taggedProfile.Username + " is not followed by this profile!")
			}
		}
	}
	return nil
}

func getProfileByUsername(username string) (*http.Response, error) {
	profileHost, profilePort := util.GetProfileHostAndPort()
	resp, err := util.CrossServiceRequest(http.MethodGet,
		util.CrossServiceProtocol+"://"+profileHost+":"+profilePort+"/get/" + username,
		nil, map[string]string{})
	return resp, err
}

func getUserFollowers(loggedUserId uint) (*http.Response, error) {
	connHost, connPort := util.GetConnectionHostAndPort()
	resp, err := util.CrossServiceRequest(http.MethodGet,
		util.CrossServiceProtocol+"://"+connHost+":"+connPort+"/connection/following/show/"+util.Uint2String(loggedUserId),
		nil, map[string]string{})
	return resp, err
}

func (service *PostService) ReadPost(id primitive.ObjectID, postType model.PostType) (model.Post, error) {
	return service.PostRepository.Read(id, postType)
}

func (service *PostService) DeletePost(id primitive.ObjectID, postType model.PostType) error {
	return service.PostRepository.Delete(id, postType)
}

func (service *PostService) UpdatePost(id primitive.ObjectID, postType model.PostType, post dto.PostDto) error {
	return service.PostRepository.Update(id, postType, post)
}

func (service *PostService) DeleteUserPosts(profileId uint) error {
	return service.PostRepository.DeleteUserPosts(profileId)
}

func (service *PostService) ChangeUsername(profileId uint, username string) error {
	return service.PostRepository.ChangeUsername(profileId, username)
}

func (service *PostService) ChangePrivacy(profileId uint, isPrivate bool) error {
	return service.PostRepository.ChangePrivacy(profileId, isPrivate)
}

func getReactionsForPosts(posts []model.Post, profileID uint) ([]dto.ResponsePostDTO, error) {
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
		util.CrossServiceProtocol+"://"+postReactionHost+":"+postReactionPort+"/get-reaction-types/" + util.Uint2String(profileID),
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
