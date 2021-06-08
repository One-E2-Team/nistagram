package dto

type UpdateUserDTO struct {
	ProfileId string `json:"profileId"`
	Email     string `json:"email"`
	Username  string `json:"username"`
}
