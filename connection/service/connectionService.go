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

func (service *Service) GetConnection(ctx context.Context, followerId, profileId uint) *model.ConnectionEdge {
	span := util.Tracer.StartSpanFromContext(ctx, "GetConnection-service")
	defer util.Tracer.FinishSpan(span)
	util.Tracer.LogFields(span, "service", fmt.Sprintf("servicing id %v %v\n", followerId, profileId))
	nextCtx := util.Tracer.ContextWithSpan(ctx, span)
	connection, _ := service.ConnectionRepository.SelectConnection(nextCtx, followerId, profileId, false)
	return connection
}

func (service *Service) GetConnectedProfiles(ctx context.Context, conn model.ConnectionEdge, excludeMuted, excludeBlocked bool) *[]dto.ConnectedProfileDTO {
	span := util.Tracer.StartSpanFromContext(ctx, "GetConnectedProfiles-service")
	defer util.Tracer.FinishSpan(span)
	util.Tracer.LogFields(span, "service", fmt.Sprintf("servicing ids %v %v\n", conn.PrimaryProfile, conn.SecondaryProfile))
	nextCtx := util.Tracer.ContextWithSpan(ctx, span)
	profiles := service.ConnectionRepository.GetConnectedProfiles(nextCtx, conn, excludeMuted, true)
	if profiles == nil {
		temp := make([]uint, 0)
		profiles = &temp
	}
	if !excludeBlocked {
		var final []uint
		blocking := service.ConnectionRepository.GetBlockedProfiles(nextCtx, conn.PrimaryProfile, false)
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
		invConnection, ok := service.ConnectionRepository.SelectConnection(nextCtx, val, conn.PrimaryProfile, false)
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

func (service *Service) GetProfilesInFollowRelationship(ctx context.Context,conn model.ConnectionEdge, excludeMuted, excludeBlocked bool, following bool) *[]dto.UserDTO {
	span := util.Tracer.StartSpanFromContext(ctx, "GetProfilesInFollowRelationship-service")
	defer util.Tracer.FinishSpan(span)
	util.Tracer.LogFields(span, "service", fmt.Sprintf("servicing ids %v %v\n", conn.PrimaryProfile, conn.SecondaryProfile))
	nextCtx := util.Tracer.ContextWithSpan(ctx, span)
	profiles := service.ConnectionRepository.GetConnectedProfiles(nextCtx, conn, excludeMuted, following)
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
		blocking := service.ConnectionRepository.GetBlockedProfiles(nextCtx, blockId, false)
		for _, val := range *profiles {
			if !contains(blocking, val) {
				final = append(final, val)
			}
		}
		profiles = &final
	}
	ret := make([]dto.UserDTO, 0)
	for _, val := range *profiles {
		p := util.GetProfile(nextCtx, val)
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

func (service *Service) UpdateConnection(ctx context.Context, id uint, conn model.ConnectionEdge) (*model.ConnectionEdge, bool) {
	span := util.Tracer.StartSpanFromContext(ctx, "UpdateConnection-service")
	defer util.Tracer.FinishSpan(span)
	util.Tracer.LogFields(span, "service", fmt.Sprintf("servicing ids %v %v\n", conn.PrimaryProfile, conn.SecondaryProfile))
	nextCtx := util.Tracer.ContextWithSpan(ctx, span)
	if id == conn.PrimaryProfile {
		return service.ConnectionRepository.UpdateConnection(nextCtx, &conn)
	} else {
		return nil, false
	}
}

func (service *Service) DeleteConnection(ctx context.Context, followerId, profileId uint) (*model.ConnectionEdge, bool) {
	span := util.Tracer.StartSpanFromContext(ctx, "DeleteConnection-service")
	defer util.Tracer.FinishSpan(span)
	util.Tracer.LogFields(span, "service", fmt.Sprintf("servicing id %v %v\n", followerId, profileId))
	nextCtx := util.Tracer.ContextWithSpan(ctx, span)
	return service.ConnectionRepository.DeleteConnection(nextCtx, followerId, profileId)
}

func (service *Service) FollowRequest(ctx context.Context, followerId, profileId uint) (*model.ConnectionEdge, bool) {
	span := util.Tracer.StartSpanFromContext(ctx, "FollowRequest-service")
	defer util.Tracer.FinishSpan(span)
	util.Tracer.LogFields(span, "service", fmt.Sprintf("servicing id %v %v\n", followerId, profileId))
	nextCtx := util.Tracer.ContextWithSpan(ctx, span)
	if service.IsInBlockingRelationship(nextCtx, followerId, profileId) {
		return nil, false
	}
	connection := service.ConnectionRepository.SelectOrCreateConnection(nextCtx, followerId, profileId)
	if connection.Approved {
		return nil, false
	}
	profile2 := util.GetProfile(nextCtx, profileId)
	if profile2 == nil {
		return nil, false
	}
	if profile2.ProfileSettings.IsPrivate == false {
		connection.Approved = true
		service.ConnectionRepository.CreateOrUpdateMessageRelationship(nextCtx, model.MessageEdge{
			PrimaryProfile:   profileId,
			SecondaryProfile: followerId,
			Approved:         true,
			NotifyMessage:    true,
		})
	} else {
		connection.ConnectionRequest = true
	}
	service.ConnectionRepository.CreateOrUpdateMessageRelationship(nextCtx, model.MessageEdge{
		PrimaryProfile:   followerId,
		SecondaryProfile: profileId,
		Approved:         true,
		NotifyMessage:    true,
	})
	resConnection, ok := service.ConnectionRepository.UpdateConnection(nextCtx, connection)
	if ok {
		return resConnection, true
	} else {
		return connection, false
	}
}

func (service *Service) ApproveConnection(ctx context.Context, followerId, profileId uint) (*model.ConnectionEdge, bool) {
	span := util.Tracer.StartSpanFromContext(ctx, "ApproveConnection-service")
	defer util.Tracer.FinishSpan(span)
	util.Tracer.LogFields(span, "service", fmt.Sprintf("servicing id %v %v\n", followerId, profileId))
	nextCtx := util.Tracer.ContextWithSpan(ctx, span)
	connection, okSelect := service.ConnectionRepository.SelectConnection(nextCtx, followerId, profileId, false)
	if okSelect && connection == nil {
		return connection, false
	}
	profile1 := util.GetProfile(nextCtx, followerId)
	profile2 := util.GetProfile(nextCtx, profileId)
	if profile1 == nil || profile2 == nil {
		return nil, false
	}
	if !connection.ConnectionRequest {
		return nil, false
	}
	connection.ConnectionRequest = false
	connection.Approved = true
	service.ConnectionRepository.CreateOrUpdateMessageRelationship(nextCtx, model.MessageEdge{
		PrimaryProfile:   profileId,
		SecondaryProfile: followerId,
		Approved:         true,
		NotifyMessage:    true,
	})
	return service.ConnectionRepository.UpdateConnection(nextCtx, connection)
}

func (service *Service) Unfollow(ctx context.Context, followerId, profileId uint) (*model.ConnectionEdge, bool) {
	span := util.Tracer.StartSpanFromContext(ctx, "Unfollow-service")
	defer util.Tracer.FinishSpan(span)
	util.Tracer.LogFields(span, "service", fmt.Sprintf("servicing id %v %v\n", followerId, profileId))
	nextCtx := util.Tracer.ContextWithSpan(ctx, span)
	if service.IsInBlockingRelationship(nextCtx, followerId, profileId) {
		return nil, false
	}
	connection, okSelect := service.ConnectionRepository.SelectConnection(nextCtx, followerId, profileId, false)
	if okSelect && connection == nil {
		return connection, false
	}
	resConnection, ok := service.ConnectionRepository.DeleteConnection(nextCtx, followerId, profileId)
	if ok {
		return resConnection, true
	} else {
		return connection, false
	}
}

func (service *Service) GetAllFollowRequests(ctx context.Context, id uint) *[]dto.UserDTO {
	span := util.Tracer.StartSpanFromContext(ctx, "GetAllFollowRequests-service")
	defer util.Tracer.FinishSpan(span)
	util.Tracer.LogFields(span, "service", fmt.Sprintf("servicing id %v\n", id))
	nextCtx := util.Tracer.ContextWithSpan(ctx, span)
	var result = service.ConnectionRepository.GetAllFollowRequests(nextCtx, id)
	var ret = make([]dto.UserDTO, 0) // 0, :)
	for _, profileId := range *result {
		var p model2.Profile
		profileHost, profilePort := util.GetProfileHostAndPort()
		resp, err := util.CrossServiceRequest(nextCtx, http.MethodGet,
			util.GetCrossServiceProtocol()+"://"+profileHost+":"+profilePort+"/get-by-id/"+util.Uint2String(profileId),
			nil, map[string]string{})
		if err != nil {
			util.Tracer.LogError(span, err)
			fmt.Println(err)
			return nil
		}
		body, err1 := ioutil.ReadAll(resp.Body)
		if err1 != nil {
			util.Tracer.LogError(span, err1)
			fmt.Println(err1)
			return nil
		}
		err = json.Unmarshal(body, &p)
		if err != nil {
			util.Tracer.LogError(span, err)
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

func (service *Service) ToggleNotifyComment(ctx context.Context, followerId, profileId uint) (*model.ConnectionEdge, bool) {
	span := util.Tracer.StartSpanFromContext(ctx, "ToggleNotifyComment-service")
	defer util.Tracer.FinishSpan(span)
	util.Tracer.LogFields(span, "service", fmt.Sprintf("servicing id %v %v\n", followerId, profileId))
	nextCtx := util.Tracer.ContextWithSpan(ctx, span)
	if service.IsInBlockingRelationship(nextCtx, followerId, profileId) {
		return nil, false
	}
	connection, okSelect := service.ConnectionRepository.SelectConnection(nextCtx, followerId, profileId, false)
	if okSelect && connection == nil {
		return connection, false
	}
	connection.NotifyComment = !connection.NotifyComment
	resConnection, ok := service.ConnectionRepository.UpdateConnection(nextCtx, connection)
	if ok {
		return resConnection, true
	} else {
		return connection, false
	}
}

func (service *Service) ToggleNotifyStory(ctx context.Context, followerId, profileId uint) (*model.ConnectionEdge, bool) {
	span := util.Tracer.StartSpanFromContext(ctx, "ToggleNotifyStory-service")
	defer util.Tracer.FinishSpan(span)
	util.Tracer.LogFields(span, "service", fmt.Sprintf("servicing id %v %v\n", followerId, profileId))
	nextCtx := util.Tracer.ContextWithSpan(ctx, span)
	if service.IsInBlockingRelationship(nextCtx, followerId, profileId) {
		return nil, false
	}
	connection, okSelect := service.ConnectionRepository.SelectConnection(nextCtx, followerId, profileId, false)
	if okSelect && connection == nil {
		return connection, false
	}
	connection.NotifyStory = !connection.NotifyStory
	resConnection, ok := service.ConnectionRepository.UpdateConnection(nextCtx, connection)
	if ok {
		return resConnection, true
	} else {
		return connection, false
	}
}

func (service *Service) ToggleNotifyPost(ctx context.Context, followerId, profileId uint) (*model.ConnectionEdge, bool) {
	span := util.Tracer.StartSpanFromContext(ctx, "ToggleNotifyPost-service")
	defer util.Tracer.FinishSpan(span)
	util.Tracer.LogFields(span, "service", fmt.Sprintf("servicing id %v %v\n", followerId, profileId))
	nextCtx := util.Tracer.ContextWithSpan(ctx, span)
	if service.IsInBlockingRelationship(nextCtx, followerId, profileId) {
		return nil, false
	}
	connection, okSelect := service.ConnectionRepository.SelectConnection(nextCtx, followerId, profileId, false)
	if okSelect && connection == nil {
		return connection, false
	}
	connection.NotifyPost = !connection.NotifyPost
	resConnection, ok := service.ConnectionRepository.UpdateConnection(nextCtx, connection)
	if ok {
		return resConnection, true
	} else {
		return connection, false
	}
}

func (service *Service) ToggleCloseFriend(ctx context.Context, followerId, profileId uint) (*model.ConnectionEdge, bool) {
	span := util.Tracer.StartSpanFromContext(ctx, "ToggleCloseFriend-service")
	defer util.Tracer.FinishSpan(span)
	util.Tracer.LogFields(span, "service", fmt.Sprintf("servicing id %v %v\n", followerId, profileId))
	nextCtx := util.Tracer.ContextWithSpan(ctx, span)
	if service.IsInBlockingRelationship(nextCtx, followerId, profileId) {
		return nil, false
	}
	connection, okSelect := service.ConnectionRepository.SelectConnection(nextCtx, followerId, profileId, false)
	if okSelect && connection == nil {
		return connection, false
	}
	connection.CloseFriend = !connection.CloseFriend
	resConnection, ok := service.ConnectionRepository.UpdateConnection(nextCtx, connection)
	if ok {
		return resConnection, true
	} else {
		return connection, false
	}
}

func (service *Service) ToggleMuted(ctx context.Context, followerId, profileId uint) (*model.ConnectionEdge, bool) {
	span := util.Tracer.StartSpanFromContext(ctx, "ToggleMuted-service")
	defer util.Tracer.FinishSpan(span)
	util.Tracer.LogFields(span, "service", fmt.Sprintf("servicing id %v %v\n", followerId, profileId))
	nextCtx := util.Tracer.ContextWithSpan(ctx, span)
	if service.IsInBlockingRelationship(nextCtx, followerId, profileId) {
		return nil, false
	}
	connection, okSelect := service.ConnectionRepository.SelectConnection(nextCtx, followerId, profileId, false)
	if okSelect && connection == nil {
		return connection, false
	}
	connection.Muted = !connection.Muted
	resConnection, ok := service.ConnectionRepository.UpdateConnection(nextCtx, connection)
	if ok {
		return resConnection, true
	} else {
		return connection, false
	}
}
