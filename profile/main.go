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

	var db *gorm.DB
	var err error
	time.Sleep(5 * time.Second)
	var dbhost, dbport, dbusername, dbpassword string = "localhost", "3306", "root", "root"
	_, ok := os.LookupEnv("DOCKER_ENV_SET")
	if ok {
		dbhost = "mysql_profile"
		dbport = "3306"
		dbusername = os.Getenv("DB_USERNAME")
		dbpassword = os.Getenv("DB_PASSWORD")
	}
	for {
		db, err = gorm.Open(mysql.Open(dbusername + ":" + dbpassword + "@tcp(" + dbhost + ":" + dbport + ")/profile?charset=utf8mb4&parseTime=True&loc=Local"))

		if err != nil {
			fmt.Printf("Cannot connect to database! Sleeping 10s and then retrying....")
			time.Sleep(10 * time.Second)
		} else {
			fmt.Printf("Connected to the database.")
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
	fmt.Printf("Starting server..")
	var port string = "8083"
	_, ok := os.LookupEnv("DOCKER_ENV_SET")
	if ok {
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
