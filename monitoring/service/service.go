package service

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"nistagram/monitoring/dto"
	"nistagram/monitoring/model"
	"nistagram/monitoring/repository"
	"nistagram/util"
	"time"
)

type MonitoringService struct {
	MonitoringRepository *repository.MonitoringRepository
}

func (service *MonitoringService) CreateEventInfluencer(ctx context.Context, eventDto dto.EventDTO) error {
	span := util.Tracer.StartSpanFromContext(ctx, "CreateEventInfluencer-service")
	defer util.Tracer.FinishSpan(span)
	util.Tracer.LogFields(span, "service", fmt.Sprintf("servicing creation od event influencer"))
	nextCtx := util.Tracer.ContextWithSpan(ctx, span)
	var interests []string
	event := &model.Event{OriginalPostId: eventDto.PostId,
		EventType: model.GetEventType(eventDto.EventType), WebSite: eventDto.WebSite,
		InfluencerId: eventDto.InfluencerId,
		Timestamp:    time.Now(), CampaignId: eventDto.CampaignId, Interests: interests}

	err := service.MonitoringRepository.Create(nextCtx, event)
	return err
}

func (service *MonitoringService) CreateEventTargetGroup(ctx context.Context, eventDto dto.EventDTO) error {
	span := util.Tracer.StartSpanFromContext(ctx, "CreateEventTargetGroup-service")
	defer util.Tracer.FinishSpan(span)
	util.Tracer.LogFields(span, "service", fmt.Sprintf("servicing creation of event target group"))
	nextCtx := util.Tracer.ContextWithSpan(ctx, span)
	var interests []string
	if eventDto.ProfileId != 0 {
		profilesInterests, err := getPersonalDataFromProfileMs(nextCtx, eventDto.ProfileId)
		if err != nil {
			util.Tracer.LogError(span, err)
			return err
		}
		campaignInterests, err := getInterestsFromCampaignMs(nextCtx, eventDto.CampaignId)
		if err != nil {
			util.Tracer.LogError(span, err)
			return err
		}
		for _, profInt := range profilesInterests {
			for _, campInt := range campaignInterests {
				if profInt == campInt {
					interests = append(interests, profInt)
				}
			}
		}
	}

	event := &model.Event{OriginalPostId: eventDto.PostId,
		EventType: model.GetEventType(eventDto.EventType), WebSite: eventDto.WebSite, InfluencerId: eventDto.InfluencerId,
		Timestamp: time.Now(), CampaignId: eventDto.CampaignId, Interests: interests}

	err := service.MonitoringRepository.Create(nextCtx, event)
	return err
}


func (service *MonitoringService) VisitSite(ctx context.Context, campaignId uint, influencerId uint,loggedUserId uint, mediaId string) (string, error){
	span := util.Tracer.StartSpanFromContext(ctx, "VisitSite-service")
	defer util.Tracer.FinishSpan(span)
	util.Tracer.LogFields(span, "service", fmt.Sprintf("servicing id %v\n", mediaId))
	nextCtx := util.Tracer.ContextWithSpan(ctx, span)

	webSite, err := getMediaByIdFromPostMs(nextCtx, mediaId)
	if err != nil {
		util.Tracer.LogError(span, err)
		return "", err
	}

	eventDto := &dto.EventDTO{EventType: "visit", PostId: "", WebSite: webSite,
		ProfileId: loggedUserId, CampaignId: campaignId, InfluencerId: influencerId}

	if influencerId != 0 {
		err := service.CreateEventInfluencer(nextCtx, *eventDto)
		if err != nil {
			util.Tracer.LogError(span, err)
			return webSite, err
		}
	} else {
		err := service.CreateEventTargetGroup(nextCtx, *eventDto)
		if err != nil {
			util.Tracer.LogError(span, err)
			return webSite, err
		}
	}

	return webSite, err
}

func (service *MonitoringService) GetCampaignStatistics(ctx context.Context, campaignId uint) (dto.StatisticsDTO, error) {
	span := util.Tracer.StartSpanFromContext(ctx, "GetCampaignStatistics-service")
	defer util.Tracer.FinishSpan(span)
	util.Tracer.LogFields(span, "service", fmt.Sprintf("servicing id %v\n", campaignId))
	nextCtx := util.Tracer.ContextWithSpan(ctx, span)
	var statistics dto.StatisticsDTO
	campaign, err := getCampaignByCampaignId(nextCtx, campaignId)
	if err != nil {
		util.Tracer.LogError(span, err)
		return statistics, err
	}

	events, err := service.MonitoringRepository.GetEventsByCampaignId(nextCtx, campaignId)
	if err != nil {
		util.Tracer.LogError(span, err)
		return statistics, err
	}
	var infIds []uint
	for _, param := range campaign.CampaignParameters {
		for _, req := range param.CampaignRequests {
			if req.InfluencerID != 0 && !util.Contains(infIds, req.InfluencerID) {
				infIds = append(infIds, req.InfluencerID)
			}
		}
	}

	for _, e := range events {
		if e.InfluencerId != 0 && !util.Contains(infIds, e.InfluencerId) {
			infIds = append(infIds, e.InfluencerId)
		}
	}

	usernames, err := getProfileUsernamesByIDs(nextCtx, infIds)
	usersMap := make(map[uint]string, 0)
	for i, id := range infIds {
		usersMap[id] = usernames[i]
	}

	for i, param := range campaign.CampaignParameters {
		for j, req := range param.CampaignRequests {
			campaign.CampaignParameters[i].CampaignRequests[j].InfluencerUsername = usersMap[req.InfluencerID]
		}
	}

	var eventDtos []dto.ShowEventDTO
	for _, e := range events {
		eventDto := &dto.ShowEventDTO{EventType: e.EventType.ToString(), InfluencerId: e.InfluencerId,
			InfluencerUsername: usersMap[e.InfluencerId],
			Interests:          e.Interests, WebSite: e.WebSite, Timestamp: e.Timestamp}
		eventDtos = append(eventDtos, *eventDto)
	}

	statistics.CampaignId = campaignId
	statistics.Campaign = campaign
	statistics.Events = eventDtos
	return statistics, err
}

