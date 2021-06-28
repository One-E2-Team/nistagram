package service

import (
	"math"
	"nistagram/connection/dto"
	"nistagram/connection/model"
	"nistagram/util"
)

func (service *Service) AddOrUpdateProfile(profile model.ProfileVertex) (*model.ProfileVertex, bool) {
	ret := service.ConnectionRepository.CreateOrUpdateProfile(profile)
	return ret, ret != nil && ret.ProfileID == profile.ProfileID && ret.Deleted == profile.Deleted
}

func (service *Service) GetRecommendations(id uint) (*[]dto.ProfileRecommendationDTO, bool) {
	var degrees *map[uint]*[]uint
	var ok bool
	degrees, ok = service.ConnectionRepository.FindConnectionDegreeForRecommendation(id)
	if !ok || degrees == nil {
		return nil, false
	}
	var ret = make([]dto.ProfileRecommendationDTO, 0)
	for key, degVals := range *degrees {
		ret = append(ret, dto.ProfileRecommendationDTO{
			Username:   util.GetProfile(key).Username,
			ProfileID:  key,
			Confidence: rankingAlgorithm(degVals),
		})
	}
	return &ret, ok
}

func rankingAlgorithm(lst *[]uint) float64 {
	var rank float64 = 0.0
	for _, val := range *lst {
		rank += math.Pow(float64(val),-2)
	}
	return rank
}