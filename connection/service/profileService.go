package service

import (
	"context"
	"fmt"
	"math"
	"nistagram/connection/dto"
	"nistagram/connection/model"
	"nistagram/util"
)

func (service *Service) AddOrUpdateProfile(profile model.ProfileVertex) (*model.ProfileVertex, bool) {
	ret := service.ConnectionRepository.CreateOrUpdateProfile(profile)
	return ret, ret != nil && ret.ProfileID == profile.ProfileID && ret.Deleted == profile.Deleted
}

func (service *Service) GetRecommendations(ctx context.Context, id uint) (*[]dto.ProfileRecommendationDTO, bool) {
	span := util.Tracer.StartSpanFromContext(ctx, "GetRecommendations-service")
	defer util.Tracer.FinishSpan(span)
	util.Tracer.LogFields(span, "service", fmt.Sprintf("servicing id %v\n", id))
	nextCtx := util.Tracer.ContextWithSpan(ctx, span)
	var degrees *map[uint]*[]uint
	var ok bool
	degrees, ok = service.ConnectionRepository.FindConnectionDegreeForRecommendation(nextCtx, id)
	if !ok || degrees == nil {
		util.Tracer.LogError(span, fmt.Errorf("could not get valid repository response"))
		return nil, false
	}
	var ret = make([]dto.ProfileRecommendationDTO, 0)
	for key, degVals := range *degrees {
		ret = append(ret, dto.ProfileRecommendationDTO{
			Username:   util.GetProfile(nextCtx, key).Username,
			ProfileID:  key,
			Confidence: rankingAlgorithm(nextCtx, degVals),
		})
	}
	return &ret, ok
}

func rankingAlgorithm(ctx context.Context, lst *[]uint) float64 {
	span := util.Tracer.StartSpanFromContext(ctx, "rankingAlgorithm-service")
	defer util.Tracer.FinishSpan(span)
	util.Tracer.LogFields(span, "service", fmt.Sprintf("servicing lost %v\n", lst))
	var rank float64 = 0.0
	for _, val := range *lst {
		rank += math.Pow(float64(val),-2)
	}
	return rank
}