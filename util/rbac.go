package util

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"reflect"
	"runtime"
	"strings"
	"time"
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
				userA := request.Header.Get("X-FORWARDED-FOR")
				//userA := request.RemoteAddr
				Logging(INFO, "Unauthorized access from " + userA + " to '"+request.Method+":"+request.RequestURI+"' in "+handlerFunctionName+".", parts[1])
				writer.Header().Set("Content-Type", "application/json")
				if returnCollection {
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
			resp, err := http.Get(CrossServiceProtocol + "://" + authHost + ":" + authPort + "/privileges/" + Uint2String(id))
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

func Logging(logType LogType, reason string, service string) {
	logFIle := "../../logs/" + service + "/"
	switch logType{
	case INFO:
		logFIle += "infoLogs"
	case WARN:
		logFIle += "warnLogs"
	case SUCCESS:
		logFIle += "successLogs"
	case ERROR:
		logFIle += "errorLogs"
	}
	logFIle += ".txt"
	file, err := os.OpenFile(logFIle, os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0666)
	defer func(file *os.File) {
		err := file.Close()
		if err != nil {

		}
	}(file)
	if err != nil {
		fmt.Println(err)
	}
	log.SetOutput(file)
	delimiter := "|"
	oneLog := time.Now().UTC().String() + delimiter + logType.ToString() + delimiter + reason
	log.Println(oneLog)
}
