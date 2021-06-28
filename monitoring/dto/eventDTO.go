package dto

type EventDTO struct {
	Type 	  		string 		`json:"type"`
	PostId	  		uint		`json:"postId"`
	ProfileId		uint		`json:"profileId"`
	CampaignId		uint		`json:"campaignId"`
	InfluencerId	uint		`json:"influencerId"`
	InfluencerUsername	string		`json:"influencerUsername"`
	WebSite			string		`json:"webSite"`
}
