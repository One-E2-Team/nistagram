package dto

import (
	"nistagram/campaign/model"
)

type CampaignWithPostDTO struct {
	Campaign model.Campaign `json:"campaign"`
	Post     PostDTO        `json:"post"`
}
