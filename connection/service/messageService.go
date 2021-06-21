package service

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"nistagram/connection/dto"
	"nistagram/connection/model"
	model2 "nistagram/profile/model"
	"nistagram/util"
)

func (service *Service) MessageConnect(followerId, profileId uint) (*model.MessageEdge, bool) {
	message1, ok1 := service.ConnectionRepository.SelectMessage(followerId, profileId)
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
	messResp, ok2 := service.ConnectionRepository.CreateOrUpdateMessageRelationship(message2)
	if !ok2 || messResp == nil {
		return nil, false
	}
	message1, ok1 = service.ConnectionRepository.CreateOrUpdateMessageRelationship(*message1)
	if !ok1 || message1 == nil {
		return nil, false
	}
	return messResp, true
}

func (service *Service) MessageRequest(followerId, profileId uint) (*model.MessageEdge, bool) {
	if service.IsInBlockingRelationship(followerId, profileId) {
		return nil, false
	}
	message, messOk := service.ConnectionRepository.SelectMessage(profileId, followerId)
	if messOk != false {
		return nil, false
	}
	if message != nil {
		return service.MessageConnect(profileId, followerId)
	}
	message, messOk = service.ConnectionRepository.SelectMessage(followerId, profileId)
	if message != nil || messOk != false {
		return nil, false
	}
	connection, connOk := service.ConnectionRepository.SelectConnection(followerId, profileId, false)
	if connection != nil || connOk != false {
		return nil, false
	}
	newMessage := model.MessageEdge{
		PrimaryProfile:   followerId,
		SecondaryProfile: profileId,
		Approved:         false,
		NotifyMessage:    true,
	}
	resMessage, ok := service.ConnectionRepository.CreateOrUpdateMessageRelationship(newMessage)
	if ok {
		return resMessage, true
	} else {
		return nil, false
	}
}

func (service *Service) DeclineMessageRequest(followerId, profileId uint) (*model.MessageEdge, bool) {
	if service.IsInBlockingRelationship(followerId, profileId) {
		return nil, false
	}
	message, messOk := service.ConnectionRepository.SelectMessage(followerId, profileId)
	if message == nil || messOk == false {
		return nil, false
	}
	connection, connOk := service.ConnectionRepository.SelectConnection(followerId, profileId, false)
	if connection != nil || connOk != false {
		return nil, false
	}
	return service.ConnectionRepository.DeleteMessage(followerId, profileId)
}

func (service *Service) ToggleNotifyMessage(followerId, profileId uint) (*model.MessageEdge, bool) {
	if service.IsInBlockingRelationship(followerId, profileId) {
		return nil, false
	}
	message, messOk := service.ConnectionRepository.SelectMessage(followerId, profileId)
	if message == nil || messOk == false {
		return nil, false
	}
	message.NotifyMessage = !message.NotifyMessage
	resMessage, ok := service.ConnectionRepository.CreateOrUpdateMessageRelationship(*message)
	if ok {
		return resMessage, true
	} else {
		return nil, false
	}
}

func (service *Service) GetMessage(followerId, profileId uint) *model.MessageEdge {
	message, ok := service.ConnectionRepository.SelectMessage(followerId, profileId)
	if ok {
		return message
	} else {
		return nil
	}
}

func (service *Service) GetAllMessageRequests(id uint) *[]dto.UserDTO {
	var result = service.ConnectionRepository.GetAllMessageRequests(id)
	var ret = make([]dto.UserDTO, 0) // 0, :)
	for _, profileId := range *result {
		var p model2.Profile
		profileHost, profilePort := util.GetProfileHostAndPort()
		resp, err := util.CrossServiceRequest(http.MethodGet,
			util.CrossServiceProtocol+"://"+profileHost+":"+profilePort+"/get-by-id/"+util.Uint2String(profileId),
			nil, map[string]string{})
		if err != nil {
			fmt.Println(err)
			return nil
		}
		body, err1 := ioutil.ReadAll(resp.Body)
		if err1 != nil {
			fmt.Println(err1)
			return nil
		}
		err = json.Unmarshal(body, &p)
		if err != nil {
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