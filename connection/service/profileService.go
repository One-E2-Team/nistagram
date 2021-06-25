package service

import (
	"nistagram/connection/dto"
	"nistagram/connection/model"
)

func (service *Service) AddOrUpdateProfile(profile model.ProfileVertex) (*model.ProfileVertex, bool) {
	ret := service.ConnectionRepository.CreateOrUpdateProfile(profile)
	return ret, ret != nil && ret.ProfileID == profile.ProfileID && ret.Deleted == profile.Deleted
}

func (service *Service) GetRecommendations(id uint) (*[]dto.UserDTO, bool) {
	return nil, false
}