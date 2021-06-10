package main

import (
	"fmt"
	"github.com/gorilla/mux"
	"github.com/neo4j/neo4j-go-driver/v4/neo4j"
	"net/http"
	"nistagram/connection/handler"
	"nistagram/connection/repository"
	"nistagram/connection/service"
	"nistagram/util"
	"os"
	"time"
)

func initDB() *neo4j.Driver {

	var (
		driver neo4j.Driver
		err    error
	)
	time.Sleep(10 * time.Second)
	var dbHost, dbPort = /*, dbusername, dbpassword*/ "localhost", "7687" //, "root", "root" // dev.db environment
	_, ok := os.LookupEnv("DOCKER_ENV_SET_PROD")                          // production environment
	if ok {
		dbHost = "graphdb_connection"
		dbPort = "7687"
		//dbusername = os.Getenv("DB_USERNAME")
		//dbpassword = os.Getenv("DB_PASSWORD")
	} else {
		_, ok := os.LookupEnv("DOCKER_ENV_SET_DEV") // dev front environment
		if ok {
			dbHost = "graphdb_connection"
			dbPort = "7687"
			//dbusername = "root"
			//dbpassword = "root"
		}
	}
	for {
		driver, err = neo4j.NewDriver("bolt://"+dbHost+":"+dbPort+"/neo4j", neo4j.NoAuth())

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
	router.HandleFunc("/connection/following/all/{id}", handler.GetFollowedProfiles).Methods("GET")
	router.HandleFunc("/connection/following/show/{id}", handler.GetFollowedProfilesNotMuted).Methods("GET")
	router.HandleFunc("/connection/following/properties/{followerId}/{profileId}", handler.GetConnection).Methods("GET")
	router.HandleFunc("/connection/following/update", handler.UpdateConnection).Methods("PUT") //frontend func
	router.HandleFunc("/connection/following/my-properties/{profileId}",
		util.RBAC(handler.GetConnectionPublic, "READ_CONNECTION_STATUS", false)).Methods("GET") // frontend func
	router.HandleFunc("/connection/following/approve/{profileId}",
		util.RBAC(handler.FollowApprove, "EDIT_CONNECTION_STATUS", false)).Methods("POST") // frontend func
	router.HandleFunc("/connection/following/request/{profileId}",
		util.RBAC(handler.FollowRequest, "CREATE_CONNECTION", false)).Methods("POST") //frontend func
	router.HandleFunc("/connection/following/request",
		util.RBAC(handler.GetAllFollowRequests, "READ_CONNECTION_REQUESTS", true)).Methods("GET") //frontend func
	router.HandleFunc("/connection/following/request/{profileId}",
		util.RBAC(handler.DeclineFollowRequest, "EDIT_CONNECTION_STATUS", false)).Methods("DELETE") //frontend func
	fmt.Println("Starting server..")
	_, port := util.GetConnectionHostAndPort()
	err := http.ListenAndServe(":"+port, router)
	if err != nil {
		fmt.Println(err)
		return
	}
}

func main() {
	db := initDB()
	defer (*db).Close()
	connectionRepo := initConnectionRepo(db)
	connectionService := initService(connectionRepo)
	handler := initHandler(connectionService)
	handleFunc(handler)
}
