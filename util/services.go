package util

import "os"

const MicroservicesProtocol = "https"
const CrossServiceProtocol = "http"
const FrontProtocol = "http"

func dockerChecker() bool {
	_, ok := os.LookupEnv("DOCKER_ENV_SET_PROD") // dev production environment
	_, ok1 := os.LookupEnv("DOCKER_ENV_SET_DEV") // dev front environment
	return ok || ok1
}

func GetAuthHostAndPort() (string, string) {
	var authHost, authPort = "localhost", "8000" // dev.db environment
	if dockerChecker() {
		authHost = "auth"
		authPort = "8080"
	}
	return authHost, authPort
}

func GetConnectionHostAndPort() (string, string) {
	var connHost, connPort = "localhost", "8085" // dev.db environment
	if dockerChecker() {
		connHost = "connection"
		connPort = "8080"
	}
	return connHost, connPort
}

func GetProfileHostAndPort() (string, string) {
	var profileHost, profilePort = "localhost", "8083" // dev.db environment
	if dockerChecker() {
		profileHost = "profile"
		profilePort = "8080"
	}
	return profileHost, profilePort
}

func GetPostHostAndPort() (string, string) {
	var postHost, postPort = "localhost", "8086" // dev.db environment
	if dockerChecker() {
		postHost = "post"
		postPort = "8080"
	}
	return postHost, postPort
}

func GetFrontHostAndPort() (string, string) {
	var frontHost, frontPort = "localhost", "3000" //dev environment
	_, ok := os.LookupEnv("DOCKER_ENV_SET_PROD")   // dev production environment
	if ok {
		frontPort = "81"
	}
	return frontHost, frontPort
}

func GetGatewayHostAndPort() (string, string) {
	var gatewayHost, gatewayPort = "localhost", "81"
	return gatewayHost, gatewayPort
}
