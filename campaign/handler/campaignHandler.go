package handler

import (
	"context"
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
	span := util.Tracer.StartSpanFromRequest("CreateCampaign-handler", r)
	defer util.Tracer.FinishSpan(span)
	util.Tracer.LogFields(span, "handler", fmt.Sprintf("handling %s\n", r.URL.Path))
	ctx := util.Tracer.ContextWithSpan(context.Background(), span)

	var campaign dto.CampaignDTO
	if err := json.NewDecoder(r.Body).Decode(&campaign); err != nil {
		util.Tracer.LogError(span, err)
		fmt.Println(err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	if result, err := handler.CampaignService.CreateCampaign(ctx, util.GetLoggedUserIDFromToken(r), campaign); err != nil {
		util.Tracer.LogError(span, err)
		w.WriteHeader(http.StatusInternalServerError)
	} else {
		w.WriteHeader(http.StatusCreated)
		json.NewEncoder(w).Encode(result)
	}

	w.Header().Set("Content-Type", "application/json")
}

func (handler CampaignHandler) UpdateCampaignParameters(w http.ResponseWriter, r *http.Request) {
	span := util.Tracer.StartSpanFromRequest("UpdateCampaignParameters-handler", r)
	defer util.Tracer.FinishSpan(span)
	util.Tracer.LogFields(span, "handler", fmt.Sprintf("handling %s\n", r.URL.Path))
	ctx := util.Tracer.ContextWithSpan(context.Background(), span)

	params := mux.Vars(r)
	campaignId := params["id"]

	var camParams dto.CampaignParametersDTO

	if err := json.NewDecoder(r.Body).Decode(&camParams); err != nil {
		util.Tracer.LogError(span, err)
		fmt.Println(err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	if err := handler.CampaignService.UpdateCampaignParameters(ctx, util.String2Uint(campaignId), camParams); err != nil {
		util.Tracer.LogError(span, err)
		w.WriteHeader(http.StatusInternalServerError)
	} else {
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte("{\"message\":\"ok\"}"))
	}

	w.Header().Set("Content-Type", "application/json")
}

func (handler CampaignHandler) DeleteCampaign(w http.ResponseWriter, r *http.Request) {
	span := util.Tracer.StartSpanFromRequest("DeleteCampaign-handler", r)
	defer util.Tracer.FinishSpan(span)
	util.Tracer.LogFields(span, "handler", fmt.Sprintf("handling %s\n", r.URL.Path))
	ctx := util.Tracer.ContextWithSpan(context.Background(), span)

	params := mux.Vars(r)
	campaignId := params["id"]
	switch err := handler.CampaignService.DeleteCampaign(ctx, util.String2Uint(campaignId)); err {
	case gorm.ErrRecordNotFound:
		util.Tracer.LogError(span, fmt.Errorf("campaign not found"))
		w.WriteHeader(http.StatusNotFound)
	case nil:
		w.WriteHeader(http.StatusOK)
	default:
		util.Tracer.LogError(span, err)
		w.WriteHeader(http.StatusInternalServerError)
	}
}

func (handler *CampaignHandler) GetMyCampaigns(w http.ResponseWriter, r *http.Request) {
	span := util.Tracer.StartSpanFromRequest("GetMyCampaigns-handler", r)
	defer util.Tracer.FinishSpan(span)
	util.Tracer.LogFields(span, "handler", fmt.Sprintf("handling %s\n", r.URL.Path))
	ctx := util.Tracer.ContextWithSpan(context.Background(), span)

	campaigns, err := handler.CampaignService.GetMyCampaigns(ctx, util.GetLoggedUserIDFromToken(r))
	if err != nil {
		util.Tracer.LogError(span, err)
		fmt.Println(err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	if js, err := json.Marshal(campaigns); err != nil {
		util.Tracer.LogError(span, err)
		w.WriteHeader(http.StatusInternalServerError)
	} else {
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write(js)
	}
	w.Header().Set("Content-Type", "application/json")
}

func (handler *CampaignHandler) GetAllInterests(w http.ResponseWriter, r *http.Request) {
	span := util.Tracer.StartSpanFromRequest("GetAllInterests-handler", r)
	defer util.Tracer.FinishSpan(span)
	util.Tracer.LogFields(span, "handler", fmt.Sprintf("handling %s\n", r.URL.Path))
	ctx := util.Tracer.ContextWithSpan(context.Background(), span)

	interests, err := handler.CampaignService.GetAllInterests(ctx)
	if err != nil {
		util.Tracer.LogError(span, err)
		fmt.Println(err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	if js, err := json.Marshal(interests); err != nil {
		util.Tracer.LogError(span, err)
		w.WriteHeader(http.StatusInternalServerError)
	} else {
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write(js)
	}
	w.Header().Set("Content-Type", "application/json")
}

func (handler CampaignHandler) GetCurrentlyValidInterests(w http.ResponseWriter, r *http.Request) {
	span := util.Tracer.StartSpanFromRequest("GetCurrentlyValidInterests-handler", r)
	defer util.Tracer.FinishSpan(span)
	util.Tracer.LogFields(span, "handler", fmt.Sprintf("handling %s\n", r.URL.Path))
	ctx := util.Tracer.ContextWithSpan(context.Background(), span)

	params := mux.Vars(r)
	campaignId := params["campaignId"]

	interests, err := handler.CampaignService.GetCurrentlyValidInterests(ctx, util.String2Uint(campaignId))
	if err != nil {
		util.Tracer.LogError(span, err)
		fmt.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	err = json.NewEncoder(w).Encode(interests)
	if err != nil {
		util.Tracer.LogError(span, err)
		fmt.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json")
}

func (handler CampaignHandler) GetCampaignByIdForMonitoring(w http.ResponseWriter, r *http.Request) {
	span := util.Tracer.StartSpanFromRequest("GetCampaignByIdForMonitoring-handler", r)
	defer util.Tracer.FinishSpan(span)
	util.Tracer.LogFields(span, "handler", fmt.Sprintf("handling %s\n", r.URL.Path))
	ctx := util.Tracer.ContextWithSpan(context.Background(), span)

	params := mux.Vars(r)
	campaignId := util.String2Uint(params["id"])

	campaign, err := handler.CampaignService.GetCampaignByIdForMonitoring(ctx, campaignId)
	if err != nil {
		util.Tracer.LogError(span, err)
		fmt.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	err = json.NewEncoder(w).Encode(campaign)
	if err != nil {
		util.Tracer.LogError(span, err)
		fmt.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json")
}

func (handler CampaignHandler) GetLastActiveParametersForCampaign(w http.ResponseWriter, r *http.Request) {
	span := util.Tracer.StartSpanFromRequest("GetLastActiveParametersForCampaign-handler", r)
	defer util.Tracer.FinishSpan(span)
	util.Tracer.LogFields(span, "handler", fmt.Sprintf("handling %s\n", r.URL.Path))
	ctx := util.Tracer.ContextWithSpan(context.Background(), span)

	params := mux.Vars(r)
	campaignId := util.String2Uint(params["id"])

	switch campaignParams, err := handler.CampaignService.GetLastActiveParametersForCampaign(ctx, campaignId); err {
	case gorm.ErrRecordNotFound:
		w.WriteHeader(http.StatusNotFound)
		return
	case nil:
		if err2 := json.NewEncoder(w).Encode(campaignParams); err2 != nil {
			util.Tracer.LogError(span, err2)
			w.WriteHeader(http.StatusInternalServerError)
		} else {
			w.WriteHeader(http.StatusOK)
		}
		return
	default:
		util.Tracer.LogError(span, err)
		w.WriteHeader(http.StatusInternalServerError)
	}
}

func (handler *CampaignHandler) GetAvailableCampaignsForUser(w http.ResponseWriter, r *http.Request) {
	span := util.Tracer.StartSpanFromRequest("GetAvailableCampaignsForUser-handler", r)
	defer util.Tracer.FinishSpan(span)
	util.Tracer.LogFields(span, "handler", fmt.Sprintf("handling %s\n", r.URL.Path))
	ctx := util.Tracer.ContextWithSpan(context.Background(), span)

	params := mux.Vars(r)
	var followingProfiles []util.FollowingProfileDTO

	if err := json.NewDecoder(r.Body).Decode(&followingProfiles); err != nil {
		util.Tracer.LogError(span, err)
		fmt.Println(err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	postIDs, influencerIDs, campaignIDs, err := handler.CampaignService.GetAvailableCampaignsForUser(ctx, util.String2Uint(params["profileID"]), followingProfiles)
	if err != nil {
		util.Tracer.LogError(span, err)
		fmt.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	if len(postIDs) != len(influencerIDs) || len(postIDs) != len(campaignIDs) || len(influencerIDs) != len(campaignIDs){
		util.Tracer.LogError(span, fmt.Errorf("bad list sizes"))
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
		util.Tracer.LogError(span, err)
		fmt.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
	_, _ = w.Write(js)
	w.Header().Set("Content-Type", "application/json")
}

func (handler *CampaignHandler) GetAcceptedCampaignsForInfluencer(w http.ResponseWriter, r *http.Request) {
	span := util.Tracer.StartSpanFromRequest("GetAcceptedCampaignsForInfluencer-handler", r)
	defer util.Tracer.FinishSpan(span)
	util.Tracer.LogFields(span, "handler", fmt.Sprintf("handling %s\n", r.URL.Path))
	ctx := util.Tracer.ContextWithSpan(context.Background(), span)

	params := mux.Vars(r)
	postIDs, influencerIDs, campaignIDs, err := handler.CampaignService.GetAcceptedCampaignsForInfluencer(ctx, util.String2Uint(params["influencerID"]))
	if err != nil {
		util.Tracer.LogError(span, err)
		fmt.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	if len(postIDs) != len(influencerIDs) || len(postIDs) != len(campaignIDs) || len(influencerIDs) != len(campaignIDs){
		util.Tracer.LogError(span, fmt.Errorf("bad list sizes"))
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
		util.Tracer.LogError(span, err)
		fmt.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
	_, _ = w.Write(js)
	w.Header().Set("Content-Type", "application/json")
}

func (handler CampaignHandler) UpdateCampaignRequest(w http.ResponseWriter, r *http.Request) {
	span := util.Tracer.StartSpanFromRequest("UpdateCampaignRequest-handler", r)
	defer util.Tracer.FinishSpan(span)
	util.Tracer.LogFields(span, "handler", fmt.Sprintf("handling %s\n", r.URL.Path))
	ctx := util.Tracer.ContextWithSpan(context.Background(), span)

	params := mux.Vars(r)
	type ReqBody struct {
		Accepted bool `json:"accepted"`
	}
	var data ReqBody
	requestId := params["id"]
	loggedUserId := util.GetLoggedUserIDFromToken(r)

	if err := json.NewDecoder(r.Body).Decode(&data); err != nil {
		util.Tracer.LogError(span, err)
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
		util.Tracer.LogError(span, fmt.Errorf("not valid request status"))
		fmt.Println("Not valid request body")
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	switch err := handler.CampaignService.UpdateCampaignRequest(ctx, loggedUserId,requestId,status); err {
	case gorm.ErrRecordNotFound:
		util.Tracer.LogError(span, fmt.Errorf("campaign request not found"))
		w.WriteHeader(http.StatusNotFound)
	case nil:
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte("{\"message\":\"ok\"}"))
	default:
		util.Tracer.LogError(span, err)
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Println(err)
	}

}

func (handler CampaignHandler) GetMySentCampaignsRequest(w http.ResponseWriter, r *http.Request) {
	span := util.Tracer.StartSpanFromRequest("GetMySentCampaignsRequest-handler", r)
	defer util.Tracer.FinishSpan(span)
	util.Tracer.LogFields(span, "handler", fmt.Sprintf("handling %s\n", r.URL.Path))
	ctx := util.Tracer.ContextWithSpan(context.Background(), span)

	switch ret, err := handler.CampaignService.GetActiveCampaignsRequestsForProfileId(ctx, util.GetLoggedUserIDFromToken(r)); err {
	case gorm.ErrRecordNotFound:
		fmt.Println(err)
		w.WriteHeader(http.StatusNotFound)
	case nil:
		js, err := json.Marshal(ret)
		if err != nil {
			util.Tracer.LogError(span, err)
			fmt.Println(err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write(js)
		w.Header().Set("Content-Type", "application/json")
	default:
		util.Tracer.LogError(span, err)
		fmt.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
	}
}

