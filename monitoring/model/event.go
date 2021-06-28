package model

import (
	"time"
)

type Event struct {
	OriginalPostId        uint      	`json:"originalPostId"`
	WebSite 			  string		`json:"webSite"`
	InfluencerUsername 	  string		`json:"influencerUsername"`
	Timestamp   		  time.Time 	`json:"timestamp"`
	CampaignId        	  uint      	`json:"campaignId"`
	Interests			  []string		`json:"interests"`
}
