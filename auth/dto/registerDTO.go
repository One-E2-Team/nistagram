package dto

type RegisterDTO struct {
	ProfileIdString string `json:"profileId"`
	Username        string `json:"username"`
	Email           string `json:"email"`
	Password        string `json:"password"`
}
