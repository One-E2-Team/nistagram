package dto

type ProfileRecommendationDTO struct {
	Username  	string	`json:"username"`
	ProfileID 	uint	`json:"profileID"`
	Confidence	float64	`json:"confidence"`
}