package main

import (
	"context"
	"fmt"
	"github.com/gorilla/mux"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"net/http"
	"nistagram/post/handler"
	"nistagram/post/repository"
	"nistagram/post/service"
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
			initCollections(client)
			return client
		}
	}
}

func initCollections(client *mongo.Client) {
	const postsCollectionName = "posts"
	const postDbName = "postdb"
	createCollection(client, postDbName, postsCollectionName)
}
func createCollection(client *mongo.Client,dbName string, collectionName string) {
	if err := client.Database(dbName).CreateCollection(context.TODO(), collectionName); err != nil {
		fmt.Println(err)
	} else {
		fmt.Println("Create " + collectionName + " collection success")
	}
}

func initPostRepo(client *mongo.Client) *repository.PostRepository {
	return &repository.PostRepository{Client: client}
}

func initService(postRepo *repository.PostRepository) *service.PostService {
	return &service.PostService{PostRepository: postRepo}
}
func initHandler(postService *service.PostService) *handler.Handler {
	return &handler.Handler{PostService: postService}
}

func handleFunc(handler *handler.Handler) {
	router := mux.NewRouter().StrictSlash(true)

	router.HandleFunc("/profile/{username}", handler.GetProfilesPosts).Methods("GET")            // frontend func
	router.HandleFunc("/public", handler.GetPublic).Methods("GET")                               // frontend func
	router.HandleFunc("/public/location/{value}", handler.SearchPublicByLocation).Methods("GET") // frontend func
	router.HandleFunc("/public/hashtag/{value}", handler.SearchPublicByHashTag).Methods("GET")   // frontend func
	router.HandleFunc("/my",
		util.RBAC(handler.GetMyPosts, "READ_NOT_ONLY_PUBLIC_POSTS", true)).Methods("GET") // frontend func
	router.HandleFunc("/agent-my", util.AgentAuth(handler.GetMyPosts)).Methods("GET") // frontend func
	router.HandleFunc("/homePage",
		util.RBAC(handler.GetPostsForHomePage, "READ_NOT_ONLY_PUBLIC_POSTS", true)).Methods("GET") // frontend func
	router.HandleFunc("/",
		util.RBAC(handler.Create, "CREATE_POST", false)).Methods("POST") // frontend func
	router.HandleFunc("/user/{loggedUserId}/privacy",
		util.MSAuth(handler.ChangePrivacy, []string{"profile"})).Methods("PUT")
	router.HandleFunc("/user/{id}",
		util.MSAuth(handler.DeleteUserPosts, []string{"profile"})).Methods("DELETE")
	router.HandleFunc("/user/{loggedUserId}/username",
		util.MSAuth(handler.ChangeUsername, []string{"profile"})).Methods("PUT")
	router.HandleFunc("/post/{id}",
		util.MSAuth(handler.GetPost, []string{"postreaction"})).Methods("GET")
	router.HandleFunc("/posts",
		util.MSAuth(handler.GetPosts, []string{"postreaction", "campaign"})).Methods("POST")
	router.HandleFunc("/make-campaign/{id}/{agentID}",
		util.MSAuth(handler.MakeCampaign, []string{"campaign"})).Methods("POST")
	router.HandleFunc("/{id}",
		util.RBAC(handler.DeletePost, "DELETE_POST", false)).Methods("DELETE") // frontend func
	router.HandleFunc("/{id}", handler.UpdatePost).Methods("PUT")
	router.HandleFunc("/media/{id}",
		util.MSAuth(handler.GetMediaById, []string{"monitoring"})).Methods("GET")
	fmt.Println("Starting server..")
	host, port := util.GetPostHostAndPort()
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
	util.TracerInit("post")
	client := initDB()
	defer closeConnection(client)
	postRepo := initPostRepo(client)
	postService := initService(postRepo)
	postHandler := initHandler(postService)
	_ = util.SetupMSAuth("post")
	util.InitMonitoring("post")
	handleFunc(postHandler)
}
