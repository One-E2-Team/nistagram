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

func (service *PostService) CreatePost(postType model.PostType,post dto.PostDto, mediaNames []string) error {
	var medias []model.Media
	for i:=0;i<len(mediaNames);i++ {
		m := model.Media{FilePath: mediaNames[i], WebSite: ""}
		medias = append(medias, m)
	}

	newPost := model.Post{PublisherId: 123, PublisherUsername: "andrej",
		PostType: postType, Medias: medias, PublishDate: time.Now(),
		Description: post.Description, IsHighlighted: post.IsHighlighted, IsCampaign: post.IsCampaign,
		IsCloseFriendsOnly: post.IsCloseFriendsOnly,
		HashTags: nil, TaggedUsers: post.TaggedUsers, IsPrivate: false, IsDeleted: false}

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