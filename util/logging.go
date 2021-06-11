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

const LogSize = 104857600

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

func getLogFileString(logType LogType) string{
	switch logType {
	case INFO:
		return "infoLogs"
	case WARN:
		return "warnLogs"
	case SUCCESS:
		return "successLogs"
	case ERROR:
		return "errorLogs"
	}
	return ""
}

func Logging(logType LogType, resourceMethod string, resourceIP string, content string, service string) {
	logFileService := "../../logs/" + service + "/"
	logFile := logFileService
	logFile += getLogFileString(logType)
	logFile += ".txt"

	file, err := os.OpenFile(logFile, os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0666)
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

	stat, err := file.Stat()

	if(stat.Size() > 1024){ //should use LogSize here
		os.Rename(logFile, logFileService + getLogFileString(logType) + time.Now().UTC().String() + ".txt")
	}
}
