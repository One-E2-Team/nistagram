package service

import (
	"context"
	"encoding/json"
	"io"
	"net/http"
	"nistagram/notification/dto"
	"nistagram/util"
)


func (service *Service) GetMessageConnections(loggedUserId uint) ([]dto.MessageConnectionDTO, error) {
	profileIds, err := service.Repository.GetConnectedProfileIds(loggedUserId)
	if err != nil{
		return nil,err
	}

	usernames, err := getProfileUsernamesByIDs(profileIds)
	if err != nil{
		return nil, err
	}

	messageApproved, err := getMessageApprovedByIDs(profileIds)
	if err != nil{
		return nil, err
	}

	var ret []dto.MessageConnectionDTO

	for i, id := range profileIds{
		mcDto := dto.MessageConnectionDTO{ProfileId: id, Username: usernames[i], MessageApproved: messageApproved[i]}
		ret = append(ret, mcDto)
	}

	return ret, nil
}


func getProfileUsernamesByIDs(profileIDs []uint) ([]string, error) {
	type data struct {
		Ids []string `json:"ids"`
	}
	req := make([]string, 0)
	for _, value := range profileIDs {
		req = append(req, util.Uint2String(value))
	}
	bodyData := data{Ids: req}
	jsonBody, err := json.Marshal(bodyData)
	if err != nil {
		return nil, err
	}
	profileHost, profilePort := util.GetProfileHostAndPort()
	resp, err := util.CrossServiceRequest(context.Background(), http.MethodPost,
		util.GetCrossServiceProtocol()+"://"+profileHost+":"+profilePort+"/get-by-ids",
		jsonBody, map[string]string{"Content-Type": "application/json;"})
	if err != nil {
		return nil, err
	}

	var ret []string

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	defer func(Body io.ReadCloser) {
		_ = Body.Close()
	}(resp.Body)

	if err = json.Unmarshal(body, &ret); err != nil {
		return nil, err
	}

	return ret, nil
}

func getMessageApprovedByIDs(profileIDs []uint) ([]bool, error) {
	type data struct {
		Ids []string `json:"ids"`
	}
	req := make([]string, 0)
	for _, value := range profileIDs {
		req = append(req, util.Uint2String(value))
	}
	bodyData := data{Ids: req}
	jsonBody, err := json.Marshal(bodyData)
	if err != nil {
		return nil, err
	}
	profileHost, profilePort := util.GetProfileHostAndPort()
	resp, err := util.CrossServiceRequest(context.Background(), http.MethodPost,
		util.GetCrossServiceProtocol()+"://"+profileHost+":"+profilePort+"/get-by-ids",
		jsonBody, map[string]string{"Content-Type": "application/json;"})
	if err != nil {
		return nil, err
	}

	var retDto []dto.MessageEdgeDTO

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	defer func(Body io.ReadCloser) {
		_ = Body.Close()
	}(resp.Body)

	if err = json.Unmarshal(body, &retDto); err != nil {
		return nil, err
	}

	var ret []bool
	for _, edge := range retDto{
		ret = append(ret, edge.Approved)
	}

	return ret, nil
}