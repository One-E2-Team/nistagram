package service

import (
	"nistagram/monitoring/dto"
	"nistagram/monitoring/model"
	"nistagram/monitoring/repository"
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
	var interests []string
	event := &model.Event{OriginalPostId: eventDto.PostId,
		EventType: model.GetEventType(eventDto.EventType), InfluencerUsername: eventDto.InfluencerUsername,
		Timestamp: time.Now(), CampaignId: eventDto.CampaignId, Interests: interests}

	err := service.MonitoringRepository.Create(event)
	return err
}