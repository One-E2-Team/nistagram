package service

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"nistagram/connection/dto"
	"nistagram/connection/model"
	model2 "nistagram/profile/model"
	"nistagram/util"
)

func (service *Service) MessageConnect(ctx context.Context, followerId, profileId uint) (*model.MessageEdge, bool) {
	span := util.Tracer.StartSpanFromContext(ctx, "MessageConnect-service")
	defer util.Tracer.FinishSpan(span)
	util.Tracer.LogFields(span, "service", fmt.Sprintf("servicing id %v %v\n", followerId, profileId))
	nextCtx := util.Tracer.ContextWithSpan(ctx, span)
	message1, ok1 := service.ConnectionRepository.SelectMessage(nextCtx, followerId, profileId)
	if !ok1 || message1 == nil || message1.Approved {
		return nil, false
	}
	message1.Approved = true
	message2 := model.MessageEdge{
		PrimaryProfile:   profileId,
		SecondaryProfile: followerId,
		Approved:         true,
		NotifyMessage:    true,
	}
	messResp, ok2 := service.ConnectionRepository.CreateOrUpdateMessageRelationship(nextCtx, message2)
	if !ok2 || messResp == nil {
		return nil, false
	}
	message1, ok1 = service.ConnectionRepository.CreateOrUpdateMessageRelationship(nextCtx, *message1)
	if !ok1 || message1 == nil {
		return nil, false
	}
	return messResp, true
}

func (service *Service) MessageRequest(ctx context.Context, followerId, profileId uint) (*model.MessageEdge, bool) {
	span := util.Tracer.StartSpanFromContext(ctx, "MessageRequest-service")
	defer util.Tracer.FinishSpan(span)
	util.Tracer.LogFields(span, "service", fmt.Sprintf("servicing id %v %v\n", followerId, profileId))
	nextCtx := util.Tracer.ContextWithSpan(ctx, span)
	if service.IsInBlockingRelationship(nextCtx, followerId, profileId) {
		return nil, false
	}
	message, messOk := service.ConnectionRepository.SelectMessage(nextCtx, profileId, followerId)
	if messOk != false {
		return nil, false
	}
	if message != nil {
		return service.MessageConnect(nextCtx, profileId, followerId)
	}
	message, messOk = service.ConnectionRepository.SelectMessage(nextCtx, followerId, profileId)
	if message != nil || messOk != false {
		return nil, false
	}
	connection, connOk := service.ConnectionRepository.SelectConnection(nextCtx, followerId, profileId, false)
	if connection != nil || connOk != false {
		return nil, false
	}
	newMessage := model.MessageEdge{
		PrimaryProfile:   followerId,
		SecondaryProfile: profileId,
		Approved:         false,
		NotifyMessage:    true,
	}
	resMessage, ok := service.ConnectionRepository.CreateOrUpdateMessageRelationship(nextCtx, newMessage)
	if ok {
		return resMessage, true
	} else {
		return nil, false
	}
}

func (service *Service) DeclineMessageRequest(ctx context.Context, followerId, profileId uint) (*model.MessageEdge, bool) {
	span := util.Tracer.StartSpanFromContext(ctx, "DeclineMessageRequest-service")
	defer util.Tracer.FinishSpan(span)
	util.Tracer.LogFields(span, "service", fmt.Sprintf("servicing id %v %v\n", followerId, profileId))
	nextCtx := util.Tracer.ContextWithSpan(ctx, span)
	if service.IsInBlockingRelationship(nextCtx, followerId, profileId) {
		return nil, false
	}
	message, messOk := service.ConnectionRepository.SelectMessage(nextCtx, followerId, profileId)
	if message == nil || messOk == false {
		return nil, false
	}
	connection, connOk := service.ConnectionRepository.SelectConnection(nextCtx, followerId, profileId, false)
	if connection != nil || connOk != false {
		return nil, false
	}
	return service.ConnectionRepository.DeleteMessage(nextCtx, followerId, profileId)
}

func (service *Service) ToggleNotifyMessage(ctx context.Context, followerId, profileId uint) (*model.MessageEdge, bool) {
	span := util.Tracer.StartSpanFromContext(ctx, "ToggleNotifyMessage-service")
	defer util.Tracer.FinishSpan(span)
	util.Tracer.LogFields(span, "service", fmt.Sprintf("servicing id %v %v\n", followerId, profileId))
	nextCtx := util.Tracer.ContextWithSpan(ctx, span)
	if service.IsInBlockingRelationship(nextCtx, followerId, profileId) {
		return nil, false
	}
	message, messOk := service.ConnectionRepository.SelectMessage(nextCtx, followerId, profileId)
	if message == nil || messOk == false {
		return nil, false
	}
	message.NotifyMessage = !message.NotifyMessage
	resMessage, ok := service.ConnectionRepository.CreateOrUpdateMessageRelationship(nextCtx, *message)
	if ok {
		return resMessage, true
	} else {
		return nil, false
	}
}

func (service *Service) GetMessage(ctx context.Context, followerId, profileId uint) *model.MessageEdge {
	span := util.Tracer.StartSpanFromContext(ctx, "GetMessage-service")
	defer util.Tracer.FinishSpan(span)
	util.Tracer.LogFields(span, "service", fmt.Sprintf("servicing id %v %v\n", followerId, profileId))
	nextCtx := util.Tracer.ContextWithSpan(ctx, span)
	message, ok := service.ConnectionRepository.SelectMessage(nextCtx, followerId, profileId)
	if ok {
		return message
	} else {
		return nil
	}
}

func (service *Service) GetAllMessageRequests(ctx context.Context, id uint) *[]dto.UserDTO {
	span := util.Tracer.StartSpanFromContext(ctx, "GetAllMessageRequests-service")
	defer util.Tracer.FinishSpan(span)
	util.Tracer.LogFields(span, "service", fmt.Sprintf("servicing id %v \n", id))
	nextCtx := util.Tracer.ContextWithSpan(ctx, span)
	var result = service.ConnectionRepository.GetAllMessageRequests(nextCtx, id)
	var ret = make([]dto.UserDTO, 0) // 0, :)
	for _, profileId := range *result {
		var p model2.Profile
		profileHost, profilePort := util.GetProfileHostAndPort()
		resp, err := util.CrossServiceRequest(nextCtx, http.MethodGet,
			util.GetCrossServiceProtocol()+"://"+profileHost+":"+profilePort+"/get-by-id/"+util.Uint2String(profileId),
			nil, map[string]string{})
		if err != nil {
			util.Tracer.LogError(span, err)
			fmt.Println(err)
			return nil
		}
		body, err1 := ioutil.ReadAll(resp.Body)
		if err1 != nil {
			util.Tracer.LogError(span, err)
			fmt.Println(err1)
			return nil
		}
		err = json.Unmarshal(body, &p)
		if err != nil {
			util.Tracer.LogError(span, err)
			fmt.Println(err)
			return nil
		}
		ret = append(ret, dto.UserDTO{
			Username:  p.Username,
			ProfileID: p.ID,
		})
		resp.Body.Close()
	}
	return &ret
}
