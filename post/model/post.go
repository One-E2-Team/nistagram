package model

import (
	"time"
)

type Post struct {
	PublisherId uint `json:"publisherId"`
	PublisherUsername string `json:"publisherUsername"`
	PostType PostType `json:"postType"`
 	Medias []Media `json:"medias"`
	PublishDate time.Time `json:"publishDate"`
	Description string `json:"description"`
	IsHighlighted bool `json:"isHighlighted"`
	IsCampaign bool `json:"isCampaign"`
	IsCloseFriendsOnly bool `json:"isCloseFriendsOnly"`
	Location string `json:"location"`
	HashTags string `json:"hashTags"`
	TaggedUsers []string `json:"taggedUsers"`
	IsPrivate bool `json:"isPrivate"`
	IsDeleted bool `json:"isDeleted"`
}

func (post *Post) AddMedia(item Media){
	post.Medias = append(post.Medias, item)
}
