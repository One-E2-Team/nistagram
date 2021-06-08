package dto

type RegisterDTO struct {
	ProfileIdString string `json:"profileId"`
	Email           string `json:"email"`
	Username        string `json:"username"`
	Password        string `json:"password"`
}
