package dto

type EventDTO struct {
	EventType 	    string 		`json:"eventType"`
	PostId	  		string		`json:"postId"`
	ProfileId		uint		`json:"profileId"`
	CampaignId		uint		`json:"campaignId"`
	InfluencerId	uint		`json:"influencerId"`
	WebSite 		string		`json:"webSite"`
}
