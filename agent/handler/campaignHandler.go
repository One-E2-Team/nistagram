package handler

import (
	"fmt"
	"net/http"
	"nistagram/agent/service"
	"nistagram/agent/util"
)

type CampaignHandler struct {
	CampaignService *service.CampaignService
}

func (handler *CampaignHandler) GetMyCampaigns(w http.ResponseWriter, r *http.Request) {
	resp, err := handler.CampaignService.GetMyCampaigns(util.GetLoggedUserIDFromToken(r))
	if err != nil {
		fmt.Println(err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	w.WriteHeader(http.StatusOK)
	_, _ = w.Write(util.GetResponseJSON(*resp))
	w.Header().Set("Content-Type", "application/json")
}
