package dto

type ReactionDTO struct {
	PostID         string `json:"postId"`
	ReactionType   string `json:"reactionType"`
	CampaignID     string `json:"campaignId"`
	InfluencerID   string `json:"influencerID"`
}
