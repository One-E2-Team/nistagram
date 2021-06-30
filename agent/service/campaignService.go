package service

import (
	"net/http"
	"nistagram/agent/util"
)

type CampaignService struct {
	
}

func (service *CampaignService) GetMyCampaigns(agentID uint) (*http.Response, error) {
	return util.NistagramRequest(http.MethodGet, "/agent-api/campaign/my-campaigns", nil, map[string]string{})
}