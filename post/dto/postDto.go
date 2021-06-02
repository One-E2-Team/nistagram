package dto


type PostDto struct {
	Description string `json:"description"`
	IsHighlighted bool `json:"isHighlighted"`
	IsCampaign bool `json:"isCampaign"`
	IsCloseFriendsOnly bool `json:"isCloseFriendsOnly"`
	Location string `json:"location"`
	HashTags []string `json:"hashTags"`
	TaggedUsers []string `json:"taggedUsers"`
}
