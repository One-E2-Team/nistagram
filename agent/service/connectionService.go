package service

import (
	"net/http"
	"nistagram/agent/util"
)

type ConnectionService struct {

}

func (service *ConnectionService) GetMyFollowedProfiles() (*http.Response, error) {
	return util.NistagramRequest(http.MethodGet, "/agent-api/connection/following/my/all-agent",
		nil, map[string]string{})
}
