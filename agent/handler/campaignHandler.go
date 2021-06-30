package handler

import (
	"fmt"
	"io"
	"net/http"
	"nistagram/agent/service"
	"nistagram/agent/util"
)

type CampaignHandler struct {
	CampaignService *service.CampaignService
}

func (handler *CampaignHandler) GetMyCampaigns(w http.ResponseWriter, r *http.Request) {
	resp, err := handler.CampaignService.GetMyCampaigns()
	if err != nil {
		fmt.Println(err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	w.WriteHeader(http.StatusOK)
	_, _ = w.Write(util.GetResponseJSON(*resp))
	w.Header().Set("Content-Type", "application/json")
}

func (handler *CampaignHandler) CreateCampaign(w http.ResponseWriter, r *http.Request) {
	body, err := io.ReadAll(r.Body)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer func(Body io.ReadCloser) {
		_ = Body.Close()
	}(r.Body)
	resp, err := handler.CampaignService.CreateCampaign(body)
	if err != nil {
		fmt.Println(err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	w.WriteHeader(http.StatusOK)
	_, _ = w.Write(util.GetResponseJSON(*resp))
	w.Header().Set("Content-Type", "application/json")
}

func (handler *CampaignHandler) GetInterests(w http.ResponseWriter, r *http.Request) {
	resp, err := handler.CampaignService.GetInterests()
	if err != nil {
		fmt.Println(err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	w.WriteHeader(http.StatusOK)
	_, _ = w.Write(util.GetResponseJSON(*resp))
	w.Header().Set("Content-Type", "application/json")
}
