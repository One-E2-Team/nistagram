package handler

import (
	"encoding/json"
	"fmt"
	"net/http"
	"nistagram/monitoring/dto"
	"nistagram/monitoring/service"
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
