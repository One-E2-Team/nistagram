package service

import "nistagram/connection/model"

func (service *Service) AddOrUpdateProfile(profile model.ProfileVertex) (*model.ProfileVertex, bool) {
	ret := service.ConnectionRepository.CreateOrUpdateProfile(profile)
	return ret, ret != nil && ret.ProfileID == profile.ProfileID && ret.Deleted == profile.Deleted
}