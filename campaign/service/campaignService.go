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

func (service *CampaignService) CreateCampaign(ctx context.Context, userId uint, campaignRequest dto.CampaignDTO) (model.Campaign, error) {
	span := util.Tracer.StartSpanFromContext(ctx, "CreateCampaign-service")
	defer util.Tracer.FinishSpan(span)
	nextCtx := util.Tracer.ContextWithSpan(ctx, span)

	err := makeCampaign(nextCtx, campaignRequest.PostID, userId)
	if err != nil {
		util.Tracer.LogError(span, err)
		return model.Campaign{}, err
	}
	campaignParams := model.CampaignParameters{
		Model:            gorm.Model{},
		Start:            campaignRequest.Start,
		End:              campaignRequest.End,
		CampaignID:       0,
		Interests:        service.getInterestsFromRequest(nextCtx, campaignRequest.Interests),
		CampaignRequests: getCampaignRequestsForProfileId(nextCtx, campaignRequest.InfluencerProfileIds),
		Timestamps:       getTimestampsFromRequest(nextCtx, campaignRequest.Timestamps),
	}

	campaign := model.Campaign{
		Model:              gorm.Model{},
		PostID:             campaignRequest.PostID,
		AgentID:            userId,
		CampaignType:       getCampaignTypeFromRequest(campaignRequest.Start, campaignRequest.End, len(campaignRequest.Timestamps)),
		Start:              campaignRequest.Start,
		CampaignParameters: []model.CampaignParameters{campaignParams},
	}
	return service.CampaignRepository.CreateCampaign(nextCtx, campaign)
}

