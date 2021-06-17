package service

import (
	"encoding/json"
	"errors"
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

func (service *PostService) GetPublic() []model.Post {
	return service.PostRepository.GetPublic()
}

func (service *PostService) GetProfilesPosts(followingProfiles []uint, targetUsername string) []model.Post {
	return service.PostRepository.GetProfilesPosts(followingProfiles, targetUsername)
}

func (service *PostService) GetPublicPostByLocation(location string) []model.Post {
	return service.PostRepository.GetPublicPostByLocation(location)
}

func (service *PostService) GetPublicPostByHashTag(hashTag string) []model.Post {
	return service.PostRepository.GetPublicPostByHashTag(hashTag)
}

func (service *PostService) GetMyPosts(loggedUserId uint) []model.Post {
	return service.PostRepository.GetMyPosts(loggedUserId)
}

func (service *PostService) GetPostsForHomePage(followingProfiles []uint) []model.Post {
	return service.PostRepository.GetPostsForHomePage(followingProfiles)
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
	defer resp.Body.Close()

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
				defer resp.Body.Close()
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
