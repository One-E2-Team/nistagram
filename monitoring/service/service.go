package service

import (
	"encoding/json"
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

func (service *MonitoringService) CreateEventInfluencer(eventDto dto.EventDTO) error{
	var interests []string
	event := &model.Event{OriginalPostId: eventDto.PostId,
		EventType: model.GetEventType(eventDto.EventType), WebSite: eventDto.WebSite,
		InfluencerId: eventDto.InfluencerId,
	Timestamp: time.Now(), CampaignId: eventDto.CampaignId, Interests: interests}

	err := service.MonitoringRepository.Create(event)
	return err
}

func (service *MonitoringService) CreateEventTargetGroup(eventDto dto.EventDTO) error{
	var interests []string
	if eventDto.ProfileId != 0{
		profilesInterests, err := getPersonalDataFromProfileMs(eventDto.ProfileId)
		if err != nil{
			return err
		}
		campaignInterests, err := getInterestsFromCampaignMs(eventDto.CampaignId)
		if err != nil{
			return err
		}
		for _, profInt := range profilesInterests{
			for _, campInt := range campaignInterests{
				if profInt == campInt{
					interests = append(interests, profInt)
				}
			}
		}
	}

	event := &model.Event{OriginalPostId: eventDto.PostId,
		EventType: model.GetEventType(eventDto.EventType), WebSite: eventDto.WebSite, InfluencerId: eventDto.InfluencerId,
		Timestamp: time.Now(), CampaignId: eventDto.CampaignId, Interests: interests}

	err := service.MonitoringRepository.Create(event)
	return err
}

func (service *MonitoringService) VisitSite(campaignId uint, influencerId uint,loggedUserId uint, mediaId uint) (string, error){

	webSite, err := getMediaByIdFromPostMs(mediaId)
	if err != nil {
		return "", err
	}

	eventDto := &dto.EventDTO{EventType: "visit", PostId: "", WebSite: webSite,
		ProfileId: loggedUserId, CampaignId: campaignId, InfluencerId: influencerId}

	if influencerId != 0{
		err := service.CreateEventInfluencer(*eventDto)
		if err != nil{
			return webSite, err
		}
	}else{
		err := service.CreateEventTargetGroup(*eventDto)
		if err != nil{
			return webSite, err
		}
	}

	return webSite, err
}

func getPersonalDataFromProfileMs(profileId uint) ([]string, error){
	profileHost, profilePort := util.GetProfileHostAndPort()
	resp, err := util.CrossServiceRequest(http.MethodGet,
		util.GetCrossServiceProtocol()+"://"+profileHost+":"+profilePort+"/personal-data/"+ util.Uint2String(profileId),
		nil, map[string]string{})

	var ret []string
	personalDataDto := &dto.PersonalDataDTO{}
	if err != nil{
		return ret, err
	}

	err = json.NewDecoder(resp.Body).Decode(&personalDataDto)

	if err != nil{
		return ret,err
	}

	for _, interest := range personalDataDto.InterestedIn{
		ret = append(ret, interest.Name)
	}

	return ret, nil
}

func getInterestsFromCampaignMs(campaignId uint) ([]string, error){
	campHost, campPort := util.GetCampaignHostAndPort()
	resp, err := util.CrossServiceRequest(http.MethodGet,
		util.GetCrossServiceProtocol()+"://"+campHost+":"+campPort+"/interests/"+ util.Uint2String(campaignId),
		nil, map[string]string{})

	var ret []string
	if err != nil{
		return ret, err
	}

	err = json.NewDecoder(resp.Body).Decode(&ret)

	return ret, err
}

func getMediaByIdFromPostMs(mediaId uint) (string, error){
	postHost, postPort := util.GetPostHostAndPort()
	resp, err := util.CrossServiceRequest(http.MethodGet,
		util.GetCrossServiceProtocol()+"://"+postHost+":"+postPort+"/media/"+ util.Uint2String(mediaId),
		nil, map[string]string{})

	var dto dto.MediaDTO
	if err != nil{
		return "", err
	}

	err = json.NewDecoder(resp.Body).Decode(&dto)

	return dto.WebSite, err
}