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
	"os"
)

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
	router.HandleFunc("/{postType}", handler.Create).Methods("POST")
	router.HandleFunc("/{postType}/{id}",handler.GetPost).Methods("GET")
	router.HandleFunc("/{postType}/{id}",handler.DeletePost).Methods("DELETE")
	router.HandleFunc("/{postType}/{id}",handler.UpdatePost).Methods("PUT")
	fmt.Printf("Starting server..")
	var port string = "8085"                     // dev.db environ
	_, ok := os.LookupEnv("DOCKER_ENV_SET_PROD") // dev production environment
	_, ok1 := os.LookupEnv("DOCKER_ENV_SET_DEV") // dev front environment
	if ok || ok1 {
		port = "8080"
	}
	http.ListenAndServe(":"+port, router)
}

func closeConnection(client *mongo.Client){
	err := client.Disconnect(context.TODO())
	printErrorTrack(err)
	fmt.Println("Connection to MongoDB closed.")
}


func main() {
	clientOptions := options.Client().ApplyURI("mongodb://root:root@localhost:8084")
	client, err := mongo.Connect(context.TODO(), clientOptions)
	printErrorTrack(err)
	fmt.Println("Connected to mongodb")
	defer closeConnection(client)


	postRepo := initPostRepo(client)
	postService := initService(postRepo)
	postHandler := initHandler(postService)
	handleFunc(postHandler)

}


func printErrorTrack(err error) {
	if err!=nil {
		fmt.Println(err)
	}
}
