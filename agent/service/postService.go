package service

import (
	"net/http"
	"nistagram/agent/util"
)

type PostService struct {

}

func (service *PostService) GetMyPosts() (*http.Response, error){
	return util.NistagramRequest(http.MethodGet, "/agent-api/post/my", nil, map[string]string{})
}
