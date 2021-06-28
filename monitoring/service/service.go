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
		EventType: model.GetEventType(eventDto.EventType), InfluencerUsername: eventDto.InfluencerUsername,
	Timestamp: time.Now(), CampaignId: eventDto.CampaignId, Interests: interests}

	err := service.MonitoringRepository.Create(event)
	return err
}

func (service *MonitoringService) CreateEventTargetGroup(eventDto dto.EventDTO) error{
	profilesInterests, err := getPersonalDataFromProfileMs(eventDto.ProfileId)
	if err != nil{
		return err
	}
	campaignInterests, err := getInterestsFromCampaignMs(eventDto.CampaignId)
	if err != nil{
		return err
	}
	var interests []string
	for _, profInt := range profilesInterests{
		for _, campInt := range campaignInterests{
			if profInt == campInt{
				interests = append(interests, profInt)
			}
		}
	}

	if err != nil{
		return err
	}

	event := &model.Event{OriginalPostId: eventDto.PostId,
		EventType: model.GetEventType(eventDto.EventType), InfluencerUsername: eventDto.InfluencerUsername,
		Timestamp: time.Now(), CampaignId: eventDto.CampaignId, Interests: interests}

	err = service.MonitoringRepository.Create(event)
	return err
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