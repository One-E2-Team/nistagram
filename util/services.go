package util

import "os"

func GetAuthHostAndPort() (string, string) {
	var authHost, authPort string = "localhost", "8000"
	_, ok := os.LookupEnv("DOCKER_ENV_SET_PROD") // dev production environment
	_, ok1 := os.LookupEnv("DOCKER_ENV_SET_DEV") // dev front environment
	if ok || ok1 {
		authHost = "auth"
		authPort = "8080"
	}
	return authHost, authPort
}

func GetConnectionHostAndPort() (string, string) {
	var connHost, connPort string = "localhost", "8085"
	_, ok := os.LookupEnv("DOCKER_ENV_SET_PROD") // dev production environment
	_, ok1 := os.LookupEnv("DOCKER_ENV_SET_DEV") // dev front environment
	if ok || ok1 {
		connHost = "connection"
		connPort = "8080"
	}
	return connHost, connPort
}