package main

import (
	"context"
	"fmt"
	"github.com/gorilla/handlers"
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
	var dbhost, dbport, dbusername, dbpassword string = "localhost", "8084", "root", "root" // dev.db environment
	_, ok := os.LookupEnv("DOCKER_ENV_SET_PROD") // production environment
	if ok {
		dbhost = "mongo1"
		dbport = "27017"
		dbusername = os.Getenv("MONGO_INITDB_ROOT_USERNAME")
		dbpassword = os.Getenv("MONGO_INITDB_ROOT_PASSWORD")
	} else {
		_, ok := os.LookupEnv("DOCKER_ENV_SET_DEV") // dev front environment
		if ok {
			dbhost = "mongo1"
			dbport = "27017"
			dbusername = os.Getenv("MONGO_INITDB_ROOT_USERNAME")
			dbpassword = os.Getenv("MONGO_INITDB_ROOT_PASSWORD")
		}
	}

	clientOptions := options.Client().ApplyURI("mongodb://"+dbusername+":"+dbpassword+"@" +dbhost+":"+dbport )
	for {
		client, err := mongo.Connect(context.TODO(), clientOptions)

		if err != nil {
			fmt.Println("Cannot connect to MongoDB! Sleeping 10s and then retrying....")
			time.Sleep(10 * time.Second)
		}else {
			fmt.Println("Connected to MongoDB")
			return client
		}
	}
	return nil

}

func initPostRepo(client *mongo.Client) *repository.PostRepository {
	return &repository.PostRepository{Client: client}
}

func initService(postRepo *repository.PostRepository) *service.PostService{
	return &service.PostService{PostRepository: postRepo}
}
func initHandler(postService *service.PostService) *handler.Handler{
	return &handler.Handler{PostService: postService}
}

func handleFunc(handler *handler.Handler) {
	router := mux.NewRouter().StrictSlash(true)
	router.HandleFunc("/", handler.GetAll).Methods("GET")
	router.HandleFunc("/public", handler.GetPublic).Methods("GET")
	router.HandleFunc("/", handler.Create).Methods("POST")
	router.HandleFunc("/user/privacy", handler.ChangePrivacy).Methods("PUT")
	router.HandleFunc("/user",handler.DeleteUserPosts).Methods("DELETE")
	router.HandleFunc("/user/username",handler.ChangeUsername).Methods("PUT")
	router.HandleFunc("/{postType}/{id}",handler.GetPost).Methods("GET")
	router.HandleFunc("/{postType}/{id}",handler.DeletePost).Methods("DELETE")
	router.HandleFunc("/{postType}/{id}",handler.UpdatePost).Methods("PUT")
	fmt.Printf("Starting server..")
	_, port := util.GetPostHostAndPort()
	http.ListenAndServe(":" + port, handlers.CORS(handlers.AllowedHeaders([]string{"X-Requested-With", "Content-Type",
		"Authorization"}), handlers.AllowedMethods([]string{"GET", "POST", "PUT", "HEAD", "OPTIONS"}),
		handlers.AllowedOrigins([]string{"*"}))(router))

}

func closeConnection(client *mongo.Client){
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
	postRepo := initPostRepo(client)
	postService := initService(postRepo)
	postHandler := initHandler(postService)
	handleFunc(postHandler)

}



