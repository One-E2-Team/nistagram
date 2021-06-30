package util

import (
	"fmt"
	"net/smtp"
	"os"
	"strconv"
)

func GetAgentHostAndPort() (string, string) {
	var agentHost, agentPort = "localhost", "82"
	return agentHost, agentPort
}

func GetAgentProtocol() string {
	if DockerChecker(){
		return "https"
	}
	return "http"
}

func GetNistagramHostAndPort() (string, string) {
	var nistagramHost, nistagramPort = "localhost", "81"
	if DockerChecker() {
		nistagramHost = "apigateway"
		nistagramPort = "80"
	}
	return nistagramHost, nistagramPort
}

func GetExistDBHostAndPort() (string, string) {
	var existHost, existPort = "localhost", "8666"
	if DockerChecker() {
		existHost = "exist"
		existPort = "8080"
	}
	return existHost, existPort
}

func GetNistagramProtocol() string {
	if DockerChecker(){
		return "https"
	}
	return "http"
}

func GetExistDBProtocol() string {
	return "http"
}

func DockerChecker() bool {
	_, ok := os.LookupEnv("DOCKER_ENV_SET_PROD") // dev production environment
	_, ok1 := os.LookupEnv("DOCKER_ENV_SET_DEV") // dev front environment
	return ok || ok1
}

func SendMail(sendTo string, subject string, mailMessage string) {
	from := os.Getenv("MAIL_USERNAME")
	password := os.Getenv("MAIL_PASSWORD")
	to := []string{sendTo}
	smtpHost := "smtp.gmail.com"
	smtpPort := "587"
	msg := []byte("To: " + sendTo + "\r\n" + "Subject: " + subject + "\r\n" + "\r\n" + mailMessage + "\r\n")
	auth := smtp.PlainAuth("", from, password, smtpHost)
	err := smtp.SendMail(smtpHost+":"+smtpPort, auth, from, to, msg)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println("Email Sent Successfully!")
}

func Uint2String(input uint) string {
	return strconv.FormatUint(uint64(input), 10)
}

func String2Uint(input string) uint {
	u64, err := strconv.ParseUint(input, 10, 32)
	if err != nil {
		fmt.Println(err)
		return 0
	}
	return uint(u64)
}

func Contains(array []uint, el uint) bool {
	for _, a := range array {
		if a == el {
			return true
		}
	}
	return false
}

func GetFrontProtocol() string {
	if DockerChecker(){
		return "https"
	}
	return "http"
}

func GetFrontHostAndPort() (string, string) {
	var frontHost, frontPort = "localhost", "3001"
	if DockerChecker() {
		frontPort = "82"
	}
	return frontHost, frontPort
}
