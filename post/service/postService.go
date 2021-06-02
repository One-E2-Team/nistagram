package service

import (
	"nistagram/post/dto"
	"nistagram/post/model"
	"nistagram/post/repository"
	"time"
)

type PostService struct {
	PostRepository *repository.PostRepository
}

func (service *PostService) CreatePost(post dto.PostDto) error {
	newPost := model.Post{PublisherId: 123, PublisherUsername: "andrej",
		PostType: model.PostType(post.PostType), Medias: nil, PublishDate: time.Now(),
		Description: post.Description, IsHighlighted: post.IsHighlighted, IsCampaign: post.IsCampaign,
		IsCloseFriendsOnly: post.IsCloseFriendsOnly,
		HashTags: nil, TaggedUsers: post.TaggedUsers, IsPrivate: false, IsDeleted: false}

	return service.PostRepository.CreatePost(&newPost)
}
