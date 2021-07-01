package dto

import "time"

type CampaignParametersDTO struct {
	End                  time.Time   `json:"end"`
	Interests            []string    `json:"interests"`
	InfluencerProfileIds []string    `json:"influencerProfileIds"`
	Timestamps           []time.Time `json:"timestamps"`
}
