package handler

import (
	"fmt"
	"io"
	"net/http"
	"nistagram/agent/service"
	"nistagram/agent/util"

	"github.com/gorilla/mux"
)

type CampaignHandler struct {
	CampaignService *service.CampaignService
}

func (handler *CampaignHandler) SaveCampaignReport(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	campaignId := util.String2Uint(vars["id"])

	err := handler.CampaignService.SaveCampaignReport(campaignId)
	if err != nil {
		fmt.Println(err)
		_, _ = w.Write([]byte("{\"message\":\"error\"}"))
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	_, _ = w.Write([]byte("{\"message\":\"ok\"}"))
	w.Header().Set("Content-Type", "application/json")
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

func (handler *CampaignHandler) EditCampaign(w http.ResponseWriter, r *http.Request) {
	body, err := io.ReadAll(r.Body)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer func(Body io.ReadCloser) {
		_ = Body.Close()
	}(r.Body)
	params := mux.Vars(r)
	resp, err := handler.CampaignService.EditCampaign(params["id"], body)
	if err != nil {
		fmt.Println(err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	w.WriteHeader(http.StatusOK)
	_, _ = w.Write(util.GetResponseJSON(*resp))
	w.Header().Set("Content-Type", "application/json")
}
