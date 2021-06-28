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

func (service *MonitoringService) CreateEvent(eventDto dto.EventDTO) error{
	event := &model.Event{OriginalPostId: eventDto.PostId,
		EventType: model.GetEventType(eventDto.Type), InfluencerUsername: eventDto.InfluencerUsername,
	Timestamp: time.Now(), CampaignId: eventDto.CampaignId, WebSite: eventDto.WebSite}

	err := service.MonitoringRepository.Create(event)
	return err
}