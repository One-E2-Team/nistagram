package dto

import "time"

type CampaignMonitoringDTO struct {
	PostID             string               	`json:"postId"`
	AgentID            uint                 	`json:"agentId"`
	CampaignType       string               	`json:"campaignType"`
	Start              time.Time            	`json:"start"`
	CampaignParameters []CampaignParametersMonitoringDTO  `json:"campaignParameters"`
}

type CampaignParametersMonitoringDTO struct{
	Start            time.Time            `json:"start" gorm:"not null"`
	End              time.Time            `json:"end" gorm:"not null"`
	Interests        []string             `json:"interests"`
	CampaignRequests []CampaignRequestDTO `json:"campaignRequests"`
	Timestamps  	 []time.Time	      `json:"timestamps" gorm:"not null"`
}

type CampaignRequestDTO struct {
	InfluencerID            uint          `json:"influencerId"`
	InfluencerUsername		string		  `json:"influencerUsername"`
	RequestStatus           string        `json:"request_status"`
}