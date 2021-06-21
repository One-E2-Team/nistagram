package service

import "nistagram/connection/model"

func (service *Service) IsInBlockingRelationship(id1, id2 uint) bool {
	lst := service.ConnectionRepository.GetBlockedProfiles(id1, false)
	if lst == nil || len(*lst) == 0 {
		return false
	}
	for _, val := range *lst {
		if val == id2 {
			return true
		}
	}
	return false
}

func (service *Service) IsBlocked(id1, id2 uint) bool {
	lst := service.ConnectionRepository.GetBlockedProfiles(id1, true)
	if lst == nil || len(*lst) == 0 {
		return false
	}
	for _, val := range *lst {
		if val == id2 {
			return true
		}
	}
	return false
}

func (service *Service) ToggleBlock(followerId, profileId uint) (*model.BlockEdge, bool) {
	block, ok := service.ConnectionRepository.SelectBlock(followerId, profileId)
	if !ok || block == nil {
		service.ConnectionRepository.DeleteConnection(followerId, profileId)
		service.ConnectionRepository.DeleteConnection(profileId, followerId)
		service.ConnectionRepository.DeleteMessage(followerId, profileId)
		service.ConnectionRepository.DeleteMessage(profileId, followerId)
		block, ok = service.ConnectionRepository.CreateBlock(followerId, profileId)
	} else {
		block, ok = service.ConnectionRepository.DeleteBlock(followerId, profileId)
	}
	return block, ok
}


func (service *Service) GetBlockingRelationships(u uint) *[]uint {
	return service.ConnectionRepository.GetBlockedProfiles(u, false)
}