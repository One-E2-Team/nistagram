package handler

import (
	"fmt"
	"github.com/gorilla/mux"
	"net/http"
	"nistagram/agent/service"
	"nistagram/agent/util"
)

type CampaignHandler struct {
	CampaignService *service.CampaignService
}

func (handler *CampaignHandler) SaveCampaignReport(w http.ResponseWriter, r *http.Request){
	vars := mux.Vars(r)
	campaignId := util.String2Uint(vars["id"])

	err := handler.CampaignService.SaveCampaignReport(campaignId)
	if err != nil{
		fmt.Println(err)
		_, _ = w.Write([]byte("{\"message\":\"error\"}"))
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	_, _ = w.Write([]byte("{\"message\":\"ok\"}"))
	w.Header().Set("Content-Type", "application/json")
}