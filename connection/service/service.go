package service

import (
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"nistagram/connection/repository"
	model2 "nistagram/profile/model"
	"nistagram/util"
)

type Service struct {
	ConnectionRepository *repository.Repository
}

func contains(s *[]uint, e uint) bool {
	for _, a := range *s {
		if a == e {
			return true
		}
	}
	return false
}

func getProfile(id uint) *model2.Profile {
	var p model2.Profile
	profileHost, profilePort := util.GetProfileHostAndPort()
	resp, err := util.CrossServiceRequest(http.MethodGet,
		util.CrossServiceProtocol+"://"+profileHost+":"+profilePort+"/get-by-id/"+util.Uint2String(id),
		nil, map[string]string{})
	if err != nil {
		fmt.Println(err)
		return nil
	}
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {

		}
	}(resp.Body)
	body, err1 := ioutil.ReadAll(resp.Body)
	if err1 != nil {
		fmt.Println(err1)
		return nil
	}
	err = json.Unmarshal(body, &p)
	if err != nil {
		fmt.Println(err)
		return nil
	}
	return &p
}