package model

import "gorm.io/gorm"

type CampaignRequest struct {
	gorm.Model
	InfluencerID         uint          `json:"influencerId" gorm:"not null"`
	RequestStatus        RequestStatus `json:"request_status" gorm:"not null"`
	CampaignParametersID uint          `json:"campaignParametersId"`
}
