package service

import "nistagram/campaign/repository"

type CampaignService struct {
	CampaignRepository *repository.CampaignRepository
}
