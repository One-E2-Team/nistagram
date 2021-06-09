package util

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
)

func RBAC(handler func(http.ResponseWriter, *http.Request), privilege string, returnCollection bool) func(http.ResponseWriter, *http.Request) {

	finalHandler := func(pass bool) func(http.ResponseWriter, *http.Request) {
		if pass {
			return handler
		} else {
			return func(writer http.ResponseWriter, request *http.Request) {
				writer.WriteHeader(http.StatusOK)
				writer.Header().Set("Content-Type", "application/json")
				if returnCollection{
					writer.Write([]byte("[{\"status\":\"fail\", \"reason\":\"unauthorized\"}]"))
				} else {
					writer.Write([]byte("{\"status\":\"fail\", \"reason\":\"unauthorized\"}"))
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
			resp, err := http.Get("http://" + authHost + ":" + authPort + "/privileges/" + Uint2String(id))
			if err != nil {
				fmt.Println(err)
				handleFunc = finalHandler(false)
			} else {
				defer resp.Body.Close()
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
