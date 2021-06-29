package main

import (
	"fmt"
	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"net/http"
	"nistagram/agent/handler"
	"nistagram/agent/model"
	"nistagram/agent/repository"
	"nistagram/agent/service"
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
		dbHost = "db_agent"
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
		db, err = gorm.Open(mysql.Open(dbUsername + ":" + dbPassword + "@tcp(" + dbHost + ":" + dbPort + ")/agent?charset=utf8mb4&parseTime=True&loc=Local"))

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
	err = db.AutoMigrate(&model.Product{})
	if err != nil {
		return nil
	}
	err = db.AutoMigrate(&model.Statistics{})
	if err != nil {
		return nil
	}
	err = db.AutoMigrate(&model.AgentProduct{})
	if err != nil {
		return nil
	}
	err = db.AutoMigrate(&model.CampaignStat{})
	if err != nil {
		return nil
	}
	err = db.AutoMigrate(&model.InfluencerStat{})
	if err != nil {
		return nil
	}
	err = db.AutoMigrate(&model.InterestStat{})
	if err != nil {
		return nil
	}
	err = db.AutoMigrate(&model.Item{})
	if err != nil {
		return nil
	}
	err = db.AutoMigrate(&model.Order{})
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

func initProductRepo(db *gorm.DB) *repository.ProductRepository {
	return &repository.ProductRepository{Database: db}
}

func initProductService(repo *repository.ProductRepository) *service.ProductService {
	return &service.ProductService{ProductRepository: repo}
}

func initProductHandler(service *service.ProductService) *handler.ProductHandler {
	return &handler.ProductHandler{ProductService: service}
}

func handlerFunc(authHandler *handler.AuthHandler, productHandler *handler.ProductHandler) {
	fmt.Println("Agent application started...")
	router := mux.NewRouter().StrictSlash(true)
	router.HandleFunc("/register", authHandler.Register).Methods("POST")
	router.HandleFunc("/login", authHandler.LogIn).Methods("POST")
	router.HandleFunc("/validate/{id}/{uuid}", authHandler.ValidateUser).Methods("GET")
	router.HandleFunc("/product",
		authHandler.AuthService.RBAC(productHandler.CreateProduct, "CREATE_PRODUCT", false)).Methods("POST")
	router.HandleFunc("/product", productHandler.GetAllProducts).Methods("GET")
	router.HandleFunc("/product/{id}", productHandler.GetProductById).Methods("GET")
	router.HandleFunc("/product/{id}",
		authHandler.AuthService.RBAC(productHandler.DeleteProduct, "DELETE_PRODUCT", false)).Methods("DELETE")
	router.HandleFunc("/product",
		authHandler.AuthService.RBAC(productHandler.UpdateProduct, "UPDATE_PRODUCT", false)).Methods("PUT")
	router.HandleFunc("/order",
		authHandler.AuthService.RBAC(productHandler.CreateOrder, "CREATE_ORDER", false)).Methods("POST")
	_, ok := os.LookupEnv("DOCKER_ENV_SET_PROD")
	_, ok1 := os.LookupEnv("DOCKER_ENV_SET_DEV")
	var agentHost, agentPort = "localhost", "9000" // dev_db
	var err error
	if ok || ok1 {
		agentHost = "agent"
		agentPort = "8080"
		err = http.ListenAndServeTLS(agentHost+":"+agentPort, "../cert.pem", "../key.pem",
			handlers.CORS(handlers.AllowedOrigins([]string{"*"}),
				handlers.AllowedHeaders([]string{"Authorization", "Content-Type"}),
				handlers.AllowedMethods([]string{"GET", "HEAD", "POST", "PUT", "OPTIONS", "DELETE"}))(router))
	} else {
		err = http.ListenAndServe(":"+agentPort, handlers.CORS(handlers.AllowedOrigins([]string{"*"}),
			handlers.AllowedHeaders([]string{"Authorization", "Content-Type"}),
			handlers.AllowedMethods([]string{"GET", "HEAD", "POST", "PUT", "OPTIONS", "DELETE"}))(router))
	}
	if err != nil{
		fmt.Println(err)
		return
	}
}

func main() {
	db := initDB()
	authRepo := initAuthRepo(db)
	authService := initAuthService(authRepo)
	authHandler := initAuthHandler(authService)
	productRepo := initProductRepo(db)
	productService := initProductService(productRepo)
	productHandler := initProductHandler(productService)
	handlerFunc(authHandler, productHandler)
}
