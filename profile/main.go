package main

import (
	"context"
	"fmt"
	"github.com/go-redis/redis/v8"
	"github.com/gorilla/mux"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"net/http"
	"nistagram/profile/handler"
	"nistagram/profile/model"
	"nistagram/profile/repository"
	"nistagram/profile/service"
	"nistagram/util"
	"os"
	"time"
)

func initDBs() (*gorm.DB, *redis.Client) {
	var (
		db  *gorm.DB
		err error
	)
	time.Sleep(5 * time.Second)
	var dbhost, dbport, dbusername, dbpassword string = "localhost", "3306", "root", "root" // dev.db environment
	_, ok := os.LookupEnv("DOCKER_ENV_SET_PROD")                                            // production environment
	if ok {
		dbhost = "db_profile"
		dbport = "3306"
		dbusername = os.Getenv("DB_USERNAME")
		dbpassword = os.Getenv("DB_PASSWORD")
	} else {
		_, ok := os.LookupEnv("DOCKER_ENV_SET_DEV") // dev front environment
		if ok {
			dbhost = "db_relational"
			dbport = "3306"
			dbusername = "root"
			dbpassword = "root"
		}
	}
	for {
		db, err = gorm.Open(mysql.Open(dbusername + ":" + dbpassword + "@tcp(" + dbhost + ":" + dbport + ")/profile?charset=utf8mb4&parseTime=True&loc=Local"))

		if err != nil {
			fmt.Println("Cannot connect to database! Sleeping 10s and then retrying....")
			time.Sleep(10 * time.Second)
		} else {
			fmt.Println("Connected to the database.")
			break
		}
	}

	db.AutoMigrate(&model.Category{})
	db.AutoMigrate(&model.Interest{})
	db.AutoMigrate(&model.PersonalData{})
	db.AutoMigrate(&model.ProfileSettings{})
	db.AutoMigrate(&model.VerificationRequest{})
	db.AutoMigrate(&model.Profile{})

	//TODO: parametrize
	client := redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "",
		DB:       0,
	})
	Ctx := context.TODO()
	if err := client.Ping(Ctx).Err(); err != nil {
		fmt.Println(err)
		return db, nil
	}
	return db, client
}

func initProfileRepo(database *gorm.DB, client *redis.Client) *repository.ProfileRepository {
	return &repository.ProfileRepository{RelationalDatabase: database, Client: client, Context: context.TODO()}
}

func initService(profileRepo *repository.ProfileRepository) *service.ProfileService {
	return &service.ProfileService{ProfileRepository: profileRepo}
}

func initHandler(service *service.ProfileService) *handler.Handler {
	return &handler.Handler{ProfileService: service}
}

func handleFunc(handler *handler.Handler) {
	router := mux.NewRouter().StrictSlash(true)
	router.HandleFunc("/", handler.Register).Methods("POST")
	router.HandleFunc("/search/{username}", handler.Search).Methods("GET")
	router.HandleFunc("/get/{username}", handler.GetProfileByUsername).Methods("GET")
	router.HandleFunc("/change-profile-settings", handler.ChangeProfileSettings).Methods("PUT")
	router.HandleFunc("/change-personal-data", handler.ChangePersonalData).Methods("PUT")
	router.HandleFunc("/interests", handler.GetAllInterests).Methods("GET")
	router.HandleFunc("/my-profile-settings", handler.GetMyProfileSettings).Methods("GET")
	router.HandleFunc("/my-personal-data", handler.GetMyPersonalData).Methods("GET")
	router.HandleFunc("/get-by-id/{id}", handler.GetProfileByID).Methods("GET")
	router.HandleFunc("/test", handler.Test).Methods("GET")
	fmt.Println("Starting server..")
	_, port := util.GetProfileHostAndPort()
	http.ListenAndServe(":"+port, router)
}

func main() {
	db, client := initDBs()
	profileRepo := initProfileRepo(db, client)
	profileService := initService(profileRepo)
	handler := initHandler(profileService)
	handleFunc(handler)
}
