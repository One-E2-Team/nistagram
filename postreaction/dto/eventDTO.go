package dto

type EventDTO struct {
	Type 	  	string 		`json:"type"`
	PostId	  	string		`json:"postId"`
	ProfileId	uint		`json:"profileId"`
}
