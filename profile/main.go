package main

import (
	"fmt"
	"github.com/gorilla/mux"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"net/http"
	"nistagram/profile/handler"
	"nistagram/profile/model"
	"nistagram/profile/repository"
	"nistagram/profile/service"
	"os"
	"time"
)

func initDB() *gorm.DB {

	var (db *gorm.DB; err error)
	time.Sleep(5 * time.Second)
	var dbhost, dbport, dbusername, dbpassword string = "localhost", "3306", "root", "root" // dev.db environment
	_, ok := os.LookupEnv("DOCKER_ENV_SET_PROD") // production environment
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

	return db
}

func initProfileRepo(database *gorm.DB) *repository.ProfileRepository {
	return &repository.ProfileRepository{Database: database}
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
	router.HandleFunc("/{username}", handler.GetProfileByUsername).Methods("GET")
	fmt.Printf("Starting server..")
	var port string = "8083" // dev.db environ
	_, ok := os.LookupEnv("DOCKER_ENV_SET_PROD") // dev production environment
	_, ok1 := os.LookupEnv("DOCKER_ENV_SET_DEV") // dev front environment
	if ok || ok1 {
		port = "8080"
	}
	http.ListenAndServe(":" + port, router)
}

func main() {
	db := initDB()
	profileRepo := initProfileRepo(db)
	profileService := initService(profileRepo)
	handler := initHandler(profileService)
	handleFunc(handler)
}
