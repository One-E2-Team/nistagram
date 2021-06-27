package dto

import "time"

type CampaignDTO struct {
	PostID string `json:"postId"`
	Start time.Time `json:"start"`
	End time.Time `json:"end"`
	Interests []string `json:"interests"`
	InfluencerUsernames []string `json:"influencerUsernames"`
	Timestamps []time.Time `json:"timestamps"`
}
