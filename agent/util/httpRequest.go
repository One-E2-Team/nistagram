package util

import (
	"bytes"
	"fmt"
	"net/http"
)

var agentJwt = ""

func SetJwt(token string) {
	agentJwt = token
}

func NistagramRequest(method string, path string, data []byte, headers map[string]string) (*http.Response, error) {
	client := &http.Client{}
	req, err := http.NewRequest(method, path, bytes.NewBuffer(data))
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
