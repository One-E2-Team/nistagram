package model

import (
	"time"
)

type Event struct {
	OriginalPostId        uint      	`json:"originalPostId"`
	EventType			  EventType		`json:"eventType"`
	InfluencerUsername 	  string		`json:"influencerUsername"`
	Timestamp   		  time.Time 	`json:"timestamp"`
	CampaignId        	  uint      	`json:"campaignId"`
	Interests			  []string		`json:"interests"`
	WebSite 			  string		`json:"webSite"`
}
