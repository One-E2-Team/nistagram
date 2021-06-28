package main

import (
	"fmt"
	"github.com/gorilla/mux"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"net/http"
	"nistagram/campaign/handler"
	"nistagram/campaign/model"
	"nistagram/campaign/repository"
	"nistagram/campaign/service"
	"nistagram/util"
	"os"
	"time"
)

func initDB() *gorm.DB {
	var (
		db  *gorm.DB
		err error
	)
	time.Sleep(5 * time.Second)
	var dbHost, dbPort, dbUsername, dbPassword = "localhost", "3306", "root", "root" // dev.db environment
	_, ok := os.LookupEnv("DOCKER_ENV_SET_PROD")                                     // production environment
	if ok {
		dbHost = "db_campaign"
		dbPort = "3306"
		dbUsername = os.Getenv("DB_USERNAME")
		dbPassword = os.Getenv("DB_PASSWORD")
	} else {
		_, ok := os.LookupEnv("DOCKER_ENV_SET_DEV") // dev front environment
		if ok {
			dbHost = "db_relational"
			dbPort = "3306"
			dbUsername = os.Getenv("DB_USERNAME")
			dbPassword = os.Getenv("DB_PASSWORD")
		}
	}
	for {
		db, err = gorm.Open(mysql.Open(dbUsername + ":" + dbPassword + "@tcp(" + dbHost + ":" + dbPort + ")/campaign?charset=utf8mb4&parseTime=True&loc=Local"))

		if err != nil {
			fmt.Println("Cannot connect to database! Sleeping 10s and then retrying....")
			time.Sleep(10 * time.Second)
		} else {
			fmt.Println("Connected to the database.")
			break
		}
	}
	err = db.AutoMigrate(&model.Interest{})
	if err != nil {
		return nil
	}
	err = db.AutoMigrate(&model.Timestamp{})
	if err != nil {
		return nil
	}
	err = db.AutoMigrate(&model.CampaignRequest{})
	if err != nil {
		return nil
	}
	err = db.AutoMigrate(&model.CampaignParameters{})
	if err != nil {
		return nil
	}
	err = db.AutoMigrate(&model.Campaign{})
	if err != nil {
		return nil
	}
	return db
}

func initAuthRepo(db *gorm.DB) *repository.CampaignRepository {
	return &repository.CampaignRepository{Database: db}
}

func initAuthService(repo *repository.CampaignRepository) *service.CampaignService {
	return &service.CampaignService{CampaignRepository: repo}
}

func initAuthHandler(service *service.CampaignService) *handler.CampaignHandler {
	return &handler.CampaignHandler{CampaignService: service}
}

func handlerFunc(handler *handler.CampaignHandler) {
	fmt.Println("Campaign server started...")
	router := mux.NewRouter().StrictSlash(true)
	router.HandleFunc("/campaign", handler.CreateCampaign).Methods("POST")
	router.HandleFunc("/campaign/{id}", handler.UpdateCampaignParameters).Methods("PUT")
	router.HandleFunc("/interests/{campaignId}",
		util.MSAuth(handler.GetCurrentlyValidInterests, []string{"monitoring"})).Methods("GET")
	host, port := util.GetCampaignHostAndPort()
	var err error
	if util.DockerChecker() {
		err = http.ListenAndServeTLS(host+":"+port, "../cert.pem", "../key.pem", router)
	} else {
		err = http.ListenAndServe(":"+port, router)
	}
	if err != nil {
		fmt.Println(err)
		return
	}
}

func main() {
	db := initDB()
	campaignRepo := initAuthRepo(db)
	campaignService := initAuthService(campaignRepo)
	campaignHandler := initAuthHandler(campaignService)
	_ = util.SetupMSAuth("campaign")
	handlerFunc(campaignHandler)
}
