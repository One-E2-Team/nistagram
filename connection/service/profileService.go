package service

import "nistagram/connection/model"

func (service *Service) AddProfile(id uint) (*model.ProfileVertex, bool) {
	profile := model.ProfileVertex{ProfileID: id}
	ret := service.ConnectionRepository.CreateProfile(profile)
	return ret, ret.ProfileID == id
}