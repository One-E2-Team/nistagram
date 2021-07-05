package handler

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"net/http"
	"nistagram/monitoring/dto"
	"nistagram/monitoring/service"
	"nistagram/util"
)

type Handler struct {
	MonitoringService *service.MonitoringService
}

func (handler *Handler) CreateEventInfluencer(w http.ResponseWriter, r *http.Request) {
	span := util.Tracer.StartSpanFromRequest("CreateEventInfluencer-handler", r)
	defer util.Tracer.FinishSpan(span)
	util.Tracer.LogFields(span, "handler", fmt.Sprintf("handling %s\n", r.URL.Path))
	ctx := util.Tracer.ContextWithSpan(context.Background(), span)
	var eventDto dto.EventDTO
	if err := json.NewDecoder(r.Body).Decode(&eventDto); err != nil{
		util.Tracer.LogError(span, err)
		fmt.Println(err)
		w.WriteHeader(http.StatusBadRequest)
		_, _ = w.Write([]byte("{\"message\":\"error\"}"))
		return
	}

	fmt.Println(eventDto)

	err := handler.MonitoringService.CreateEventInfluencer(ctx, eventDto)
	if err != nil{
		util.Tracer.LogError(span, err)
		fmt.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		_, _ = w.Write([]byte("{\"message\":\"error\"}"))
		return
	}

	w.WriteHeader(http.StatusOK)
	_, _ = w.Write([]byte("{\"message\":\"ok\"}"))
	w.Header().Set("Content-Type", "application/json")
}

func (handler *Handler) CreateEventTargetGroup(w http.ResponseWriter, r *http.Request) {
	span := util.Tracer.StartSpanFromRequest("CreateEventTargetGroup-handler", r)
	defer util.Tracer.FinishSpan(span)
	util.Tracer.LogFields(span, "handler", fmt.Sprintf("handling %s\n", r.URL.Path))
	ctx := util.Tracer.ContextWithSpan(context.Background(), span)
	var eventDto dto.EventDTO
	if err := json.NewDecoder(r.Body).Decode(&eventDto); err != nil{
		util.Tracer.LogError(span, err)
		fmt.Println(err)
		w.WriteHeader(http.StatusBadRequest)
		_, _ = w.Write([]byte("{\"message\":\"error\"}"))
		return
	}

	fmt.Println(eventDto)

	err := handler.MonitoringService.CreateEventTargetGroup(ctx, eventDto)
	if err != nil{
		util.Tracer.LogError(span, err)
		fmt.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		_, _ = w.Write([]byte("{\"message\":\"error\"}"))
		return
	}

	w.WriteHeader(http.StatusOK)
	_, _ = w.Write([]byte("{\"message\":\"ok\"}"))
	w.Header().Set("Content-Type", "application/json")
}

func (handler *Handler) VisitSite(w http.ResponseWriter, r *http.Request) {
	span := util.Tracer.StartSpanFromRequest("VisitSite-handler", r)
	defer util.Tracer.FinishSpan(span)
	util.Tracer.LogFields(span, "handler", fmt.Sprintf("handling %s\n", r.URL.Path))
	ctx := util.Tracer.ContextWithSpan(context.Background(), span)
	vars := mux.Vars(r)
	campaignId := util.String2Uint(vars["campaignId"])
	influencerId := util.String2Uint(vars["influencerId"])
	mediaId := vars["mediaId"]

	loggedUserId := util.GetLoggedUserIDFromToken(r)

	website, err := handler.MonitoringService.VisitSite(ctx, campaignId, influencerId, loggedUserId, mediaId)
	if err != nil{
		util.Tracer.LogError(span, err)
		fmt.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		_, _ = w.Write([]byte("{\"message\":\"error\"}"))
		return
	}
	
	http.Redirect(w, r, website, http.StatusTemporaryRedirect)
}

func (handler *Handler) GetCampaignStatistics(w http.ResponseWriter, r *http.Request) {
	span := util.Tracer.StartSpanFromRequest("GetCampaignStatistics-handler", r)
	defer util.Tracer.FinishSpan(span)
	util.Tracer.LogFields(span, "handler", fmt.Sprintf("handling %s\n", r.URL.Path))
	ctx := util.Tracer.ContextWithSpan(context.Background(), span)
	vars := mux.Vars(r)
	campaignId := util.String2Uint(vars["campaignId"])

	result, err := handler.MonitoringService.GetCampaignStatistics(ctx, campaignId)
	if err != nil{
		util.Tracer.LogError(span, err)
		fmt.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		_, _ = w.Write([]byte("{\"message\":\"error\"}"))
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(result)
	w.Header().Set("Content-Type", "application/json")
}
