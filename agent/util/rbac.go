package util

import (
	"net/http"
	"nistagram/agent/service"
)

func RBAC(handler func(http.ResponseWriter, *http.Request), authService *service.AuthService ,privilege string, returnCollection bool) func(http.ResponseWriter, *http.Request) {

	finalHandler := func(pass bool) func(http.ResponseWriter, *http.Request) {
		if pass {
			return handler
		} else {
			return func(writer http.ResponseWriter, request *http.Request) {
				writer.WriteHeader(http.StatusOK)
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
			validPrivileges := authService.GetPrivileges(id)
			valid := false
			for _, val := range *validPrivileges {
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
		handleFunc(writer, request)
	}
}