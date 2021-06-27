package service

import (
	"encoding/json"
	"fmt"
	"gorm.io/gorm"
	"io"
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
	//TODO: dobavi interese
	//TODO: dobavi usere
	//TODO: tip kampanje
	campaignParams := model.CampaignParameters{
		Model:            gorm.Model{},
		Start:            campaignRequest.Start,
		End:              campaignRequest.End,
		CampaignID:       0,
		Interests:        service.getInterestsFromRequest(campaignRequest.Interests),
		CampaignRequests: getCampaignRequestsFromUsernames(campaignRequest.InfluencerUsernames),
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

func getCampaignRequestsFromUsernames(usernames []string) []model.CampaignRequest{
	ids, err := getProfileIdsByUsernames(usernames)

	if err != nil {
		fmt.Println(err)
		return make([]model.CampaignRequest,0)
	}
	ret := make([]model.CampaignRequest,0)
	for _ , value := range ids{
		ret = append(ret,model.CampaignRequest{
			Model:                gorm.Model{},
			InfluencerID:         value,
			RequestStatus:        model.SENT,
			CampaignParametersID: 0,
		})
	}
	return ret
}

func (service *CampaignService) getInterestsFromRequest(interests []string) []model.Interest{
	return service.CampaignRepository.GetInterests(interests)
}

func getProfileIdsByUsernames(usernames []string) ([]uint, error) {
	profileHost, profilePort := util.GetProfileHostAndPort()

	type data struct {
		Usernames []string `json:"usernames"`
	}

	bodyData := data{Usernames: usernames}
	jsonBody, err := json.Marshal(bodyData)
	if err != nil {
		return nil, err
	}

	resp, err := util.CrossServiceRequest(http.MethodPost,
		util.GetCrossServiceProtocol()+"://"+profileHost+":"+profilePort+"/get-by-usernames",
		jsonBody,  map[string]string{"Content-Type": "application/json;"})

	if err != nil {
		return nil, err
	}
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	defer func(Body io.ReadCloser) {
		_ = Body.Close()
	}(resp.Body)
	var ids []string
	if err = json.Unmarshal(body, &ids); err != nil {
		return nil, err
	}
	var ret []uint
	for _,value := range ids {
		ret = append(ret,util.String2Uint(value))
	}

	return ret, err
}

func getCampaignTypeFromRequest(start time.Time,end time.Time, timestampsLength int) model.CampaignType{
	if start.Equal(end) && timestampsLength == 1 {
		return model.ONE_TIME
	}else {
		return model.REPEATABLE
	}
}



