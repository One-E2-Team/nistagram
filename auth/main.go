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
		dbHost = "db_auth"
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
	util.InitMonitoring("auth", router)

	router.HandleFunc("/login", handler.LogIn).Methods("POST")                          //frontend func
	router.HandleFunc("/login/apitoken", handler.LogInAgentAPI).Methods("POST")
	router.HandleFunc("/agent/test", util.AgentAuth(func(writer http.ResponseWriter, request *http.Request) {
		writer.WriteHeader(http.StatusOK)
		_, _ = writer.Write([]byte("{\"energy\":\"big dick energy\"}"))
		writer.Header().Set("Content-Type", "application/json")
	})).Methods("GET")
	router.HandleFunc("/request-recovery", handler.RequestPassRecovery).Methods("POST") //frontend func
	router.HandleFunc("/recover", handler.ChangePassword).Methods("POST")               //frontend func
	/*router.HandleFunc("/validate/{id}/{uuid}/{qruuid}", handler.ValidateUser).Methods("GET") //frontend func*/
	router.HandleFunc("/validate/{id}/{uuid}", handler.ValidateUser).Methods("GET") //frontend func
	router.HandleFunc("/api-token",
		util.RBAC(handler.GetAgentAPIToken, "READ_API_TOKEN", false)).Methods("GET") //frontend func
	router.HandleFunc("/register",
		util.MSAuth(handler.Register, []string{"profile"})).Methods("POST")
	router.HandleFunc("/update-user",
		util.MSAuth(handler.UpdateUser, []string{"profile"})).Methods("POST")
	router.HandleFunc("/privileges/{profileId}",
		util.MSAuth(handler.GetPrivileges,
			[]string{"auth", "connection", "post", "profile", "postreaction", "campaign", "notification", "monitoring"})).Methods("GET")
	router.HandleFunc("/ban/{profileID}",
		util.MSAuth(handler.BanUser, []string{"profile"})).Methods("DELETE")
	router.HandleFunc("/make-agent/{profileID}",
		util.MSAuth(handler.MakeAgent, []string{"profile"})).Methods("PUT")
	host, port := util.GetAuthHostAndPort()
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
	util.TracerInit("auth")
	db := initDB()
	authRepo := initAuthRepo(db)
	authService := initAuthService(authRepo)
	authHandler := initAuthHandler(authService)
	_ = util.SetupMSAuth("auth")
	handlerFunc(authHandler)
}
