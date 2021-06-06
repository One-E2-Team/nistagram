package util

import "os"

func GetAuthHostAndPort() (string, string) {
	var authHost, authPort string = "localhost", "8000" //dev.db environment
	_, ok := os.LookupEnv("DOCKER_ENV_SET_PROD") // dev production environment
	_, ok1 := os.LookupEnv("DOCKER_ENV_SET_DEV") // dev front environment
	if ok || ok1 {
		authHost = "auth"
		authPort = "8080"
	}
	return authHost, authPort
}

func GetConnectionHostAndPort() (string, string) {
	var connHost, connPort string = "localhost", "8085" //dev.db environment
	_, ok := os.LookupEnv("DOCKER_ENV_SET_PROD") // dev production environment
	_, ok1 := os.LookupEnv("DOCKER_ENV_SET_DEV") // dev front environment
	if ok || ok1 {
		connHost = "connection"
		connPort = "8080"
	}
	return connHost, connPort
}

func GetProfileHostAndPort() (string, string) {
	var profileHost, profilePort string = "localhost", "8083" //dev.db environment
	_, ok := os.LookupEnv("DOCKER_ENV_SET_PROD") // dev production environment
	_, ok1 := os.LookupEnv("DOCKER_ENV_SET_DEV") // dev front environment
	if ok || ok1 {
		profileHost = "connection"
		profilePort = "8080"
	}
	return profileHost, profilePort
}