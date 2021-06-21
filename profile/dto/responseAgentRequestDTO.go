package dto

type ResponseAgentRequestDTO struct {
	Username  string `json:"username"`
	ProfileID uint   `json:"profileId"`
	Email     string `json:"email"`
	Website   string `json:"website"`
}

