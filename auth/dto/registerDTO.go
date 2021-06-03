package dto

type RegisterDTO struct {
	ProfileIdString string `json:"profileId"`
	Email           string `json:"email"`
	Password        string `json:"password"`
}
