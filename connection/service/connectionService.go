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

func (service *ConnectionService) GetConnection(followerId, profileId uint) *model.Connection {
	connection, _ := service.ConnectionRepository.SelectConnection(followerId, profileId, false)
	return connection
}

func (service *ConnectionService) FollowRequest(followerId, profileId uint) (*model.Connection, bool) {
	connection := service.ConnectionRepository.SelectOrCreateConnection(followerId, profileId)
	// TODO: private not private
	connection.ConnectionRequest = true
	resConnection, ok := service.ConnectionRepository.UpdateConnection(connection)
	if ok {
		return resConnection, true
	} else {
		return connection, false
	}
}

func (service *ConnectionService) Block(followerId, profileId uint) (*model.Connection, bool){
	connection := service.ConnectionRepository.SelectOrCreateConnection(followerId, profileId)
	connection.Block = true
	resConnection, ok := service.ConnectionRepository.UpdateConnection(connection)
	if ok {
		return resConnection, true
	} else {
		return connection, false
	}
}

func (service *ConnectionService) MessageConnect(followerId, profileId uint) (*model.Connection, bool){
	connection := service.ConnectionRepository.SelectOrCreateConnection(followerId, profileId)
	// TODO: obostrano
	connection.ConnectionRequest = false
	connection.MessageConnected = true
	resConnection, ok := service.ConnectionRepository.UpdateConnection(connection)
	if ok {
		return resConnection, true
	} else {
		return connection, false
	}
}

func (service *ConnectionService) MessageRequest(followerId, profileId uint) (*model.Connection, bool){
	connection := service.ConnectionRepository.SelectOrCreateConnection(followerId, profileId)
	connection.MessageRequest = true
	resConnection, ok := service.ConnectionRepository.UpdateConnection(connection)
	if ok {
		return resConnection, true
	} else {
		return connection, false
	}
}

func (service *ConnectionService) ApproveConnection(followerId, profileId uint) (*model.Connection, bool){
	connection, okSelect := service.ConnectionRepository.SelectConnection(followerId, profileId, false)
	if okSelect && connection == nil {
		return connection, false
	}
	connection.ConnectionRequest = false
	connection.Approved = true
	resConnection, ok := service.ConnectionRepository.UpdateConnection(connection)
	if ok {
		return resConnection, true
	} else {
		return connection, false
	}
}

func (service *ConnectionService) ToggleNotifyComment(followerId, profileId uint) (*model.Connection, bool){
	connection, okSelect := service.ConnectionRepository.SelectConnection(followerId, profileId, false)
	if okSelect && connection == nil {
		return connection, false
	}
	connection.NotifyComment = !connection.NotifyComment
	resConnection, ok := service.ConnectionRepository.UpdateConnection(connection)
	if ok {
		return resConnection, true
	} else {
		return connection, false
	}
}

func (service *ConnectionService) ToggleNotifyMessage(followerId, profileId uint) (*model.Connection, bool){
	connection, okSelect := service.ConnectionRepository.SelectConnection(followerId, profileId, false)
	if okSelect && connection == nil {
		return connection, false
	}
	connection.NotifyMessage = !connection.NotifyMessage
	resConnection, ok := service.ConnectionRepository.UpdateConnection(connection)
	if ok {
		return resConnection, true
	} else {
		return connection, false
	}
}

func (service *ConnectionService) ToggleNotifyStory(followerId, profileId uint) (*model.Connection, bool){
	connection, okSelect := service.ConnectionRepository.SelectConnection(followerId, profileId, false)
	if okSelect && connection == nil {
		return connection, false
	}
	connection.NotifyStory = !connection.NotifyStory
	resConnection, ok := service.ConnectionRepository.UpdateConnection(connection)
	if ok {
		return resConnection, true
	} else {
		return connection, false
	}
}

func (service *ConnectionService) ToggleNotifyPost(followerId, profileId uint) (*model.Connection, bool){
	connection, okSelect := service.ConnectionRepository.SelectConnection(followerId, profileId, false)
	if okSelect && connection == nil {
		return connection, false
	}
	connection.NotifyPost = !connection.NotifyPost
	resConnection, ok := service.ConnectionRepository.UpdateConnection(connection)
	if ok {
		return resConnection, true
	} else {
		return connection, false
	}
}

func (service *ConnectionService) ToggleCloseFriend(followerId, profileId uint) (*model.Connection, bool){
	connection, okSelect := service.ConnectionRepository.SelectConnection(followerId, profileId, false)
	if okSelect && connection == nil {
		return connection, false
	}
	connection.CloseFriend = !connection.CloseFriend
	resConnection, ok := service.ConnectionRepository.UpdateConnection(connection)
	if ok {
		return resConnection, true
	} else {
		return connection, false
	}
}

func (service *ConnectionService) ToggleMuted(followerId, profileId uint) (*model.Connection, bool){
	connection, okSelect := service.ConnectionRepository.SelectConnection(followerId, profileId, false)
	if okSelect && connection == nil {
		return connection, false
	}
	connection.Muted = !connection.Muted
	resConnection, ok := service.ConnectionRepository.UpdateConnection(connection)
	if ok {
		return resConnection, true
	} else {
		return connection, false
	}
}

func (service *ConnectionService) GetConnectedProfiles(conn model.Connection, excludeMuted bool) *[]model.Profile {
	ret := service.ConnectionRepository.GetConnectedProfiles(conn, excludeMuted)
	if ret == nil {
		temp := make([]model.Profile, 0)
		return &temp
	}
	return ret
}

