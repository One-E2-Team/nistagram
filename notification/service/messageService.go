package service

import (
	"context"
	"errors"
	"fmt"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"nistagram/notification/model"
	"nistagram/util"
)

func (service *Service) CreateMessage(ctx context.Context, message *model.Message) error{
	span := util.Tracer.StartSpanFromContext(ctx, "CreateMessage-service")
	defer util.Tracer.FinishSpan(span)
	util.Tracer.LogFields(span, "service", fmt.Sprintf("servicing id %v %v\n", message.SenderId, message.ReceiverId))
	nextCtx := util.Tracer.ContextWithSpan(ctx, span)
	message.ID = primitive.NewObjectID()
	return service.Repository.CreateMessage(nextCtx, message)
}

func (service *Service) Seen(ctx context.Context, messageId string) error{
	span := util.Tracer.StartSpanFromContext(ctx, "Seen-service")
	defer util.Tracer.FinishSpan(span)
	util.Tracer.LogFields(span, "service", fmt.Sprintf("servicing id %v\n", messageId))
	nextCtx := util.Tracer.ContextWithSpan(ctx, span)
	return service.Repository.Seen(nextCtx, messageId)
}

func (service *Service) GetAllMessages(ctx context.Context, firstId uint, secondId uint) ([]model.Message,error){
	span := util.Tracer.StartSpanFromContext(ctx, "GetAllMessages-service")
	defer util.Tracer.FinishSpan(span)
	util.Tracer.LogFields(span, "service", fmt.Sprintf("servicing id %v %v\n", firstId, secondId))
	nextCtx := util.Tracer.ContextWithSpan(ctx, span)
	messages, err := service.Repository.GetAllMessages(nextCtx, firstId, secondId)
	if err != nil{
		util.Tracer.LogError(span, err)
		return nil,err
	}

	ejectSeenOneOf(&messages)

	messages = sortMessages(messages)

	return messages, nil
}

func (service *Service) DeleteMessage(ctx context.Context, loggedUserId uint, messageId string) error{
	span := util.Tracer.StartSpanFromContext(ctx, "DeleteMessage-service")
	defer util.Tracer.FinishSpan(span)
	util.Tracer.LogFields(span, "service", fmt.Sprintf("servicing id %v %v\n", loggedUserId, messageId))
	nextCtx := util.Tracer.ContextWithSpan(ctx, span)
	message, err := service.Repository.GetMessageById(nextCtx, messageId)
	if err != nil{
		util.Tracer.LogError(span, err)
		return err
	}

	if message.SenderId == loggedUserId || message.ReceiverId == loggedUserId{
		err = service.Repository.DeleteMessage(nextCtx, messageId)
	}else{
		util.Tracer.LogError(span, fmt.Errorf("This user is not allowed to delete message."))
		return errors.New("This user is not allowed to delete message.")
	}

	return err
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