package dto

type SponsoredPostsDTO struct {
	PostID       string `json:"postId"`
	InfluencerID uint   `json:"influencerId"`
	CampaignID   uint   `json:"campaignId"`
}

