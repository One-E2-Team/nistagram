package model

import (
	"time"
)

type Event struct {
	OriginalPostId        string      	`json:"originalPostId"`
	EventType			  EventType		`json:"eventType"`
	InfluencerUsername 	  string		`json:"influencerUsername"`
	Timestamp   		  time.Time 	`json:"timestamp"`
	CampaignId        	  uint      	`json:"campaignId"`
	Interests			  []string		`json:"interests"`
	WebSite 			  string		`json:"webSite"`
}

func (event *Event) AddInterest(i string) {
	event.Interests = append(event.Interests, i)
}
