package util

import "os"

const FrontProtocol = "https"

func DockerChecker() bool {
	_, ok := os.LookupEnv("DOCKER_ENV_SET_PROD") // dev production environment
	_, ok1 := os.LookupEnv("DOCKER_ENV_SET_DEV") // dev front environment
	return ok || ok1
}

func GetCrossServiceProtocol() string {
	if DockerChecker(){
		return "https"
	}
	return "http"
}

func GetMicroservicesProtocol() string {
	if DockerChecker(){
		return "https"
	}
	return "http"
}

func GetFrontProtocol() string {
	if DockerChecker(){
		return "https"
	}
	return "http"
}

func GetAuthHostAndPort() (string, string) {
	var authHost, authPort = "localhost", "8000" // dev.db environment
	if DockerChecker() {
		authHost = "auth"
		authPort = "8080"
	}
	return authHost, authPort
}

func GetConnectionHostAndPort() (string, string) {
	var connHost, connPort = "localhost", "8085" // dev.db environment
	if DockerChecker() {
		connHost = "connection"
		connPort = "8080"
	}
	return connHost, connPort
}

func GetProfileHostAndPort() (string, string) {
	var profileHost, profilePort = "localhost", "8083" // dev.db environment
	if DockerChecker() {
		profileHost = "profile"
		profilePort = "8080"
	}
	return profileHost, profilePort
}

func GetPostHostAndPort() (string, string) {
	var postHost, postPort = "localhost", "8086" // dev.db environment
	if DockerChecker() {
		postHost = "post"
		postPort = "8080"
	}
	return postHost, postPort
}

func GetPostReactionHostAndPort() (string, string) {
	var postReactionHost, postReactionPort = "localhost", "8087" // dev.db environment
	if DockerChecker() {
		postReactionHost = "postreaction"
		postReactionPort = "8080"
	}
	return postReactionHost, postReactionPort
}

func GetCampaignHostAndPort() (string, string) {
	var campaignHost, campaignPort = "localhost", "8088" // dev.db environment
	if DockerChecker() {
		campaignHost = "campaign"
		campaignPort = "8080"
	}
	return campaignHost, campaignPort
}

func GetFrontHostAndPort() (string, string) {
	var frontHost, frontPort = "localhost", "3000"
	if DockerChecker() {
		frontPort = "81"
	}
	return frontHost, frontPort
}

func GetGatewayHostAndPort() (string, string) {
	var gatewayHost, gatewayPort = "localhost", "81"
	return gatewayHost, gatewayPort
}

func GetMonitoringHostAndPort() (string, string) {
	var monitoringHost, monitoringPort = "localhost", "8089" // dev.db environment
	if DockerChecker() {
		monitoringHost = "monitoring"
		monitoringPort = "8080"
	}
	return monitoringHost, monitoringPort
}

func GetNotificationHostAndPort() (string, string) {
	var notificationHost, notificationPort = "localhost", "8090" // dev.db environment
	if DockerChecker() {
		notificationHost = "notification"
		notificationPort = "8080"
	}
	return notificationHost, notificationPort
}
