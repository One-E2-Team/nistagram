package dto

type ReactionDTO struct {
	PostID         		 string `json:"postId"`
	ReactionType   		 string `json:"reactionType"`
	CampaignID     		 uint `json:"campaignId"`
	InfluencerID   		 uint `json:"influencerID"`
	InfluencerUsername   string `json:"influencerUsername"`
}
