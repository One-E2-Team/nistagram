package service

import (
	"context"
	"fmt"
	"nistagram/connection/model"
	"nistagram/util"
)

func (service *Service) IsInBlockingRelationship(ctx context.Context, id1, id2 uint) bool {
	span := util.Tracer.StartSpanFromContext(ctx, "IsInBlockingRelationship-service")
	defer util.Tracer.FinishSpan(span)
	util.Tracer.LogFields(span, "service", fmt.Sprintf("servicing id %v %v\n", id1, id2))
	nextCtx := util.Tracer.ContextWithSpan(ctx, span)
	lst := service.ConnectionRepository.GetBlockedProfiles(nextCtx, id1, false)
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

func (service *Service) IsBlocked(ctx context.Context, id1, id2 uint) bool {
	span := util.Tracer.StartSpanFromContext(ctx, "IsBlocked-service")
	defer util.Tracer.FinishSpan(span)
	util.Tracer.LogFields(span, "service", fmt.Sprintf("servicing id %v %v\n", id1, id2))
	nextCtx := util.Tracer.ContextWithSpan(ctx, span)
	lst := service.ConnectionRepository.GetBlockedProfiles(nextCtx, id1, true)
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

func (service *Service) ToggleBlock(ctx context.Context, followerId, profileId uint) (*model.BlockEdge, bool) {
	span := util.Tracer.StartSpanFromContext(ctx, "ToggleBlock-service")
	defer util.Tracer.FinishSpan(span)
	util.Tracer.LogFields(span, "service", fmt.Sprintf("servicing id %v %v\n",followerId, profileId))
	nextCtx := util.Tracer.ContextWithSpan(ctx, span)
	block, ok := service.ConnectionRepository.SelectBlock(nextCtx, followerId, profileId)
	if !ok || block == nil {
		service.ConnectionRepository.DeleteConnection(nextCtx, followerId, profileId)
		service.ConnectionRepository.DeleteConnection(nextCtx, profileId, followerId)
		service.ConnectionRepository.DeleteMessage(nextCtx, followerId, profileId)
		service.ConnectionRepository.DeleteMessage(nextCtx, profileId, followerId)
		block, ok = service.ConnectionRepository.CreateBlock(nextCtx, followerId, profileId)
	} else {
		block, ok = service.ConnectionRepository.DeleteBlock(nextCtx, followerId, profileId)
	}
	return block, ok
}


func (service *Service) GetBlockingRelationships(ctx context.Context, u uint) *[]uint {
	span := util.Tracer.StartSpanFromContext(ctx, "GetBlockingRelationships-service")
	defer util.Tracer.FinishSpan(span)
	util.Tracer.LogFields(span, "service", fmt.Sprintf("servicing id %v\n", u))
	nextCtx := util.Tracer.ContextWithSpan(ctx, span)
	return service.ConnectionRepository.GetBlockedProfiles(nextCtx, u, false)
}