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
)

func test(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Test")
}

func initDB() *gorm.DB {
	db, err := gorm.Open(mysql.Open("root:root@tcp(127.0.0.1:3307)/auth?charset=utf8mb4&parseTime=True&loc=Local"))
	if err != nil {
		fmt.Printf("Cannot connect to database!")
		return nil
	}
	db.AutoMigrate(&model.Privilege{})
	db.AutoMigrate(&model.Role{})
	db.AutoMigrate(&model.User{})
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
	http.ListenAndServe(":8000", router)
}

func main() {
	db := initDB()
	authRepo := initAuthRepo(db)
	authService := initAuthService(authRepo)
	authHandler := initAuthHandler(authService)
	handlerFunc(authHandler)
}
