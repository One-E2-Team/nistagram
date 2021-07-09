package service

import (
	"context"
	"errors"
	"fmt"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"nistagram/notification/model"
	"nistagram/util"
	"sort"
	"time"
)

func (service *Service) CreateMessage(ctx context.Context, message *model.Message) error{
	span := util.Tracer.StartSpanFromContext(ctx, "CreateMessage-service")
	defer util.Tracer.FinishSpan(span)
	util.Tracer.LogFields(span, "service", fmt.Sprintf("servicing id %v %v\n", message.SenderId, message.ReceiverId))
	nextCtx := util.Tracer.ContextWithSpan(ctx, span)
	message.ID = primitive.NewObjectID()
	message.MessageSeen = false
	message.Timestamp = time.Now()
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
	if len(messages) == 0{
		return messages, nil
	}

	ejectSeenOneOf(&messages)

	fmt.Println("m1: ", messages)
	sort.Slice(messages, func(i,j int) bool{
		return messages[i].Timestamp.Before(messages[j].Timestamp)
	})
	fmt.Println("m2: ", messages)

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

func (service *Service) GetNotifications(ctx context.Context, receiverId uint) ([]string,error){
	span := util.Tracer.StartSpanFromContext(ctx, "GetNotifications-service")
	defer util.Tracer.FinishSpan(span)
	util.Tracer.LogFields(span, "service", fmt.Sprintf("servicing id  %v\n", receiverId))
	nextCtx := util.Tracer.ContextWithSpan(ctx, span)
	messages, err := service.Repository.GetUnseenMessages(nextCtx, receiverId)
	if err != nil{
		util.Tracer.LogError(span, err)
		return nil,err
	}

	fmt.Println("unseen messages: ", messages)

	senderIds := make([]uint, 0)
	for _, msg := range messages{
		if !util.Contains(senderIds, msg.SenderId){
			senderIds = append(senderIds, msg.SenderId)
		}
	}

	fmt.Println("Sender ids: ", senderIds)

	senderUsernames, err := getProfileUsernamesByIDs(nextCtx, senderIds)
	if err != nil{
		util.Tracer.LogError(span, err)
		return nil, err
	}

	fmt.Println("Sender usernames: ", senderUsernames)

	return senderUsernames, nil
}

func (service *Service) SeenMessage(ctx context.Context, receiverId uint,senderId uint) error{
	span := util.Tracer.StartSpanFromContext(ctx, "SeenMessage-service")
	defer util.Tracer.FinishSpan(span)
	util.Tracer.LogFields(span, "service", fmt.Sprintf("servicing id %v %v\n", receiverId, senderId))
	nextCtx := util.Tracer.ContextWithSpan(ctx, span)
	return service.Repository.SeenMessage(nextCtx, receiverId, senderId)
}

func ejectSeenOneOf(messages *[]model.Message){
	for _, m := range *messages{
		if m.OneOf && m.Seen{
			m.MediaPath = ""
		}
	}
}

