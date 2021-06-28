package dto

type CommentDTO struct {
	PostID  	   string `json:"postId"`
	Content 	   string `json:"content"`
	CampaignID     string `json:"campaignId"`
	InfluencerID   string `json:"influencerID"`
	InfluencerUsername   string `json:"influencerUsername"`
}
