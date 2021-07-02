package service

import (
	"context"
	"encoding/json"
	"fmt"
	"gorm.io/gorm"
	"io"
	"net/http"
	"nistagram/campaign/dto"
	"nistagram/campaign/model"
	"nistagram/campaign/repository"
	"nistagram/util"
	"sync"
	"time"
)

type CampaignService struct {
	CampaignRepository *repository.CampaignRepository
}

func (service *CampaignService) CreateCampaign(userId uint, campaignRequest dto.CampaignDTO) (model.Campaign, error) {
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

	campaign := model.Campaign{
		Model:              gorm.Model{},
		PostID:             campaignRequest.PostID,
		AgentID:            userId,
		CampaignType:       getCampaignTypeFromRequest(campaignRequest.Start, campaignRequest.End, len(campaignRequest.Timestamps)),
		Start:              campaignRequest.Start,
		CampaignParameters: []model.CampaignParameters{campaignParams},
	}
	return service.CampaignRepository.CreateCampaign(campaign)
}

func getTimestampsFromRequest(timestamps []time.Time) []model.Timestamp {
	ret := make([]model.Timestamp, 0)
	for _, value := range timestamps {
		ret = append(ret, model.Timestamp{
			Model:                gorm.Model{},
			CampaignParametersID: 0,
			Time:                 value,
		})
	}
	return ret
}

func getCampaignRequestsForProfileId(profileIds []string) []model.CampaignRequest {
	ret := make([]model.CampaignRequest, 0)
	for _, value := range profileIds {
		ret = append(ret, model.CampaignRequest{
			Model:                gorm.Model{},
			InfluencerID:         util.String2Uint(value),
			RequestStatus:        model.SENT,
			CampaignParametersID: 0,
		})
	}
	return ret
}

