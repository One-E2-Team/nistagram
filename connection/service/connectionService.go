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

func (service *Service) GetConnection(followerId, profileId uint) *model.ConnectionEdge {
	connection, _ := service.ConnectionRepository.SelectConnection(followerId, profileId, false)
	return connection
}

func (service *Service) GetConnectedProfiles(conn model.ConnectionEdge, excludeMuted, excludeBlocked bool) *[]dto.ConnectedProfileDTO {
	profiles := service.ConnectionRepository.GetConnectedProfiles(conn, excludeMuted, true)
	if profiles == nil {
		temp := make([]uint, 0)
		profiles = &temp
	}
	if !excludeBlocked {
		var final []uint
		blocking := service.ConnectionRepository.GetBlockedProfiles(conn.PrimaryProfile, false)
		for _, val := range *profiles {
			if !contains(blocking, val) {
				final = append(final, val)
			}
		}
		profiles = &final
	}
	ret := make([]dto.ConnectedProfileDTO, 0)
	for _, val := range *profiles {
		var closeFriend bool
		invConnection, ok := service.ConnectionRepository.SelectConnection(val, conn.PrimaryProfile, false)
		if !ok || invConnection == nil {
			closeFriend = false
		} else {
			closeFriend = invConnection.CloseFriend
		}
		ret = append(ret, dto.ConnectedProfileDTO{
			ProfileID:   val,
			CloseFriend: closeFriend,
		})
	}
	return &ret
}

func (service *Service) GetProfilesInFollowRelationship(conn model.ConnectionEdge, excludeMuted, excludeBlocked bool, following bool) *[]dto.UserDTO {
	profiles := service.ConnectionRepository.GetConnectedProfiles(conn, excludeMuted, following)
	if profiles == nil {
		temp := make([]uint, 0)
		profiles = &temp
	}
	if !excludeBlocked {
		var final []uint
		blockId := conn.PrimaryProfile
		if !following {
			blockId = conn.SecondaryProfile
		}
		blocking := service.ConnectionRepository.GetBlockedProfiles(blockId, false)
		for _, val := range *profiles {
			if !contains(blocking, val) {
				final = append(final, val)
			}
		}
		profiles = &final
	}
	ret := make([]dto.UserDTO, 0)
	for _, val := range *profiles {
		p := util.GetProfile(context.Background(), val)
		if p == nil {
			continue
		}
		ret = append(ret, dto.UserDTO{
			ProfileID: val,
			Username:  p.Username,
		})
	}
	return &ret
}

func (service *Service) UpdateConnection(id uint, conn model.ConnectionEdge) (*model.ConnectionEdge, bool) {
	if id == conn.PrimaryProfile {
		return service.ConnectionRepository.UpdateConnection(&conn)
	} else {
		return nil, false
	}
}

func (service *Service) DeleteConnection(followerId, profileId uint) (*model.ConnectionEdge, bool) {
	return service.ConnectionRepository.DeleteConnection(followerId, profileId)
}

func (service *Service) FollowRequest(followerId, profileId uint) (*model.ConnectionEdge, bool) {
	if service.IsInBlockingRelationship(followerId, profileId) {
		return nil, false
	}
	connection := service.ConnectionRepository.SelectOrCreateConnection(followerId, profileId)
	if connection.Approved {
		return nil, false
	}
	profile2 := util.GetProfile(context.Background(), profileId)
	if profile2 == nil {
		return nil, false
	}
	if profile2.ProfileSettings.IsPrivate == false {
		connection.Approved = true
		service.ConnectionRepository.CreateOrUpdateMessageRelationship(model.MessageEdge{
			PrimaryProfile:   profileId,
			SecondaryProfile: followerId,
			Approved:         true,
			NotifyMessage:    true,
		})
	} else {
		connection.ConnectionRequest = true
	}
	service.ConnectionRepository.CreateOrUpdateMessageRelationship(model.MessageEdge{
		PrimaryProfile:   followerId,
		SecondaryProfile: profileId,
		Approved:         true,
		NotifyMessage:    true,
	})
	resConnection, ok := service.ConnectionRepository.UpdateConnection(connection)
	if ok {
		return resConnection, true
	} else {
		return connection, false
	}
}

func (service *Service) ApproveConnection(followerId, profileId uint) (*model.ConnectionEdge, bool) {
	connection, okSelect := service.ConnectionRepository.SelectConnection(followerId, profileId, false)
	if okSelect && connection == nil {
		return connection, false
	}
	profile1 := util.GetProfile(context.Background(), followerId)
	profile2 := util.GetProfile(context.Background(), profileId)
	if profile1 == nil || profile2 == nil {
		return nil, false
	}
	if !connection.ConnectionRequest {
		return nil, false
	}
	connection.ConnectionRequest = false
	connection.Approved = true
	service.ConnectionRepository.CreateOrUpdateMessageRelationship(model.MessageEdge{
		PrimaryProfile:   profileId,
		SecondaryProfile: followerId,
		Approved:         true,
		NotifyMessage:    true,
	})
	return service.ConnectionRepository.UpdateConnection(connection)
}

func (service *Service) Unfollow(followerId, profileId uint) (*model.ConnectionEdge, bool) {
	if service.IsInBlockingRelationship(followerId, profileId) {
		return nil, false
	}
	connection, okSelect := service.ConnectionRepository.SelectConnection(followerId, profileId, false)
	if okSelect && connection == nil {
		return connection, false
	}
	resConnection, ok := service.ConnectionRepository.DeleteConnection(followerId, profileId)
	if ok {
		return resConnection, true
	} else {
		return connection, false
	}
}

func (service *Service) GetAllFollowRequests(id uint) *[]dto.UserDTO {
	var result = service.ConnectionRepository.GetAllFollowRequests(id)
	var ret = make([]dto.UserDTO, 0) // 0, :)
	for _, profileId := range *result {
		var p model2.Profile
		profileHost, profilePort := util.GetProfileHostAndPort()
		resp, err := util.CrossServiceRequest(context.Background(), http.MethodGet,
			util.GetCrossServiceProtocol()+"://"+profileHost+":"+profilePort+"/get-by-id/"+util.Uint2String(profileId),
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

func (service *Service) ToggleNotifyComment(followerId, profileId uint) (*model.ConnectionEdge, bool) {
	if service.IsInBlockingRelationship(followerId, profileId) {
		return nil, false
	}
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

func (service *Service) ToggleNotifyStory(followerId, profileId uint) (*model.ConnectionEdge, bool) {
	if service.IsInBlockingRelationship(followerId, profileId) {
		return nil, false
	}
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

func (service *Service) ToggleNotifyPost(followerId, profileId uint) (*model.ConnectionEdge, bool) {
	if service.IsInBlockingRelationship(followerId, profileId) {
		return nil, false
	}
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

func (service *Service) ToggleCloseFriend(followerId, profileId uint) (*model.ConnectionEdge, bool) {
	if service.IsInBlockingRelationship(followerId, profileId) {
		return nil, false
	}
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

func (service *Service) ToggleMuted(followerId, profileId uint) (*model.ConnectionEdge, bool) {
	if service.IsInBlockingRelationship(followerId, profileId) {
		return nil, false
	}
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