func getPersonalDataFromProfileMs(ctx context.Context, profileId uint) ([]string, error) {
	span := util.Tracer.StartSpanFromContext(ctx, "getPersonalDataFromProfileMs-service")
	defer util.Tracer.FinishSpan(span)
	util.Tracer.LogFields(span, "service", fmt.Sprintf("servicing id %v\n", profileId))
	nextCtx := util.Tracer.ContextWithSpan(ctx, span)
	profileHost, profilePort := util.GetProfileHostAndPort()
	resp, err := util.CrossServiceRequest(nextCtx, http.MethodGet,
		util.GetCrossServiceProtocol()+"://"+profileHost+":"+profilePort+"/personal-data/"+util.Uint2String(profileId),
		nil, map[string]string{})

	var ret []string
	personalDataDto := &dto.PersonalDataDTO{}
	if err != nil {
		util.Tracer.LogError(span, err)
		return ret, err
	}

	err = json.NewDecoder(resp.Body).Decode(&personalDataDto)

	if err != nil {
		util.Tracer.LogError(span, err)
		return ret, err
	}

	for _, interest := range personalDataDto.InterestedIn {
		ret = append(ret, interest.Name)
	}

	return ret, nil
}

func getInterestsFromCampaignMs(ctx context.Context, campaignId uint) ([]string, error) {
	span := util.Tracer.StartSpanFromContext(ctx, "getInterestsFromCampaignMs-service")
	defer util.Tracer.FinishSpan(span)
	util.Tracer.LogFields(span, "service", fmt.Sprintf("servicing id %v\n", campaignId))
	nextCtx := util.Tracer.ContextWithSpan(ctx, span)
	campHost, campPort := util.GetCampaignHostAndPort()
	resp, err := util.CrossServiceRequest(nextCtx, http.MethodGet,
		util.GetCrossServiceProtocol()+"://"+campHost+":"+campPort+"/interests/"+util.Uint2String(campaignId),
		nil, map[string]string{})

	var ret []string
	if err != nil {
		util.Tracer.LogError(span, err)
		return ret, err
	}

	err = json.NewDecoder(resp.Body).Decode(&ret)

	return ret, err
}

func getMediaByIdFromPostMs(ctx context.Context, mediaId string) (string, error) {
	span := util.Tracer.StartSpanFromContext(ctx, "getMediaByIdFromPostMs-service")
	defer util.Tracer.FinishSpan(span)
	util.Tracer.LogFields(span, "service", fmt.Sprintf("servicing id %v\n", mediaId))
	nextCtx := util.Tracer.ContextWithSpan(ctx, span)
	postHost, postPort := util.GetPostHostAndPort()
	resp, err := util.CrossServiceRequest(nextCtx, http.MethodGet,
		util.GetCrossServiceProtocol()+"://"+postHost+":"+postPort+"/media/"+mediaId,
		nil, map[string]string{})

	var dto dto.MediaDTO
	if err != nil {
		util.Tracer.LogError(span, err)
		return "", err
	}

	err = json.NewDecoder(resp.Body).Decode(&dto)

	return dto.WebSite, err
}

func getCampaignByCampaignId(ctx context.Context, campaignId uint) (dto.CampaignMonitoringDTO, error) {
	span := util.Tracer.StartSpanFromContext(ctx, "getCampaignByCampaignId-service")
	defer util.Tracer.FinishSpan(span)
	util.Tracer.LogFields(span, "service", fmt.Sprintf("servicing id %v\n", campaignId))
	nextCtx := util.Tracer.ContextWithSpan(ctx, span)
	campHost, campPort := util.GetCampaignHostAndPort()
	resp, err := util.CrossServiceRequest(nextCtx, http.MethodGet,
		util.GetCrossServiceProtocol()+"://"+campHost+":"+campPort+"/campaign/monitoring/"+util.Uint2String(campaignId),
		nil, map[string]string{})

	var ret dto.CampaignMonitoringDTO
	if err != nil {
		util.Tracer.LogError(span, err)
		return ret, err
	}

	err = json.NewDecoder(resp.Body).Decode(&ret)
	return ret, err
}

func getProfileUsernamesByIDs(ctx context.Context, profileIDs []uint) ([]string, error) {
	span := util.Tracer.StartSpanFromContext(ctx, "getProfileUsernamesByIDs-service")
	defer util.Tracer.FinishSpan(span)
	util.Tracer.LogFields(span, "service", fmt.Sprintf("servicing id %v\n", profileIDs))
	nextCtx := util.Tracer.ContextWithSpan(ctx, span)
	type data struct {
		Ids []string `json:"ids"`
	}
	req := make([]string, 0)
	for _, value := range profileIDs {
		req = append(req, util.Uint2String(value))
	}
	bodyData := data{Ids: req}
	jsonBody, err := json.Marshal(bodyData)
	if err != nil {
		util.Tracer.LogError(span, err)
		return nil, err
	}
	profileHost, profilePort := util.GetProfileHostAndPort()
	resp, err := util.CrossServiceRequest(nextCtx, http.MethodPost,
		util.GetCrossServiceProtocol()+"://"+profileHost+":"+profilePort+"/get-by-ids",
		jsonBody, map[string]string{"Content-Type": "application/json;"})
	if err != nil {
		util.Tracer.LogError(span, err)
		return nil, err
	}

	var ret []string

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		util.Tracer.LogError(span, err)
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
