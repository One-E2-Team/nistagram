package service

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"nistagram/connection/model"
	"nistagram/connection/repository"
	model2 "nistagram/profile/model"
	"nistagram/util"
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

func getProfile(id uint) *model2.Profile{
	var p model2.Profile
	profileHost, profilePort := util.GetProfileHostAndPort()
	resp, err := http.Get("http://" + profileHost+":" + profilePort + "/get-by-id/" + util.Uint2String(id))
	if err != nil {
		fmt.Println(err)
		return nil
	}
	defer resp.Body.Close()
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
	return &p
}

func (service *ConnectionService) FollowRequest(followerId, profileId uint) (*model.Connection, bool) {
	connection := service.ConnectionRepository.SelectOrCreateConnection(followerId, profileId)
	conn2, ok2 := service.ConnectionRepository.SelectConnection(followerId, profileId, false)
	profile1 := getProfile(followerId)
	profile2 := getProfile(profileId)
	if !ok2 || profile1 == nil || profile2 == nil {
		return nil, false
	}
	if connection.Block == true || (conn2 != nil && conn2.Block == true) {
		return nil, false
	}
	if profile1.ProfileSettings.IsPrivate == false && profile2.ProfileSettings.IsPrivate == false {
		connection.Approved = true
	} else {
		connection.ConnectionRequest = true
	}
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
	connection, ok := service.ConnectionRepository.SelectConnection(followerId, profileId, false)
	conn2, ok2 := service.ConnectionRepository.SelectConnection(profileId, followerId, false)
	if !connection.MessageRequest || (!ok || !ok2){
		return nil, false
	}
	connection.MessageRequest = false
	connection.MessageConnected = true
	conn2.MessageRequest = false
	conn2.MessageConnected = true
	service.ConnectionRepository.UpdateConnection(connection)
	resConnection, ok1 := service.ConnectionRepository.UpdateConnection(conn2)
	if ok1 {
		return resConnection, true
	} else {
		return conn2, false
	}
}

func (service *ConnectionService) MessageRequest(followerId, profileId uint) (*model.Connection, bool){
	connection := service.ConnectionRepository.SelectOrCreateConnection(followerId, profileId)
	if connection.MessageConnected {
		return nil, false
	}
	connection.MessageRequest = true
	conn2 := service.ConnectionRepository.SelectOrCreateConnection(profileId, followerId)
	if conn2.MessageConnected {
		return nil, false
	}
	if !conn2.Approved {
		conn2.MessageRequest = false
	}
	resConnection, ok := service.ConnectionRepository.UpdateConnection(connection)
	service.ConnectionRepository.UpdateConnection(conn2)
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
	conn2, ok2 := service.ConnectionRepository.SelectConnection(profileId, followerId, false)
	profile1 := getProfile(followerId)
	profile2 := getProfile(profileId)
	if !ok2 || profile1 == nil || profile2 == nil {
		return nil, false
	}
	if connection.Block == true || (conn2 != nil && conn2.Block == true) {
		return nil, false
	}
	if conn2 == nil {
		conn2 = service.ConnectionRepository.SelectOrCreateConnection(profileId, followerId)
	}
	if !connection.ConnectionRequest {
		return nil, false
	}
	connection.ConnectionRequest = false
	connection.Approved = true
	conn2.ConnectionRequest = false
	conn2.Approved = true
	service.ConnectionRepository.UpdateConnection(connection)
	var ok bool
	conn2, ok = service.ConnectionRepository.UpdateConnection(conn2)
	if ok {
		return conn2, true
	} else {
		return conn2, false
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

