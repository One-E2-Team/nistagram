package main

import (
	"context"
	"fmt"
	"net/http"
	"nistagram/monitoring/handler"
	"nistagram/monitoring/repository"
	"nistagram/monitoring/service"
	"nistagram/util"
	"os"
	"time"

	"github.com/gorilla/mux"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func initDB() *mongo.Client {
	var dbHost, dbPort, dbUsername, dbPassword = "localhost", "8084", "root", "root" // dev.db environment
	_, ok := os.LookupEnv("DOCKER_ENV_SET_PROD")                                     // production environment
	if ok {
		dbHost = "mongo1"
		dbPort = "27017"
		dbUsername = os.Getenv("DB_USERNAME")
		dbPassword = os.Getenv("DB_PASSWORD")
	} else {
		_, ok := os.LookupEnv("DOCKER_ENV_SET_DEV") // dev front environment
		if ok {
			dbHost = "mongo1"
			dbPort = "27017"
			dbUsername = os.Getenv("DB_USERNAME")
			dbPassword = os.Getenv("DB_PASSWORD")
		}
	}

	clientOptions := options.Client().ApplyURI("mongodb://" + dbUsername + ":" + dbPassword + "@" + dbHost + ":" + dbPort)
	for {
		client, err := mongo.Connect(context.TODO(), clientOptions)

		if err != nil {
			fmt.Println("Cannot connect to MongoDB! Sleeping 10s and then retrying....")
			time.Sleep(10 * time.Second)
		} else {
			fmt.Println("Connected to MongoDB")
			return client
		}
	}
}

func initRepo(client *mongo.Client) *repository.MonitoringRepository {
	return &repository.MonitoringRepository{Client: client}
}

func initService(monitoringRepo *repository.MonitoringRepository) *service.MonitoringService {
	return &service.MonitoringService{MonitoringRepository: monitoringRepo}
}
func initHandler(monitoringService *service.MonitoringService) *handler.Handler {
	return &handler.Handler{MonitoringService: monitoringService}
}

func handleFunc(handler *handler.Handler) {
	router := mux.NewRouter().StrictSlash(true)
	util.InitMonitoring("monitoring", router)

	router.HandleFunc("/influencer",
		util.MSAuth(handler.CreateEventInfluencer, []string{"postreaction"})).Methods("POST")
	router.HandleFunc("/target-group",
		util.MSAuth(handler.CreateEventTargetGroup, []string{"postreaction"})).Methods("POST")
	router.HandleFunc("/statistics/{campaignId}", util.AgentAuth(handler.GetCampaignStatistics)).Methods("GET")
	router.HandleFunc("/redirect/{campaignId}/{influencerId}/{mediaId}", handler.VisitSite).Methods("GET")

	fmt.Println("Starting server..")
	host, port := util.GetMonitoringHostAndPort()
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

func closeConnection(client *mongo.Client) {
	err := client.Disconnect(context.TODO())
	if err != nil {
		fmt.Println("Failed to close MongoDB.")
		return
	}
	fmt.Println("Connection to MongoDB closed.")
}

func main() {
	util.TracerInit("monitoring")
	client := initDB()
	defer closeConnection(client)
	repo := initRepo(client)
	service := initService(repo)
	handler := initHandler(service)
	_ = util.SetupMSAuth("monitoring")
	handleFunc(handler)
}
