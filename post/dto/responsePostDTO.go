package dto

import "nistagram/post/model"

type ResponsePostDTO struct {
	Post               model.Post       `json:"post"`
	Reaction 		   string           `json:"reaction"`
	CampaignId		   uint				`json:"campaignId"`
	InfluencerId	   uint				`json:"influencerId"`
	InfluencerUsername string			`json:"influencerUsername"`
}
