package service

import (
	"fmt"
	"gorm.io/gorm"
	"net/http"
	"nistagram/campaign/dto"
	"nistagram/campaign/model"
	"nistagram/campaign/repository"
	"nistagram/util"
	"time"
)

type CampaignService struct {
	CampaignRepository *repository.CampaignRepository
}

func (service *CampaignService) CreateCampaign (userId uint, campaignRequest dto.CampaignDTO) (model.Campaign,error){
	err := makeCampaign(campaignRequest.PostID, userId)
	if err != nil {
		return model.Campaign{}, err
	}
	campaignParams := model.CampaignParameters{
		Model:            gorm.Model{},
		Start:            campaignRequest.Start,
		End:              campaignRequest.End,
		CampaignID:       0,
		Interests:        service.getInterestsFromRequest(campaignRequest.Interests),
		CampaignRequests: getCampaignRequestsForProfileId(campaignRequest.InfluencerProfileIds),
		Timestamps:       getTimestampsFromRequest(campaignRequest.Timestamps),
	}

	campaign:=model.Campaign{
		Model:              gorm.Model{},
		PostID:             campaignRequest.PostID,
		AgentID:            userId,
		CampaignType:       getCampaignTypeFromRequest(campaignRequest.Start,campaignRequest.End,len(campaignRequest.Timestamps)),
		Start:              campaignRequest.Start,
		CampaignParameters: []model.CampaignParameters{campaignParams},
	}
	return service.CampaignRepository.CreateCampaign(campaign)
}

func getTimestampsFromRequest(timestamps []time.Time) []model.Timestamp{
	ret:= make([]model.Timestamp,0)
	for _,value:= range timestamps{
		ret = append(ret,model.Timestamp{
			Model:                gorm.Model{},
			CampaignParametersID: 0,
			Time:                 value,
		})
	}
	return ret
}

func getCampaignRequestsForProfileId(profileIds []string) []model.CampaignRequest{
	ret := make([]model.CampaignRequest,0)
	for _ , value := range profileIds{
		ret = append(ret,model.CampaignRequest{
			Model:                gorm.Model{},
			InfluencerID:         util.String2Uint(value),
			RequestStatus:        model.SENT,
			CampaignParametersID: 0,
		})
	}
	return ret
}

func (service *CampaignService) getInterestsFromRequest(interests []string) []model.Interest{
	return service.CampaignRepository.GetInterests(interests)
}

func (service *CampaignService) UpdateCampaignParameters(id uint, params dto.CampaignParametersDTO) error {
	newParams := model.CampaignParameters{
		Model:            gorm.Model{},
		Start:            time.Time{},
		End:              params.End,
		CampaignID:       id,
		Interests:        service.getInterestsFromRequest(params.Interests),
		CampaignRequests: getCampaignRequestsForProfileId(params.InfluencerProfileIds),
		Timestamps:       getTimestampsFromRequest(params.Timestamps),
	}
	return service.CampaignRepository.UpdateCampaignParameters(newParams)
}

func getCampaignTypeFromRequest(start time.Time,end time.Time, timestampsLength int) model.CampaignType{
	if start.Equal(end) && timestampsLength == 1 {
		return model.ONE_TIME
	}else {
		return model.REPEATABLE
	}
}

func makeCampaign(postID string, loggedUserID uint) error {
	postHost, postPort := util.GetPostHostAndPort()
	resp, err := util.CrossServiceRequest(http.MethodPost,
		util.GetCrossServiceProtocol()+"://"+postHost+":"+postPort+"make-campaign/"+postID + "/" + util.Uint2String(loggedUserID),
		nil, map[string]string{})
	if err != nil {
		return err
	}
	if resp.StatusCode != 200 {
		return fmt.Errorf("BAD_POST_ID")
	}
	return nil
}



