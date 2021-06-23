package model

import "gorm.io/gorm"

type CampaignStat struct {
	gorm.Model
	CampaignID          uint        		 `json:"campaignID"`
	PostLink            string               `json:"postLink"`
	InfluencerStat      []InfluencerStat     `json:"influencerStat"`
	InterestStat        []InterestStat       `json:"interestStat"`
	StatisticsID        uint        	     `json:"statisticsID"`
}

type InfluencerStat struct {
	gorm.Model
	CampaignStatID			 uint		 `json:"campaignStatID"`
	InfluencerUsername       string      `json:"influencerUsername"`
	StatisticsID             uint        `json:"statisticsID"`
}

type InterestStat struct {
	gorm.Model
	CampaignStatID			 uint		 `json:"campaignStatID"`
	Interest       			 string      `json:"interest"`
	StatisticsID             uint        `json:"statisticsID"`
}
