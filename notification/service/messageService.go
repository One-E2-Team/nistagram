package service

import (
	"nistagram/notification/model"
)

func (service *Service) CreateMessage(message *model.Message) error{
	return service.Repository.CreateMessage(message)
}

func (service *Service) Seen(messageId string) error{
	return service.Repository.Seen(messageId)
}

func (service *Service) GetAllMessages(firstId uint, secondId uint) ([]model.Message,error){
	messages, err := service.Repository.GetAllMessages(firstId, secondId)
	if err != nil{
		return nil,err
	}

	ejectSeenOneOf(&messages)

	messages = sortMessages(messages)

	return messages, nil
}

func sortMessages(messages []model.Message) []model.Message {
	for i := 0; i < len(messages); i++{
		for j := i + 1; j < len(messages); j++{
			if messages[j].Timestamp.Before(messages[i].Timestamp){
				t := messages[i]
				messages[i] = messages[j]
				messages[j] = t
			}
		}
	}
	return messages
}

func ejectSeenOneOf(messages *[]model.Message){
	for _, m := range *messages{
		if m.OneOf && m.Seen{
			m.MediaPath = ""
		}
	}
}