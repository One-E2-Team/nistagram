package main

import (
	"fmt"
	"github.com/gorilla/mux"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"net/http"
	"nistagram/auth/handler"
	"nistagram/auth/model"
	"nistagram/auth/repository"
	"nistagram/auth/service"
	"os"
	"time"
)

func initDB() *gorm.DB {
	var (
		db  *gorm.DB
		err error
	)
	time.Sleep(5 * time.Second)
	var dbHost, dbPort, dbUsername, dbPassword string = "localhost", "3306", "root", "root" // dev.db environment
	_, ok := os.LookupEnv("DOCKER_ENV_SET_PROD")                                            // production environment
	if ok {
		dbHost = "db_auth"
		dbPort = "3306"
		dbUsername = os.Getenv("DB_USERNAME")
		dbPassword = os.Getenv("DB_PASSWORD")
	} else {
		_, ok := os.LookupEnv("DOCKER_ENV_SET_DEV") // dev front environment
		if ok {
			dbHost = "db_relational"
			dbPort = "3306"
			dbUsername = "root"
			dbPassword = "root"
		}
	}
	for {
		db, err = gorm.Open(mysql.Open(dbUsername + ":" + dbPassword + "@tcp(" + dbHost + ":" + dbPort + ")/auth?charset=utf8mb4&parseTime=True&loc=Local"))

		if err != nil {
			fmt.Println("Cannot connect to database! Sleeping 10s and then retrying....")
			time.Sleep(10 * time.Second)
		} else {
			fmt.Println("Connected to the database.")
			break
		}
	}
	err = db.AutoMigrate(&model.Privilege{})
	if err != nil {
		return nil
	}
	err = db.AutoMigrate(&model.Role{})
	if err != nil {
		return nil
	}
	err = db.AutoMigrate(&model.User{})
	if err != nil {
		return nil
	}
	return db
}

func initAuthRepo(db *gorm.DB) *repository.AuthRepository {
	return &repository.AuthRepository{Database: db}
}

func initAuthService(repo *repository.AuthRepository) *service.AuthService {
	return &service.AuthService{AuthRepository: repo}
}

func initAuthHandler(service *service.AuthService) *handler.AuthHandler {
	return &handler.AuthHandler{AuthService: service}
}

func handlerFunc(handler *handler.AuthHandler) {
	fmt.Println("Auth server started...")
	router := mux.NewRouter().StrictSlash(true)
	router.HandleFunc("/login", handler.LogIn).Methods("POST")
	router.HandleFunc("/register", handler.Register).Methods("POST")
	router.HandleFunc("/request-recovery", handler.RequestPassRecovery).Methods("POST")
	router.HandleFunc("/recover", handler.ChangePassword).Methods("POST")
	router.HandleFunc("/update-user", handler.UpdateUser).Methods("POST")
	router.HandleFunc("/validate/{id}/{uuid}", handler.ValidateUser).Methods("GET")
	var port string = "8000"                     // dev.db environ
	_, ok := os.LookupEnv("DOCKER_ENV_SET_PROD") // dev production environment
	_, ok1 := os.LookupEnv("DOCKER_ENV_SET_DEV") // dev front environment
	if ok || ok1 {
		port = "8080"
	}
	err := http.ListenAndServe(":"+port, router)
	if err != nil {
		fmt.Println(err)
		return
	}
}

func main() {
	db := initDB()
	authRepo := initAuthRepo(db)
	authService := initAuthService(authRepo)
	authHandler := initAuthHandler(authService)
	handlerFunc(authHandler)
}
