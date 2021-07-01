package dto

import "time"

type CampaignRequestForInfluencerDTO struct {
	CampaignId  	uint			`json:"campaign_id"`
	RequestId 		uint			`json:"request_id"`
	Post            PostDTO         `json:"post"`
	Timestamps      []time.Time     `json:"timestamps"`
	Start       	time.Time       `json:"start"`
	End             time.Time       `json:"end"`
}
