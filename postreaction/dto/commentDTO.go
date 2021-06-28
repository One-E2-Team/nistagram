package dto

type CommentDTO struct {
	PostID  	   string `json:"postId"`
	Content 	   string `json:"content"`
	CampaignID     uint `json:"campaignId"`
	InfluencerID   uint `json:"influencerID"`
	InfluencerUsername   string `json:"influencerUsername"`
}
