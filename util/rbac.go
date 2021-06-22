package util

import (
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"reflect"
	"runtime"
	"strings"
)

func RBAC(handler func(http.ResponseWriter, *http.Request), privilege string, returnCollection bool) func(http.ResponseWriter, *http.Request) {

	finalHandler := func(pass bool) func(http.ResponseWriter, *http.Request) {
		if pass {
			return handler
		} else {
			return func(writer http.ResponseWriter, request *http.Request) {
				writer.WriteHeader(http.StatusOK)
				handlerFunctionName := runtime.FuncForPC(reflect.ValueOf(handler).Pointer()).Name()
				parts := strings.Split(handlerFunctionName, "/")
				Logging(WARN, handlerFunctionName, GetIPAddress(request), "Unauthorized access", parts[1])
				writer.Header().Set("Content-Type", "application/json")
				if returnCollection {
					_, _ = writer.Write([]byte("[{\"status\":\"fail\", \"reason\":\"unauthorized\"}]"))
				} else {
					_, _ = writer.Write([]byte("{\"status\":\"fail\", \"reason\":\"unauthorized\"}"))
				}
			}
		}
	}

	return func(writer http.ResponseWriter, request *http.Request) {
		var handleFunc func(http.ResponseWriter, *http.Request)
		id := GetLoggedUserIDFromToken(request)
		if id == 0 {
			handleFunc = finalHandler(false)
		} else {
			authHost, authPort := GetAuthHostAndPort()
			resp, err := CrossServiceRequest(http.MethodGet,
				GetCrossServiceProtocol()+"://"+authHost+":"+authPort+"/privileges/"+Uint2String(id),
				nil, map[string]string{})
			if err != nil {
				fmt.Println(err)
				handleFunc = finalHandler(false)
			} else {
				defer func(Body io.ReadCloser) {
					err := Body.Close()
					if err != nil {

					}
				}(resp.Body)
				body, err1 := ioutil.ReadAll(resp.Body)
				if err1 != nil {
					fmt.Println(err1)
					handleFunc = finalHandler(false)
				} else {
					var validPrivileges []string
					err = json.Unmarshal(body, &validPrivileges)
					if err != nil {
						fmt.Println(err)
						handleFunc = finalHandler(false)
					} else {
						valid := false
						for _, val := range validPrivileges {
							if val == privilege {
								valid = true
								break
							}
						}
						if valid {
							handleFunc = finalHandler(true)
						} else {
							handleFunc = finalHandler(false)
						}
					}
				}
			}
		}
		handleFunc(writer, request)
	}
}
