package util

import (
	"bytes"
	"fmt"
	"io"
	"net/http"
	"net/url"
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

func ExistDBRequest(method string, urlPath string, data []byte, headers map[string]string) (*http.Response, error) {
	client := &http.Client{}
	existHost, existPort := GetExistDBHostAndPort()
	path := GetExistDBProtocol() + "://" + existHost + ":" + existPort
	req, err := http.NewRequest(method, path + urlPath, bytes.NewBuffer(data))
	if err != nil {
		fmt.Println(err)
		return nil, err
	}
	req.SetBasicAuth("admin", "")
	req.Header.Set("Content-Type", "application/xml")
	req.Header.Set("Content-Length", fmt.Sprintf("%d", len(data)))
	for key, value := range headers {
		req.Header.Set(key, value)
	}
	return client.Do(req)
}

func GenerateFuckingExistDBHTTPRequestParametersQuery(data map[string]string) string{
	ret := "?"
	first := true
	for key, value := range data {
		if !first {
			ret += "&"
		} else {
			first = false
		}
		ret += url.QueryEscape(key) + "=" + url.QueryEscape(value)
	}
	if first {
		return ""
	} else {
		return ret
	}
}

func GetResponseJSON(response http.Response) []byte{
	body, err := io.ReadAll(response.Body)
	if err != nil {
		fmt.Println(err)
		return make([]byte, 0)
	}
	defer func(Body io.ReadCloser) {
		_ = Body.Close()
	}(response.Body)
	return body
}