func (service *CampaignService) getInterestsFromRequest(interests []string) []model.Interest {
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

func (service *CampaignService) DeleteCampaign(id uint) error {
	return service.CampaignRepository.DeleteCampaign(id)
}

func (service *CampaignService) GetMyCampaigns(agentID uint) ([]dto.CampaignWithPostDTO, error) {
	campaigns, err := service.CampaignRepository.GetMyCampaigns(agentID)
	if err != nil {
		return nil, err
	}
	postIDs := make([]string, 0)
	for _, value := range campaigns {
		postIDs = append(postIDs, value.PostID)
	}

	posts, err := getPostsByPostsIds(postIDs)
	if err != nil {
		return nil, err
	}
	ret := make([]dto.CampaignWithPostDTO, 0)
	for i := 0; i < len(posts); i++ {
		currentPost := posts[i]
		postDTO := dto.PostDTO{
			PostType:           currentPost.PostType,
			Medias:             currentPost.Medias,
			PublishDate:        currentPost.PublishDate,
			Description:        currentPost.Description,
			IsHighlighted:      currentPost.IsHighlighted,
			IsCloseFriendsOnly: currentPost.IsCloseFriendsOnly,
			Location:           currentPost.Location,
			HashTags:           currentPost.HashTags,
			TaggedUsers:        currentPost.TaggedUsers,
			IsPrivate:          currentPost.IsPrivate,
			IsDeleted:          currentPost.IsDeleted,
		}
		ret = append(ret, dto.CampaignWithPostDTO{Campaign: campaigns[i], Post: postDTO})
	}
	return ret, err
}

func (service *CampaignService) GetAllInterests() ([]string, error) {
	return service.CampaignRepository.GetAllInterests()
}

func getCampaignTypeFromRequest(start time.Time, end time.Time, timestampsLength int) model.CampaignType {
	if start.Equal(end) && timestampsLength == 1 {
		return model.ONE_TIME
	} else {
		return model.REPEATABLE
	}
}

func makeCampaign(postID string, loggedUserID uint) error {
	postHost, postPort := util.GetPostHostAndPort()
	resp, err := util.CrossServiceRequest(context.Background(), http.MethodPost,
		util.GetCrossServiceProtocol()+"://"+postHost+":"+postPort+"/make-campaign/"+postID+"/"+util.Uint2String(loggedUserID),
		nil, map[string]string{})
	if err != nil {
		return err
	}
	if resp.StatusCode != 200 {
		return fmt.Errorf("BAD_POST_ID")
	}
	return nil
}

func (service *CampaignService) GetCurrentlyValidInterests(campaignId uint) ([]string, error) {
	var ret []string
	parameters, err := service.CampaignRepository.GetParametersByCampaignId(campaignId)

	for _, i := range parameters.Interests {
		ret = append(ret, i.Name)
	}

	return ret, err
}

func (service *CampaignService) GetLastActiveParametersForCampaign(id uint) (model.CampaignParameters, error) {
	return service.CampaignRepository.GetLastActiveParametersForCampaign(id)
}

func (service *CampaignService) GetCampaignByIdForMonitoring(campaignId uint) (dto.CampaignMonitoringDTO, error) {

	var ret dto.CampaignMonitoringDTO
	var retParams []dto.CampaignParametersMonitoringDTO

	campaign, err := service.CampaignRepository.GetCampaignById(campaignId)
	if err != nil {
		return ret, err
	}

	for _, param := range campaign.CampaignParameters {
		var paramDto dto.CampaignParametersMonitoringDTO
		var interests []string
		var timestamps []time.Time
		var requests []dto.CampaignRequestDTO
		for _, interest := range param.Interests {
			interests = append(interests, interest.Name)
		}
		for _, ts := range param.Timestamps {
			timestamps = append(timestamps, ts.Time)
		}
		for _, request := range param.CampaignRequests {
			reqDto := dto.CampaignRequestDTO{InfluencerID: request.InfluencerID,
				InfluencerUsername: "", RequestStatus: request.RequestStatus.ToString()}
			requests = append(requests, reqDto)
		}
		paramDto.Interests = interests
		paramDto.Timestamps = timestamps
		paramDto.Start = param.Start
		paramDto.End = param.End
		paramDto.CampaignRequests = requests

		retParams = append(retParams, paramDto)
	}

	ret.PostID = campaign.PostID
	ret.AgentID = campaign.AgentID
	ret.Start = campaign.Start
	ret.CampaignType = campaign.CampaignType.ToString()
	ret.CampaignParameters = retParams

	return ret, nil
}

func (service *CampaignService) GetAvailableCampaignsForUser(loggedUserID uint, followingProfiles []util.FollowingProfileDTO) ([]string, []uint, []uint, error) {
	interests, err := getProfileInterests(loggedUserID)
	if err != nil {
		return nil, nil, nil, err
	}
	allActiveParams, err := service.CampaignRepository.GetAllActiveParameters()
	if err != nil {
		return nil, nil, nil, err
	}
	campaignIDs := make([]uint, 0)
	retInfluencerIDs := make([]uint, 0)
	var wg sync.WaitGroup
	for _, params := range allActiveParams {
		wg.Add(2)
		go func() {
			defer wg.Done()
			if campaignParamsContainsInterest(params, interests) {
				campaignIDs = append(campaignIDs, params.CampaignID)
				retInfluencerIDs = append(retInfluencerIDs, 0)
			}
		}()
		go func() {
			defer wg.Done()
			test, influencerIDs := campaignParamsContainsFollowedProfiles(params, followingProfiles)
			if test {
				for _, influencerID := range influencerIDs {
					campaignIDs = append(campaignIDs, params.CampaignID)
					retInfluencerIDs = append(retInfluencerIDs, influencerID)
				}
			}
		}()
		wg.Wait()
	}
	if len(campaignIDs) == 0 {
		return make([]string, 0), make([]uint, 0), make([]uint, 0), nil
	}
	postIDs, err := service.CampaignRepository.GetPostIDsFromCampaignIDs(campaignIDs)
	return postIDs, retInfluencerIDs, campaignIDs, err
}

func (service *CampaignService) UpdateCampaignRequest(loggedUserId uint,requestId string, status model.RequestStatus) error {
	if service.CampaignRepository.GetCampaignRequestInfluencerId(util.String2Uint(requestId)) != loggedUserId {
		return fmt.Errorf("not allowed")
	}
	return service.CampaignRepository.UpdateCampaignRequest(requestId, status)
}

func (service *CampaignService) GetActiveCampaignsRequestsForProfileId(profileId int) error {
	res, err := service.CampaignRepository.GetDistinctCampaignParamsIdForProfileId(profileId)
	if err != nil {
		return err
	}
	res, err = service.CampaignRepository.GetActiveCampaignIdsForCampaignParamsIds(res)
	if err != nil {
		return err
	}
	fmt.Println(res)
	return nil
}

func (service *CampaignService) GetCampaignRequestInfluencerId(requestId uint) uint {
	return service.CampaignRepository.GetCampaignRequestInfluencerId(requestId)
}

func getPostsByPostsIds(postsIds []string) ([]dto.PostDTO, error) {
	var ret []dto.PostDTO
	type data struct {
		Ids []string `json:"ids"`
	}
	bodyData := data{Ids: postsIds}
	jsonBody, err := json.Marshal(bodyData)
	if err != nil {
		return nil, err
	}
	postHost, postPort := util.GetPostHostAndPort()
	resp, err := util.CrossServiceRequest(context.Background(), http.MethodPost,
		util.GetCrossServiceProtocol()+"://"+postHost+":"+postPort+"/posts",
		jsonBody, map[string]string{})

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

	if err = json.Unmarshal(body, &ret); err != nil {
		return nil, err
	}
	return ret, nil
}

func getProfileInterests(loggedUserID uint) ([]string, error) {
	profileHost, profilePort := util.GetProfileHostAndPort()
	resp, err := util.CrossServiceRequest(context.Background(), http.MethodGet,
		util.GetCrossServiceProtocol()+"://"+profileHost+":"+profilePort+"/profile-interests/"+util.Uint2String(loggedUserID),
		nil, map[string]string{})

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
	var ret []string
	if err = json.Unmarshal(body, &ret); err != nil {
		return nil, err
	}
	return ret, nil
}

func campaignParamsContainsInterest(params model.CampaignParameters, interests []string) bool {
	for _, param := range params.Interests {
		for _, interest := range interests {
			if param.Name == interest {
				return true
			}
		}
	}
	return false
}

func campaignParamsContainsFollowedProfiles(params model.CampaignParameters, followingProfiles []util.FollowingProfileDTO) (bool, []uint) {
	influencers := make([]uint, 0)
	contains := false
	for _, campaignRequest := range params.CampaignRequests {
		for _, potentialInfluencer := range followingProfiles {
			if /*campaignRequest.RequestStatus == model.ACCEPTED &&*/ campaignRequest.InfluencerID == potentialInfluencer.ProfileID {
				contains = true
				influencers = append(influencers, campaignRequest.InfluencerID)
			}
		}
	}
	return contains, influencers
}
