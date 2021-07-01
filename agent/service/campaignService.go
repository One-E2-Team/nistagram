package service

import (
	"net/http"
	"nistagram/agent/util"
)

type CampaignService struct {
}

func (service *CampaignService) GetMyCampaigns() (*http.Response, error) {
	return util.NistagramRequest(http.MethodGet, "/agent-api/campaign/my-campaigns",
		nil, map[string]string{})
}

func (service *CampaignService) CreateCampaign(requestBody []byte) (*http.Response, error) {
	return util.NistagramRequest(http.MethodPost, "/agent-api/campaign/create",
		requestBody, map[string]string{"Content-Type": "application/json"})
}

func (service *CampaignService) GetInterests() (*http.Response, error) {
	return util.NistagramRequest(http.MethodGet, "/agent-api/campaign/interests",
		nil, map[string]string{})
}

func (service *CampaignService) GetActiveParams(campaignID string) (*http.Response, error) {
	return util.NistagramRequest(http.MethodGet, "/agent-api/campaign/"+campaignID+"/params/active",
		nil, map[string]string{})
}

func (service *CampaignService) EditCampaign(postID string, requestBody []byte) (*http.Response, error) {
	return util.NistagramRequest(http.MethodPut, "/agent-api/campaign/update/"+postID,
		requestBody, map[string]string{"Content-Type": "application/json"})
}
