package main

import (
	"context"
	"fmt"
	"github.com/gorilla/mux"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"net/http"
	"nistagram/postreaction/handler"
	"nistagram/postreaction/repository"
	"nistagram/postreaction/service"
	"nistagram/util"
	"os"
	"time"
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

func initPostRepo(client *mongo.Client) *repository.PostReactionRepository {
	return &repository.PostReactionRepository{Client: client}
}

func initService(postRepo *repository.PostReactionRepository) *service.PostReactionService {
	return &service.PostReactionService{PostReactionRepository: postRepo}
}
func initHandler(postService *service.PostReactionService) *handler.PostReactionHandler {
	return &handler.PostReactionHandler{PostReactionService: postService}
}

func handleFunc(handler *handler.PostReactionHandler) {
	router := mux.NewRouter().StrictSlash(true)
	router.HandleFunc("/react",
		util.RBAC(handler.ReactOnPost, "REACT_ON_POST", false)).Methods("POST") //frontend func
	router.HandleFunc("/react/{postID}",
		util.RBAC(handler.DeleteReaction, "REACT_ON_POST", false)).Methods("DELETE") //frontend func
	router.HandleFunc("/comment",
		util.RBAC(handler.CommentPost, "REACT_ON_POST", false)).Methods("POST") //frontend func
	router.HandleFunc("/report",
		util.RBAC(handler.ReportPost, "REPORT_POST", false)).Methods("POST") //frontend func
	router.HandleFunc("/my-reactions/{type}",
		util.RBAC(handler.GetMyReactions, "READ_REACTIONS", true)).Methods("GET") //frontend func
	router.HandleFunc("/all-reactions/{postID}", handler.GetAllReactions).Methods("GET") //frontend func
	router.HandleFunc("/get-reaction-types/{profileID}",
		util.MSAuth(handler.GetReactionTypes, []string{"post"})).Methods("POST")
	router.HandleFunc("/report",
		util.RBAC(handler.GetAllReports, "READ_REPORTS", true)).Methods("GET")	  //frontend func
	router.HandleFunc("/report/{postId}",
		util.MSAuth(handler.DeletePostsReports, []string{"post"})).Methods("DELETE")
	fmt.Println("Post reaction server started...")
	host, port := util.GetPostReactionHostAndPort()
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
	client := initDB()
	defer closeConnection(client)
	postReactionRepo := initPostRepo(client)
	postReactionService := initService(postReactionRepo)
	postReactionHandler := initHandler(postReactionService)
	_ = util.SetupMSAuth("postreaction")
	handleFunc(postReactionHandler)
}
