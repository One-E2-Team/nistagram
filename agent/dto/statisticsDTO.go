package dto

import "time"

type StatisticsDTO struct {
	CampaignId		uint				  `json:"campaignId"`
	Campaign 		CampaignMonitoringDTO `json:"campaign"`
	Events 	 		[]ShowEventDTO		  `json:"events"`
}

type CampaignMonitoringDTO struct {
	PostID             string               	`json:"postId"`
	AgentID            uint                 	`json:"agentId"`
	CampaignType       string               	`json:"campaignType"`
	Start              time.Time            	`json:"start"`
	CampaignParameters []CampaignParametersMonitoringDTO  `json:"campaignParameters"`
}

type CampaignParametersMonitoringDTO struct{
	Start            time.Time            `json:"start"`
	End              time.Time            `json:"end"`
	Interests        []string             `json:"interests"`
	CampaignRequests []CampaignRequestDTO `json:"campaignRequests"`
	Timestamps  	 []time.Time	      `json:"timestamps"`
}

type CampaignRequestDTO struct {
	InfluencerID            uint          `json:"influencerId"`
	InfluencerUsername		string		  `json:"influencerUsername"`
	RequestStatus           string        `json:"request_status"`
}

type ShowEventDTO struct {
	EventType 	    		string 		`json:"eventType"`
	InfluencerId			uint		`json:"influencerId"`
	InfluencerUsername		string		`json:"influencerUsername"`
	Interests			    []string	`json:"interests"`
	WebSite 				string		`json:"webSite"`
	Timestamp   		    time.Time 	`json:"timestamp"`
}
