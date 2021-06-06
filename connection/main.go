package main

import (
	"fmt"
	"github.com/gorilla/mux"
	"github.com/neo4j/neo4j-go-driver/v4/neo4j"
	"net/http"
	"nistagram/connection/handler"
	"nistagram/connection/repository"
	"nistagram/connection/service"
	"os"
	"time"
)

func initDB() *neo4j.Driver {

	var (driver neo4j.Driver; err error)
	time.Sleep(10 * time.Second)
	var dbhost, dbport/*, dbusername, dbpassword*/ string = "localhost", "7687"//, "root", "root" // dev.db environment
	_, ok := os.LookupEnv("DOCKER_ENV_SET_PROD") // production environment
	if ok {
		dbhost = "graphdb_connection"
		dbport = "7687"
		//dbusername = os.Getenv("DB_USERNAME")
		//dbpassword = os.Getenv("DB_PASSWORD")
	} else {
		_, ok := os.LookupEnv("DOCKER_ENV_SET_DEV") // dev front environment
		if ok {
			dbhost = "graphdb_connection"
			dbport = "7687"
			//dbusername = "root"
			//dbpassword = "root"
		}
	}
	for {
		driver, err = neo4j.NewDriver("bolt://" + dbhost + ":" + dbport + "/neo4j", neo4j.NoAuth())

		if err != nil {
			fmt.Println("Cannot connect to database! Sleeping 10s and then retrying....")
			time.Sleep(10 * time.Second)
		} else {
			fmt.Println("Connected to the database.")
			break
		}
	}

	return &driver
}

func initConnectionRepo(databaseDriver *neo4j.Driver) *repository.ConnectionRepository {
	return &repository.ConnectionRepository{DatabaseDriver: databaseDriver}
}

func initService(connectionRepo *repository.ConnectionRepository) *service.ConnectionService {
	return &service.ConnectionService{ConnectionRepository: connectionRepo}
}

func initHandler(service *service.ConnectionService) *handler.Handler {
	return &handler.Handler{ConnectionService: service}
}

func handleFunc(handler *handler.Handler) {
	router := mux.NewRouter().StrictSlash(true)
	router.HandleFunc("/profile/{id}", handler.AddProfile).Methods("POST")
	router.HandleFunc("/connection/{followerId}/{profileId}", handler.FollowRequest).Methods("POST")
	router.HandleFunc("/connection/{followerId}/{profileId}", handler.GetConnection).Methods("GET")
	router.HandleFunc("/connection/following/all/{id}", handler.GetFollowedProfiles).Methods("GET")
	router.HandleFunc("/connection/following/show/{id}", handler.GetFollowedProfilesNotMuted).Methods("GET")
	//router.HandleFunc("/{username}", handler.GetProfileByUsername).Methods("GET")
	fmt.Printf("Starting server..")
	var port string = "8085" // dev.db environ
	_, ok := os.LookupEnv("DOCKER_ENV_SET_PROD") // dev production environment
	_, ok1 := os.LookupEnv("DOCKER_ENV_SET_DEV") // dev front environment
	if ok || ok1 {
		port = "8080"
	}
	http.ListenAndServe(":" + port, router)
}

func main() {
	db := initDB()
	defer (*db).Close()
	connectionRepo := initConnectionRepo(db)
	connectionService := initService(connectionRepo)
	handler := initHandler(connectionService)
	handleFunc(handler)
}