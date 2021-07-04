package main

import (
	"context"
	"fmt"
	"net/http"
	"nistagram/notification/handler"
	"nistagram/notification/repository"
	"nistagram/notification/service"
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

func initRepo(client *mongo.Client) *repository.Repository {
	return &repository.Repository{Client: client}
}
func initService(repository *repository.Repository) *service.Service {
	return &service.Service{Repository: repository}
}
func initHandler(service *service.Service) *handler.Handler {
	return &handler.Handler{Service: service}
}

func handleFunc(handler *handler.Handler) {
	router := mux.NewRouter().StrictSlash(true)

	router.HandleFunc("/messaging", util.RBAC(handler.MessagingWebSocket, "MESSAGING", false)).Methods("GET")

	router.HandleFunc("/connections", util.RBAC(handler.GetMessageConnections, "MESSAGING", true)).Methods("GET")

	router.HandleFunc("/message/{id}", util.RBAC(handler.DeleteMessage, "MESSAGING", false)).Methods("DELETE")

	router.HandleFunc("/message/{id}", util.RBAC(handler.GetAllMesssages, "MESSAGING", true)).Methods("GET")

	router.HandleFunc("/file", util.RBAC(handler.SaveFile, "MESSAGING", true)).Methods("POST")


	fmt.Println("Starting server..")
	host, port := util.GetNotificationHostAndPort()
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
	util.TracerInit("notification")
	client := initDB()
	defer closeConnection(client)
	repo := initRepo(client)
	service := initService(repo)
	handler := initHandler(service)
	_ = util.SetupMSAuth("notification")
	handleFunc(handler)
}
