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

func (service *PostService) CreatePost(postType model.PostType,post dto.PostDto) error {
	newPost := model.Post{PublisherId: 123, PublisherUsername: "andrej",
		PostType: model.PostType(postType), Medias: nil, PublishDate: time.Now(),
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

func (service *PostService) UpadtePost(id primitive.ObjectID,postType model.PostType,post dto.PostDto) error {
	return service.PostRepository.Update(id,postType,post)
}