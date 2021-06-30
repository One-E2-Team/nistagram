package dto

type PersonalDataDTO struct {
	InterestedIn []Interest `json:"interestedIn"`
}

type Interest struct {
	Name string `json:"name"`
}
