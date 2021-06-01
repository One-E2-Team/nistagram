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
	"time"
)

func initDB() *gorm.DB {

	var db *gorm.DB
	var err error

	for {
		db, err = gorm.Open(mysql.Open("root:root@tcp(127.0.0.1:3306)/profile?charset=utf8mb4&parseTime=True&loc=Local"))

		if err != nil {
			fmt.Printf("Cannot connect to database! Sleeping 10s and then retrying....")
			time.Sleep(10 * time.Second)
		} else {
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
	http.ListenAndServe(":8083", router)
}

func main() {
	db := initDB()
	profileRepo := initProfileRepo(db)
	profileService := initService(profileRepo)
	handler := initHandler(profileService)
	handleFunc(handler)
}
