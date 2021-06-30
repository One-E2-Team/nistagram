package util

import (
	"bytes"
	"fmt"
	"net/http"
	"strings"
)

var agentJwt = ""

func SetJwt(token string) {
	agentJwt = strings.ReplaceAll(token, "\"", "")
}

func NistagramRequest(method string, urlPath string, data []byte, headers map[string]string) (*http.Response, error) {
	client := &http.Client{}
	nistagramHost, nistagramPort := GetNistagramHostAndPort()
	path := GetNistagramProtocol() + "://" + nistagramHost + ":" + nistagramPort
	req, err := http.NewRequest(method, path + urlPath, bytes.NewBuffer(data))
	if err != nil {
		fmt.Println(err)
		return nil, err
	}
	req.Header.Set("Authorization", "Bearer "+agentJwt)
	for key, value := range headers {
		req.Header.Set(key, value)
	}
	return client.Do(req)
}
