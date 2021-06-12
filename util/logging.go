package util

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"strings"
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

func getLogFileString(logType LogType) string {
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

func encodeInputString(input string) string {
	input = strings.ReplaceAll(input, "|", "%file_separator%")
	return input
}

func Logging(logType LogType, resourceMethod string, resourceIP string, content string, service string) {
	logFileService := "../../logs/" + service + "/"
	logFile := logFileService
	logFile += getLogFileString(logType)
	logFile += ".txt"

	file, err := os.OpenFile(logFile, os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0666)

	if err != nil {
		fmt.Println(err)
	}
	log.SetOutput(file)
	delimiter := "|"
	oneLog := delimiter + time.Now().UTC().String() + delimiter + logType.ToString() + delimiter +
		resourceMethod + delimiter + encodeInputString(resourceIP) + delimiter +
		encodeInputString(content) + delimiter
	log.Println(oneLog)

	stat, err := file.Stat()

	if stat.Size() > 1024 { //TODO: should use LogSize here
		err = file.Close()
		if err != nil {
			return
		}
		newFile, err := os.OpenFile(logFileService+"/archive/"+getLogFileString(logType)+
			time.Now().UTC().String()+".txt", os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0666)
		if err != nil {
			fmt.Println(err)
			return
		}
		defer func(newFile *os.File) {
			err := newFile.Close()
			if err != nil {

			}
		}(newFile)

		file, err = os.OpenFile(logFile, os.O_RDONLY, 0444)
		_, err = io.Copy(newFile, file)
		if err != nil {
			fmt.Println(err)
			return
		}
		if err := os.Truncate(file.Name(), 0); err != nil {
			fmt.Printf("Failed to truncate: %v", err)
		}
	}
	err = file.Close()
	if err != nil {
		return
	}
}
