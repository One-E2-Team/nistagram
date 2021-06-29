package handler

import (
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
	var eventDto dto.EventDTO
	if err := json.NewDecoder(r.Body).Decode(&eventDto); err != nil{
		fmt.Println(err)
		w.WriteHeader(http.StatusBadRequest)
		_, _ = w.Write([]byte("{\"message\":\"error\"}"))
		return
	}

	fmt.Println(eventDto)

	err := handler.MonitoringService.CreateEventInfluencer(eventDto)
	if err != nil{
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
	var eventDto dto.EventDTO
	if err := json.NewDecoder(r.Body).Decode(&eventDto); err != nil{
		fmt.Println(err)
		w.WriteHeader(http.StatusBadRequest)
		_, _ = w.Write([]byte("{\"message\":\"error\"}"))
		return
	}

	fmt.Println(eventDto)

	err := handler.MonitoringService.CreateEventTargetGroup(eventDto)
	if err != nil{
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
	vars := mux.Vars(r)
	campaignId := util.String2Uint(vars["campaignId"])
	influencerId := util.String2Uint(vars["influencerId"])
	mediaId := util.String2Uint(vars["mediaId"])

	loggedUserId := util.GetLoggedUserIDFromToken(r)

	website, err := handler.MonitoringService.VisitSite(campaignId, influencerId, loggedUserId, mediaId)
	if err != nil{
		fmt.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		_, _ = w.Write([]byte("{\"message\":\"error\"}"))
		return
	}

	http.Redirect(w, r, website, http.StatusSeeOther)
	w.Header().Set("Content-Type", "application/json")
}
