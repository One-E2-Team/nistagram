package service

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
	"nistagram/post/dto"
	"nistagram/post/model"
	"nistagram/post/repository"
	"time"
)

type PostService struct {
	PostRepository *repository.PostRepository
}

func (service *PostService) GetAll() ([]model.Post){
	return service.PostRepository.GetAll()
}

func (service *PostService) GetPublic() ([]model.Post){
	return service.PostRepository.GetPublic()
}

func (service *PostService) GetProfilesPosts(followingProfiles []uint, targetUsername string) []model.Post{
	return service.PostRepository.GetProfilesPosts(followingProfiles, targetUsername)
}

func (service *PostService) GetPublicPostByLocation(location string) ([]model.Post){
	return service.PostRepository.GetPublicPostByLocation(location)
}

func (service *PostService) GetPublicPostByHashTag(hashTag string) ([]model.Post){
	return service.PostRepository.GetPublicPostByHashTag(hashTag)
}

func (service *PostService) GetMyPosts(loggedUserId uint) ([]model.Post){
	return service.PostRepository.GetMyPosts(loggedUserId)
}

func (service *PostService) GetPostsForHomePage(followingProfiles []uint) []model.Post{
	return service.PostRepository.GetPostsForHomePage(followingProfiles)
}

func (service *PostService) CreatePost(postType model.PostType,post dto.PostDto, mediaNames []string, profile dto.ProfileDto) error {
	var medias []model.Media
	for i:=0;i<len(mediaNames);i++ {
		m := model.Media{FilePath: mediaNames[i], WebSite: ""}
		medias = append(medias, m)
	}

	newPost := model.Post{PublisherId: profile.ProfileId, PublisherUsername: profile.Username,
		PostType: postType, Medias: medias, PublishDate: time.Now(),
		Description: post.Description, IsHighlighted: post.IsHighlighted, IsCampaign: false,
		IsCloseFriendsOnly: post.IsCloseFriendsOnly, Location: post.Location,
		HashTags: post.HashTags, TaggedUsers: post.TaggedUsers, IsPrivate: profile.ProfileSettings.IsPrivate, IsDeleted: false}

	return service.PostRepository.Create(&newPost)
}

func (service *PostService) ReadPost(id primitive.ObjectID, postType model.PostType) (model.Post,error) {
	return service.PostRepository.Read(id,postType)
}

func (service *PostService) DeletePost(id primitive.ObjectID, postType model.PostType)  error {
	return service.PostRepository.Delete(id,postType)
}

func (service *PostService) UpdatePost(id primitive.ObjectID,postType model.PostType,post dto.PostDto) error {
	return service.PostRepository.Update(id,postType,post)
}

func (service *PostService) DeleteUserPosts(profileId uint) error {
	return service.PostRepository.DeleteUserPosts(profileId)
}

func (service *PostService) ChangeUsername(profileId uint, username string) error {
	return service.PostRepository.ChangeUsername(profileId,username)
}

func (service *PostService) ChangePrivacy(profileId uint, isPrivate bool) error {
	return service.PostRepository.ChangePrivacy(profileId, isPrivate);
}