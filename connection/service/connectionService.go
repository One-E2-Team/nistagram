package service

import (
	"nistagram/connection/model"
	"nistagram/connection/repository"
)

type ConnectionService struct {
	ConnectionRepository *repository.ConnectionRepository
}

func (service *ConnectionService) AddProfile(id uint) (*model.Profile, bool) {
	profile := model.Profile{ProfileID: id}
	ret := service.ConnectionRepository.CreateProfile(profile)
	return ret, ret.ProfileID == id
}

func (service *ConnectionService) FollowRequest(followerId, profileId uint) (*model.Connection, bool) {
	service.ConnectionRepository.SelectOrCreateConnection(followerId, profileId)
	return &model.Connection{}, false
}