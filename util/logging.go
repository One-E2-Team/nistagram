package util

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"time"
)

type LogType int

const (
	INFO LogType = iota
	WARN
	SUCCESS
	ERROR
)

func (e LogType) ToString() string {
	switch e {
	case INFO:
		return "INFO"
	case WARN:
		return "WARN"
	case SUCCESS:
		return "SUCCESS"
	case ERROR:
		return "ERROR"
	default:
		return fmt.Sprintf("%d", int(e))
	}
}

func GetIPAddress(r *http.Request) string {
	return r.Header.Get("X-FORWARDED-FOR")
}

func GetLoggingStringFromID(id uint) string {
	return "profileId: '" + Uint2String(id) + "'"
}

func Logging(logType LogType, resourceMethod string, resourceIP string, content string, service string) {
	logFIle := "../../logs/" + service + "/"
	switch logType {
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
	oneLog := time.Now().UTC().String() + delimiter + logType.ToString() + delimiter + resourceMethod + delimiter + resourceIP + delimiter + content
	log.Println(oneLog)
}
