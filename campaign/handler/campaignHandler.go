package handler

import (
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"gorm.io/gorm"
	"net/http"
	"nistagram/campaign/dto"
	"nistagram/campaign/model"
	"nistagram/campaign/service"
	"nistagram/util"
)

type CampaignHandler struct {
	CampaignService *service.CampaignService
}

func (handler CampaignHandler) CreateCampaign(w http.ResponseWriter, r *http.Request) {
	var campaign dto.CampaignDTO
	if err := json.NewDecoder(r.Body).Decode(&campaign); err != nil {
		fmt.Println(err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	if result, err := handler.CampaignService.CreateCampaign(util.GetLoggedUserIDFromToken(r), campaign); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
	} else {
		w.WriteHeader(http.StatusCreated)
		json.NewEncoder(w).Encode(result)
	}

	w.Header().Set("Content-Type", "application/json")
}

func (handler CampaignHandler) UpdateCampaignParameters(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	campaignId := params["id"]

	var camParams dto.CampaignParametersDTO

	if err := json.NewDecoder(r.Body).Decode(&camParams); err != nil {
		fmt.Println(err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	if err := handler.CampaignService.UpdateCampaignParameters(util.String2Uint(campaignId), camParams); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
	} else {
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte("{\"message\":\"ok\"}"))
	}

	w.Header().Set("Content-Type", "application/json")
}

func (handler CampaignHandler) DeleteCampaign(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	campaignId := params["id"]
	switch err := handler.CampaignService.DeleteCampaign(util.String2Uint(campaignId)); err {
	case gorm.ErrRecordNotFound:
		w.WriteHeader(http.StatusNotFound)
	case nil:
		w.WriteHeader(http.StatusOK)
	default:
		w.WriteHeader(http.StatusInternalServerError)
	}
}

func (handler *CampaignHandler) GetMyCampaigns(w http.ResponseWriter, r *http.Request) {
	campaigns, err := handler.CampaignService.GetMyCampaigns(util.GetLoggedUserIDFromToken(r))
	if err != nil {
		fmt.Println(err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	if js, err := json.Marshal(campaigns); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
	} else {
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write(js)
	}
	w.Header().Set("Content-Type", "application/json")
}

func (handler *CampaignHandler) GetAllInterests(w http.ResponseWriter, r *http.Request) {
	interests, err := handler.CampaignService.GetAllInterests()
	if err != nil {
		fmt.Println(err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	if js, err := json.Marshal(interests); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
	} else {
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write(js)
	}
	w.Header().Set("Content-Type", "application/json")
}

func (handler CampaignHandler) GetCurrentlyValidInterests(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	campaignId := params["campaignId"]

	interests, err := handler.CampaignService.GetCurrentlyValidInterests(util.String2Uint(campaignId))
	if err != nil {
		fmt.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	err = json.NewEncoder(w).Encode(interests)
	if err != nil {
		fmt.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json")
}

func (handler CampaignHandler) GetCampaignByIdForMonitoring(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	campaignId := util.String2Uint(params["id"])

	campaign, err := handler.CampaignService.GetCampaignByIdForMonitoring(campaignId)
	if err != nil {
		fmt.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	err = json.NewEncoder(w).Encode(campaign)
	if err != nil {
		fmt.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json")
}

func (handler CampaignHandler) GetLastActiveParametersForCampaign(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	campaignId := util.String2Uint(params["id"])

	switch campaignParams, err := handler.CampaignService.GetLastActiveParametersForCampaign(campaignId); err {
	case gorm.ErrRecordNotFound:
		w.WriteHeader(http.StatusNotFound)
		return
	case nil:
		if err2 := json.NewEncoder(w).Encode(campaignParams); err2 != nil {
			w.WriteHeader(http.StatusInternalServerError)
		} else {
			w.WriteHeader(http.StatusOK)
		}
		return
	default:
		w.WriteHeader(http.StatusInternalServerError)
	}
}

func (handler *CampaignHandler) GetAvailableCampaignsForUser(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	var followingProfiles []util.FollowingProfileDTO

	if err := json.NewDecoder(r.Body).Decode(&followingProfiles); err != nil {
		fmt.Println(err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	postIDs, influencerIDs, campaignIDs, err := handler.CampaignService.GetAvailableCampaignsForUser(util.String2Uint(params["profileID"]), followingProfiles)
	if err != nil {
		fmt.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	if len(postIDs) != len(influencerIDs) || len(postIDs) != len(campaignIDs) || len(influencerIDs) != len(campaignIDs){
		fmt.Println("BAD LIST SIZES")
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	ret := make([]dto.SponsoredPostsDTO, 0)
	for i, postID := range postIDs {
		ret = append(ret, dto.SponsoredPostsDTO{
			PostID:       postID,
			InfluencerID: influencerIDs[i],
			CampaignID:   campaignIDs[i],
		})
	}
	js, err := json.Marshal(ret)
	if err != nil {
		fmt.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
	_, _ = w.Write(js)
	w.Header().Set("Content-Type", "application/json")
}

func (handler CampaignHandler) UpdateCampaignRequest(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	type ReqBody struct {
		Accepted bool `json:"accepted"`
	}
	var data ReqBody
	requestId := params["id"]

	if err := json.NewDecoder(r.Body).Decode(&data); err != nil {
		fmt.Println(err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	var status model.RequestStatus
	switch data.Accepted {
	case true:
		status = model.ACCEPTED
	case false:
		status = model.DECLINED
	default:
		fmt.Println("Not valid request body")
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	switch err := handler.CampaignService.UpdateCampaignRequest(requestId,status); err {
	case gorm.ErrRecordNotFound:
		w.WriteHeader(http.StatusNotFound)
	case nil:
		w.WriteHeader(http.StatusOK)
	default:
		w.WriteHeader(http.StatusInternalServerError)
	}

}

