package service

import (
	"nistagram/connection/model"
	"nistagram/connection/repository"
)

type ConnectionService struct {
	ConnectionRepository *repository.ConnectionRepository
}

func (service *ConnectionService) AddProfile(id uint, username string) (*model.Profile, bool) {
	profile := model.Profile{ProfileID: id, Username: username}
	ret := service.ConnectionRepository.CreateProfile(profile)
	return ret, ret.ProfileID == id
}