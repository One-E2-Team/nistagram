package main

import (
	"fmt"
	"net/http"
	"nistagram/connection/handler"
	"nistagram/connection/repository"
	"nistagram/connection/service"
	"nistagram/util"
	"os"
	"time"

	"github.com/gorilla/mux"
	"github.com/neo4j/neo4j-go-driver/v4/neo4j"
)

func initDB() *neo4j.Driver {

	var (
		driver neo4j.Driver
		err    error
	)
	time.Sleep(10 * time.Second)
	var dbHost, dbPort, dbusername, dbpassword = "localhost", "7687", "neo4j", "root" // dev.db environment
	_, ok := os.LookupEnv("DOCKER_ENV_SET_PROD")                                       // production environment
	if ok {
		dbHost = "graphdb_connection"
		dbPort = "7687"
		dbusername = os.Getenv("DB_USERNAME")
		dbpassword = os.Getenv("DB_PASSWORD")
	} else {
		_, ok := os.LookupEnv("DOCKER_ENV_SET_DEV") // dev front environment
		if ok {
			dbHost = "graphdb_connection"
			dbPort = "7687"
			dbusername = os.Getenv("DB_USERNAME")
			dbpassword = os.Getenv("DB_PASSWORD")
		}
	}
	for {
		driver, err = neo4j.NewDriver("bolt://"+dbHost+":"+dbPort+"/neo4j", neo4j.BasicAuth(dbusername, dbpassword, "Neo4j"))

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

func initConnectionRepo(databaseDriver *neo4j.Driver) *repository.Repository {
	return &repository.Repository{DatabaseDriver: databaseDriver}
}

func initService(connectionRepo *repository.Repository) *service.Service {
	return &service.Service{ConnectionRepository: connectionRepo}
}

func initHandler(service *service.Service) *handler.Handler {
	return &handler.Handler{ConnectionService: service}
}

func handleFunc(handler *handler.Handler) {
	router := mux.NewRouter().StrictSlash(true)
	router.HandleFunc("/profile/{id}", util.MSAuth(handler.AddProfile, []string{"profile"})).Methods("POST")

	router.HandleFunc("/profile/{id}", util.MSAuth(handler.DeleteProfile, []string{"profile"})).Methods("DELETE")

	router.HandleFunc("/profile/{id}", util.MSAuth(handler.ReActivateProfile, []string{"profile"})).Methods("PATCH")

	router.HandleFunc("/connection/block/relationships/{id}", util.MSAuth(handler.GetBlockingRelationships, []string{"post"})).Methods("GET")

	router.HandleFunc("/connection/following/all/{id}", handler.GetFollowedProfiles).Methods("GET")

	router.HandleFunc("/connection/following/show/{id}",
		util.MSAuth(handler.GetFollowedProfilesNotMuted, []string{"post", "profile", "postreaction"})).Methods("GET")

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

	router.HandleFunc("/connection/block/{profileId}",
		util.RBAC(handler.ToggleBlockProfile, "EDIT_CONNECTION_STATUS", false)).Methods("PUT") //frontend func

	router.HandleFunc("/connection/closeFriend/{profileId}",
		util.RBAC(handler.ToggleCloseFriendProfile, "EDIT_CONNECTION_STATUS", false)).Methods("PUT") //frontend func

	router.HandleFunc("/connection/mute/{profileId}",
		util.RBAC(handler.ToggleMuteProfile, "EDIT_CONNECTION_STATUS", false)).Methods("PUT") //frontend func

	router.HandleFunc("/connection/notify/post/{profileId}",
		util.RBAC(handler.ToggleNotifyPostProfile, "EDIT_CONNECTION_STATUS", false)).Methods("PUT") //frontend func

	router.HandleFunc("/connection/notify/story/{profileId}",
		util.RBAC(handler.ToggleNotifyStoryProfile, "EDIT_CONNECTION_STATUS", false)).Methods("PUT") //frontend func

	router.HandleFunc("/connection/notify/comment/{profileId}",
		util.RBAC(handler.ToggleNotifyCommentProfile, "EDIT_CONNECTION_STATUS", false)).Methods("PUT") //frontend func

	router.HandleFunc("/connection/block/{profileId}",
		util.RBAC(handler.IsBlocked, "READ_CONNECTION_STATUS", false)).Methods("GET") //frontend func

	router.HandleFunc("/connection/block/am/{profileId}",
		util.RBAC(handler.AmBlocked, "READ_CONNECTION_STATUS", false)).Methods("GET") //frontend func

	router.HandleFunc("/connection/unfollow/{profileId}",
		util.RBAC(handler.UnfollowProfile, "EDIT_CONNECTION_STATUS", false)).Methods("PUT") //frontend func

	router.HandleFunc("/connection/messaging/connect/{profileId}",
		util.RBAC(handler.MessageConnect, "EDIT_CONNECTION_STATUS", false)).Methods("POST") // frontend func

	router.HandleFunc("/connection/messaging/request/{profileId}",
		util.RBAC(handler.MessageRequest, "CREATE_CONNECTION", false)).Methods("POST") //frontend func

	router.HandleFunc("/connection/messaging/my-properties/{profileId}",
		util.RBAC(handler.GetMessageRelationship, "READ_CONNECTION_STATUS", false)).Methods("GET") // frontend func

	router.HandleFunc("/connection/notify/message/{profileId}",
		util.RBAC(handler.ToggleNotifyMessageProfile, "EDIT_CONNECTION_STATUS", false)).Methods("PUT") //frontend func

	router.HandleFunc("/connection/messaging/decline/{profileId}",
		util.RBAC(handler.DeclineMessageRequest, "EDIT_CONNECTION_STATUS", false)).Methods("DELETE") //frontend func

	router.HandleFunc("/connection/messaging/request",
		util.RBAC(handler.GetAllMessageRequests, "READ_CONNECTION_REQUESTS", true)).Methods("GET") //frontend func

	router.HandleFunc("/connection/recommendation",
		util.RBAC(handler.RecommendProfiles, "READ_CONNECTION_STATUS", false)).Methods("GET") // frontend func

	fmt.Println("Starting server..")
	host, port := util.GetConnectionHostAndPort()
	var err error
	if util.DockerChecker() {
		err = http.ListenAndServeTLS(":"+port, "../cert.pem", "../key.pem", router)
	} else {
		err = http.ListenAndServe(host+":"+port, router)
	}
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
	connectionHandler := initHandler(connectionService)
	_ = util.SetupMSAuth("connection")
	handleFunc(connectionHandler)
}
