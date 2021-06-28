package util

import (
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
			validPrivileges, ok := GetUserPrivileges(id)
			if !ok {
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
		handleFunc(writer, request)
	}
}