func getTimestampsFromRequest(ctx context.Context, timestamps []time.Time) []model.Timestamp {
	span := util.Tracer.StartSpanFromContext(ctx, "getTimestampsFromRequest-service")
	defer util.Tracer.FinishSpan(span)

	util.Tracer.LogFields(span, "service", fmt.Sprintf("servicing convert times to timestamps model \n",))

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

func getCampaignRequestsForProfileId(ctx context.Context, profileIds []string) []model.CampaignRequest {
	span := util.Tracer.StartSpanFromContext(ctx, "getCampaignRequestsForProfileId-service")
	defer util.Tracer.FinishSpan(span)

	util.Tracer.LogFields(span, "service", fmt.Sprintf("servicing get campaign request for profile ids \n",))

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

func (service *CampaignService) getInterestsFromRequest(ctx context.Context, interests []string) []model.Interest {
	span := util.Tracer.StartSpanFromContext(ctx, "getInterestsFromRequest-service")
	defer util.Tracer.FinishSpan(span)
	nextCtx := util.Tracer.ContextWithSpan(ctx, span)

	return service.CampaignRepository.GetInterests(nextCtx, interests)
}

func (service *CampaignService) UpdateCampaignParameters(ctx context.Context, id uint, params dto.CampaignParametersDTO) error {
	span := util.Tracer.StartSpanFromContext(ctx, "UpdateCampaignParameters-service")
	defer util.Tracer.FinishSpan(span)
	util.Tracer.LogFields(span, "service", fmt.Sprintf("servicing campaign id %v\n", id))
	nextCtx := util.Tracer.ContextWithSpan(ctx, span)

	newParams := model.CampaignParameters{
		Model:            gorm.Model{},
		Start:            time.Time{},
		End:              params.End,
		CampaignID:       id,
		Interests:        service.getInterestsFromRequest(nextCtx, params.Interests),
		CampaignRequests: getCampaignRequestsForProfileId(nextCtx, params.InfluencerProfileIds),
		Timestamps:       getTimestampsFromRequest(nextCtx, params.Timestamps),
	}
	return service.CampaignRepository.UpdateCampaignParameters(nextCtx, newParams)
}

func (service *CampaignService) DeleteCampaign(ctx context.Context, id uint) error {
	span := util.Tracer.StartSpanFromContext(ctx, "DeleteCampaign-service")
	defer util.Tracer.FinishSpan(span)
	util.Tracer.LogFields(span, "service", fmt.Sprintf("servicing campaign id %v\n", id))
	nextCtx := util.Tracer.ContextWithSpan(ctx, span)

	return service.CampaignRepository.DeleteCampaign(nextCtx, id)
}

func (service *CampaignService) GetMyCampaigns(ctx context.Context, agentID uint) ([]dto.CampaignWithPostDTO, error) {
	span := util.Tracer.StartSpanFromContext(ctx, "GetMyCampaigns-service")
	defer util.Tracer.FinishSpan(span)
	util.Tracer.LogFields(span, "service", fmt.Sprintf("servicing agent id %v\n", agentID))
	nextCtx := util.Tracer.ContextWithSpan(ctx, span)

	campaigns, err := service.CampaignRepository.GetMyCampaigns(nextCtx, agentID)
	if err != nil {
		util.Tracer.LogError(span, err)
		return nil, err
	}
	postIDs := make([]string, 0)
	for _, value := range campaigns {
		postIDs = append(postIDs, value.PostID)
	}
	util.Tracer.LogFields(span, "service", fmt.Sprintf("servicing get posts for post ids\n"))
	posts, err := getPostsByPostsIds(nextCtx, postIDs)
	if err != nil {
		util.Tracer.LogError(span, err)
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

func (service *CampaignService) GetAllInterests(ctx context.Context) ([]string, error) {
	span := util.Tracer.StartSpanFromContext(ctx, "GetAllInterests-service")
	defer util.Tracer.FinishSpan(span)
	nextCtx := util.Tracer.ContextWithSpan(ctx, span)

	return service.CampaignRepository.GetAllInterests(nextCtx)
}

func getCampaignTypeFromRequest(start time.Time, end time.Time, timestampsLength int) model.CampaignType {
	if start.Equal(end) && timestampsLength == 1 {
		return model.ONE_TIME
	} else {
		return model.REPEATABLE
	}
}

func makeCampaign(ctx context.Context, postID string, loggedUserID uint) error {
	span := util.Tracer.StartSpanFromContext(ctx, "makeCampaign-service")
	defer util.Tracer.FinishSpan(span)
	util.Tracer.LogFields(span, "service", fmt.Sprintf("servicing post id %s for logged user id %v\n", postID, loggedUserID))
	nextCtx := util.Tracer.ContextWithSpan(ctx, span)

	postHost, postPort := util.GetPostHostAndPort()
	resp, err := util.CrossServiceRequest(nextCtx, http.MethodPost,
		util.GetCrossServiceProtocol()+"://"+postHost+":"+postPort+"/make-campaign/"+postID+"/"+util.Uint2String(loggedUserID),
		nil, map[string]string{})
	if err != nil {
		util.Tracer.LogError(span, err)
		return err
	}
	if resp.StatusCode != 200 {
		util.Tracer.LogError(span, fmt.Errorf("bad post id"))
		return fmt.Errorf("BAD_POST_ID")
	}
	return nil
}

func (service *CampaignService) GetCurrentlyValidInterests(ctx context.Context, campaignId uint) ([]string, error) {
	span := util.Tracer.StartSpanFromContext(ctx, "GetCurrentlyValidInterests-service")
	defer util.Tracer.FinishSpan(span)
	util.Tracer.LogFields(span, "service", fmt.Sprintf("servicing campaign id %v\n", campaignId))
	nextCtx := util.Tracer.ContextWithSpan(ctx, span)

	var ret []string
	parameters, err := service.CampaignRepository.GetParametersByCampaignId(nextCtx, campaignId)

	if err != nil{
		util.Tracer.LogError(span, err)
	}

	for _, i := range parameters.Interests {
		ret = append(ret, i.Name)
	}

	return ret, err
}

func (service *CampaignService) GetLastActiveParametersForCampaign(ctx context.Context, id uint) (model.CampaignParameters, error) {
	span := util.Tracer.StartSpanFromContext(ctx, "GetLastActiveParametersForCampaign-service")
	defer util.Tracer.FinishSpan(span)
	util.Tracer.LogFields(span, "service", fmt.Sprintf("servicing id %v\n", id))
	nextCtx := util.Tracer.ContextWithSpan(ctx, span)

	return service.CampaignRepository.GetLastActiveParametersForCampaign(nextCtx, id)
}

func (service *CampaignService) GetCampaignByIdForMonitoring(ctx context.Context, campaignId uint) (dto.CampaignMonitoringDTO, error) {
	span := util.Tracer.StartSpanFromContext(ctx, "GetCampaignByIdForMonitoring-service")
	defer util.Tracer.FinishSpan(span)
	util.Tracer.LogFields(span, "service", fmt.Sprintf("servicing campaign id %v\n", campaignId))
	nextCtx := util.Tracer.ContextWithSpan(ctx, span)

	var ret dto.CampaignMonitoringDTO
	var retParams []dto.CampaignParametersMonitoringDTO

	campaign, err := service.CampaignRepository.GetCampaignById(nextCtx, campaignId)
	if err != nil {
		util.Tracer.LogError(span, err)
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

func (service *CampaignService) GetAvailableCampaignsForUser(ctx context.Context, loggedUserID uint, followingProfiles []util.FollowingProfileDTO) ([]string, []uint, []uint, error) {
	span := util.Tracer.StartSpanFromContext(ctx, "GetAvailableCampaignsForUser-service")
	defer util.Tracer.FinishSpan(span)
	util.Tracer.LogFields(span, "service", fmt.Sprintf("servicing logged user id  %v\n", loggedUserID))
	nextCtx := util.Tracer.ContextWithSpan(ctx, span)

	util.Tracer.LogFields(span, "service", fmt.Sprintf("servicing get all interests for profile %v\n",loggedUserID))
	interests, err := getProfileInterests(nextCtx, loggedUserID)
	if err != nil {
		util.Tracer.LogError(span, err)
		return nil, nil, nil, err
	}
	allActiveParams, err := service.CampaignRepository.GetAllActiveParameters(nextCtx)
	if err != nil {
		util.Tracer.LogError(span, err)
		return nil, nil, nil, err
	}
	potentialInfluencers := make([]uint, 0)
	for _, followingProfile := range followingProfiles {
		potentialInfluencers = append(potentialInfluencers, followingProfile.ProfileID)
	}
	campaignIDs := make([]uint, 0)
	retInfluencerIDs := make([]uint, 0)
	var wg sync.WaitGroup
	for _, params := range allActiveParams {
		util.Tracer.LogFields(span, "service", fmt.Sprintf("servicing check if campaign with inserted params is active\n"))
		if campaignIsAvailableNow(params) {
			wg.Add(2)
			go func() {
				defer wg.Done()
				util.Tracer.LogFields(span, "service", fmt.Sprintf("servicing if inserted campaign paramters contains some of inserted interests\n"))
				if campaignParamsContainsInterest(params, interests) {
					campaignIDs = append(campaignIDs, params.CampaignID)
					retInfluencerIDs = append(retInfluencerIDs, 0)
				}
			}()
			go func() {
				defer wg.Done()
				util.Tracer.LogFields(span, "service", fmt.Sprintf("servicing search for active influencers in campaign parameters from list of inserted profile ids\n"))
				test, influencerIDs := campaignParamsContainsInfluencerIDs(params, potentialInfluencers)
				if test {
					for _, influencerID := range influencerIDs {
						campaignIDs = append(campaignIDs, params.CampaignID)
						retInfluencerIDs = append(retInfluencerIDs, influencerID)
					}
				}
			}()
			wg.Wait()
		}
	}
	if len(campaignIDs) == 0 {
		return make([]string, 0), make([]uint, 0), make([]uint, 0), nil
	}
	postIDs, err := service.CampaignRepository.GetPostIDsFromCampaignIDs(nextCtx, campaignIDs)
	return postIDs, retInfluencerIDs, campaignIDs, err
}

func (service *CampaignService) GetAcceptedCampaignsForInfluencer(ctx context.Context, influencerID uint) ([]string, []uint, []uint, error) {
	span := util.Tracer.StartSpanFromContext(ctx, "GetAcceptedCampaignsForInfluencer-service")
	defer util.Tracer.FinishSpan(span)
	util.Tracer.LogFields(span, "service", fmt.Sprintf("servicing influencer id %v\n", influencerID))
	nextCtx := util.Tracer.ContextWithSpan(ctx, span)

	allActiveParams, err := service.CampaignRepository.GetAllActiveParameters(nextCtx)
	if err != nil {
		util.Tracer.LogError(span, err)
		return nil, nil, nil, err
	}
	campaignIDs := make([]uint, 0)
	retInfluencerIDs := make([]uint, 0)
	for _, params := range allActiveParams {
		util.Tracer.LogFields(span, "service", fmt.Sprintf("servicing search for active influencers in campaign parameters from list of inserted profile ids "))
		test, _ := campaignParamsContainsInfluencerIDs(params, []uint{influencerID})
		if test {
			campaignIDs = append(campaignIDs, params.CampaignID)
			retInfluencerIDs = append(retInfluencerIDs, influencerID)
		}
	}
	if len(campaignIDs) == 0 {
		return make([]string, 0), make([]uint, 0), make([]uint, 0), nil
	}
	postIDs, err := service.CampaignRepository.GetPostIDsFromCampaignIDs(nextCtx, campaignIDs)
	return postIDs, retInfluencerIDs, campaignIDs, err
}

func (service *CampaignService) UpdateCampaignRequest(ctx context.Context, loggedUserId uint, requestId string, status model.RequestStatus) error {
	span := util.Tracer.StartSpanFromContext(ctx, "UpdateCampaignRequest-service")
	defer util.Tracer.FinishSpan(span)
	util.Tracer.LogFields(span, "service", fmt.Sprintf("servicing logged user id  %v, request id %s, status %v\n", loggedUserId, requestId, status))
	nextCtx := util.Tracer.ContextWithSpan(ctx, span)

	if service.CampaignRepository.GetCampaignRequestInfluencerId(nextCtx, util.String2Uint(requestId)) != loggedUserId {
		util.Tracer.LogError(span, fmt.Errorf("not allowed, logged user id and influencer id not match"))
		return fmt.Errorf("not allowed")
	}
	return service.CampaignRepository.UpdateCampaignRequest(nextCtx, requestId, status)
}

func (service *CampaignService) GetActiveCampaignsRequestsForProfileId(ctx context.Context, profileId uint) ([]dto.CampaignRequestForInfluencerDTO, error) {
	span := util.Tracer.StartSpanFromContext(ctx, "GetActiveCampaignsRequestsForProfileId-service")
	defer util.Tracer.FinishSpan(span)
	util.Tracer.LogFields(span, "service", fmt.Sprintf("servicing profile id %v\n", profileId))
	nextCtx := util.Tracer.ContextWithSpan(ctx, span)

	allActiveParams, err := service.CampaignRepository.GetAllActiveParameters(nextCtx)
	if err != nil {
		util.Tracer.LogError(span, err)
		return nil, err
	}
	ret := make([]dto.CampaignRequestForInfluencerDTO, 0)
	campaignIDs := make([]uint, 0)
	for _, params := range allActiveParams {
		util.Tracer.LogFields(span, "service", fmt.Sprintf("servicing campaign params that contains request for profile id  %v\n", profileId))
		test, requestID := campaignParamsContainsRequestForProfileID(params, profileId)
		if test {
			timestamps := make([]time.Time, 0)
			for _, timestamp := range params.Timestamps {
				timestamps = append(timestamps, timestamp.Time)
			}
			campaignIDs = append(campaignIDs, params.CampaignID)
			ret = append(ret, dto.CampaignRequestForInfluencerDTO{CampaignId: params.CampaignID, RequestId: requestID,
				Post: dto.PostDTO{}, Timestamps: timestamps, Start: params.Start, End: params.End})
		}
	}
	postIDs, err := service.CampaignRepository.GetPostIDsFromCampaignIDs(nextCtx, campaignIDs)
	if err != nil {
		util.Tracer.LogError(span, err)
		return nil, err
	}
	util.Tracer.LogFields(span, "service", fmt.Sprintf("servicing get posts for post ids\n"))
	posts, err := getPostsByPostsIds(nextCtx, postIDs)
	if err != nil {
		util.Tracer.LogError(span, err)
		return nil, err
	}
	if len(posts) != len(ret) {
		util.Tracer.LogError(span, fmt.Errorf("bad list sizes"))
		return nil, fmt.Errorf("BAD LIST SIZES")
	}
	for i, post := range posts {
		ret[i].Post = post
	}
	return ret, nil
}

func (service *CampaignService) GetCampaignRequestInfluencerId(ctx context.Context, requestId uint) uint {
	span := util.Tracer.StartSpanFromContext(ctx, "GetCampaignRequestInfluencerId-service")
	defer util.Tracer.FinishSpan(span)
	util.Tracer.LogFields(span, "service", fmt.Sprintf("servicing request id %v\n", requestId))
	nextCtx := util.Tracer.ContextWithSpan(ctx, span)

	return service.CampaignRepository.GetCampaignRequestInfluencerId(nextCtx, requestId)
}

func getPostsByPostsIds(ctx context.Context, postsIds []string) ([]dto.PostDTO, error) {
	span := util.Tracer.StartSpanFromContext(ctx, "getPostsByPostsIds-service")
	defer util.Tracer.FinishSpan(span)
	nextCtx := util.Tracer.ContextWithSpan(ctx, span)

	var ret []dto.PostDTO
	type data struct {
		Ids []string `json:"ids"`
	}
	bodyData := data{Ids: postsIds}
	jsonBody, err := json.Marshal(bodyData)
	if err != nil {
		util.Tracer.LogError(span, err)
		return nil, err
	}
	postHost, postPort := util.GetPostHostAndPort()
	resp, err := util.CrossServiceRequest(nextCtx, http.MethodPost,
		util.GetCrossServiceProtocol()+"://"+postHost+":"+postPort+"/posts",
		jsonBody, map[string]string{})

	if err != nil {
		util.Tracer.LogError(span, err)
		return nil, err
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		util.Tracer.LogError(span, err)
		return nil, err
	}

	defer func(Body io.ReadCloser) {
		_ = Body.Close()
	}(resp.Body)

	if err = json.Unmarshal(body, &ret); err != nil {
		util.Tracer.LogError(span, err)
		return nil, err
	}
	return ret, nil
}

func getProfileInterests(ctx context.Context, loggedUserID uint) ([]string, error) {
	span := util.Tracer.StartSpanFromContext(ctx, "getProfileInterests-service")
	defer util.Tracer.FinishSpan(span)
	util.Tracer.LogFields(span, "service", fmt.Sprintf("servicing logged user id %v\n", loggedUserID))
	nextCtx := util.Tracer.ContextWithSpan(ctx, span)

	profileHost, profilePort := util.GetProfileHostAndPort()
	resp, err := util.CrossServiceRequest(nextCtx, http.MethodGet,
		util.GetCrossServiceProtocol()+"://"+profileHost+":"+profilePort+"/profile-interests/"+util.Uint2String(loggedUserID),
		nil, map[string]string{})

	if err != nil {
		util.Tracer.LogError(span, err)
		return nil, err
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		util.Tracer.LogError(span, err)
		return nil, err
	}

	defer func(Body io.ReadCloser) {
		_ = Body.Close()
	}(resp.Body)
	var ret []string
	if err = json.Unmarshal(body, &ret); err != nil {
		util.Tracer.LogError(span, err)
		return nil, err
	}
	return ret, nil
}

func campaignIsAvailableNow(params model.CampaignParameters) bool {
	for _, timestamp := range params.Timestamps {
		timestampTime := timestamp.Time.UTC()
		now := time.Now().UTC()
		nowDate := time.Date(1, 1, 1, now.Hour(), now.Minute(), 0, 0, time.UTC)
		timestampTimeDate := time.Date(1, 1, 1, timestampTime.Hour(), timestampTime.Minute(), 0, 0, time.UTC)
		if (timestampTimeDate.Equal(now) || timestampTimeDate.Before(nowDate)) &&
			(timestampTimeDate.Add(1*time.Hour).Equal(now) || timestampTimeDate.Add(1*time.Hour).After(nowDate)) {
			return true
		}
	}
	return false
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

func campaignParamsContainsInfluencerIDs(params model.CampaignParameters, influencerIDs []uint) (bool, []uint) {
	influencers := make([]uint, 0)
	contains := false
	for _, campaignRequest := range params.CampaignRequests {
		for _, potentialInfluencer := range influencerIDs {
			if campaignRequest.RequestStatus == model.ACCEPTED && campaignRequest.InfluencerID == potentialInfluencer {
				contains = true
				influencers = append(influencers, campaignRequest.InfluencerID)
			}
		}
	}
	return contains, influencers
}

func campaignParamsContainsRequestForProfileID(params model.CampaignParameters, profileID uint) (bool, uint) {

	for _, campaignRequest := range params.CampaignRequests {
		if campaignRequest.InfluencerID == profileID && campaignRequest.RequestStatus == model.SENT {
			return true, campaignRequest.ID
		}
	}
	return false, 0
}
