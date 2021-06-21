package dto

type ProcessAgentRequest struct {
	ProfileID string `json:"profileId"`
	Accept    bool   `json:"accept"`
}
