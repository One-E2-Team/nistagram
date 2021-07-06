package main

import (
	"context"
	"fmt"
	"net/http"
	"nistagram/profile/handler"
	"nistagram/profile/model"
	"nistagram/profile/repository"
	"nistagram/profile/service"
	"nistagram/util"
	"nistagram/util/saga"
	"os"
	"time"

	"github.com/go-redis/redis/v8"
	"github.com/gorilla/mux"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func initDBs() (*gorm.DB, *redis.Client) {
	var (
		db  *gorm.DB
		err error
	)
	time.Sleep(5 * time.Second)
	var dbHost, dbPort, dbUsername, dbPassword = "localhost", "3306", "root", "root" // dev.db environment
	_, ok := os.LookupEnv("DOCKER_ENV_SET_PROD")                                     // production environment
	if ok {
		dbHost = "db_profile"
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
		db, err = gorm.Open(mysql.Open(dbUsername + ":" + dbPassword + "@tcp(" + dbHost + ":" + dbPort + ")/profile?charset=utf8mb4&parseTime=True&loc=Local"))

		if err != nil {
			fmt.Println("Cannot connect to database! Sleeping 10s and then retrying....")
			time.Sleep(10 * time.Second)
		} else {
			fmt.Println("Connected to the database.")
			break
		}
	}

	err = db.AutoMigrate(&model.Category{})
	if err != nil {
		return nil, nil
	}
	err = db.AutoMigrate(&model.Interest{})
	if err != nil {
		return nil, nil
	}
	err = db.AutoMigrate(&model.PersonalData{})
	if err != nil {
		return nil, nil
	}
	err = db.AutoMigrate(&model.ProfileSettings{})
	if err != nil {
		return nil, nil
	}
	err = db.AutoMigrate(&model.VerificationRequest{})
	if err != nil {
		return nil, nil
	}
	err = db.AutoMigrate(&model.Profile{})
	if err != nil {
		return nil, nil
	}
	err = db.AutoMigrate(&model.AgentRequest{})
	if err != nil {
		return nil, nil
	}

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
	util.InitMonitoring("profile", router)

	router.HandleFunc("/", handler.Register).Methods("POST")               // frontend func
	router.HandleFunc("/agent", util.RBAC(handler.RegisterAgent,
		"CREATE_AGENT", false)).Methods("POST")			 //frontend func
	router.HandleFunc("/search/{username}", handler.Search).Methods("GET") // frontend func
	router.HandleFunc("/search-for-tag/{username}", handler.SearchForTag).Methods("GET")
	router.HandleFunc("/get/{username}", handler.GetProfileByUsername).Methods("GET") // frontend func
	router.HandleFunc("/get-by-username/{username}",
		util.MSAuth(handler.GetProfileByUsername, []string{"post"})).Methods("GET")
	router.HandleFunc("/interests", handler.GetAllInterests).Methods("GET")           // frontend func
	router.HandleFunc("/categories", handler.GetAllCategories).Methods("GET")         // frontend func
	router.HandleFunc("/verification-request",
		util.RBAC(handler.CreateVerificationRequest, "CREATE_VERIFICATION_REQUEST", false)).Methods("POST") // frontend func
	router.HandleFunc("/verification-request",
		util.RBAC(handler.UpdateVerificationRequest, "UPDATE_VERIFICATION_REQUEST", false)).Methods("PUT") //frontend func
	router.HandleFunc("/change-profile-settings",
		util.RBAC(handler.ChangeProfileSettings, "EDIT_PROFILE_DATA", false)).Methods("PUT") // frontend func
	router.HandleFunc("/change-personal-data",
		util.RBAC(handler.ChangePersonalData, "EDIT_PROFILE_DATA", false)).Methods("PUT") // frontend func
	router.HandleFunc("/my-profile-settings",
		util.RBAC(handler.GetMyProfileSettings, "READ_PROFILE_DATA", false)).Methods("GET") // frontend func
	router.HandleFunc("/my-personal-data",
		util.RBAC(handler.GetMyPersonalData, "READ_PROFILE_DATA", false)).Methods("GET") // frontend func
	router.HandleFunc("/verification-requests",
		util.RBAC(handler.GetVerificationRequests, "READ_VERIFICATION_REQUESTS", true)).Methods("GET") // frontend func
	router.HandleFunc("/get-by-id/{id}",
		util.MSAuth(handler.GetProfileByID, []string{"connection", "post"})).Methods("GET")
	router.HandleFunc("/get-by-ids",
		util.MSAuth(handler.GetProfileUsernamesByIDs, []string{"postreaction", "monitoring", "post", "notification"})).Methods("POST")
	router.HandleFunc("/personal-data/{id}",
		util.MSAuth(handler.GetPersonalDataByProfileId, []string{"monitoring"})).Methods("GET")
	router.HandleFunc("/{id}",
		util.RBAC(handler.DeleteProfile, "DELETE_PROFILE", false)).Methods("DELETE") //frontend func
	router.HandleFunc("/send-agent-request",
		util.RBAC(handler.SendAgentRequest, "CREATE_AGENT_REQUEST", false)).Methods("POST") //frontend func
	router.HandleFunc("/agent-requests",
		util.RBAC(handler.GetAgentRequests, "READ_AGENT_REQUEST", true)).Methods("GET") //frontend func
	router.HandleFunc("/get-by-interests", handler.GetByInterests).Methods("POST")
	router.HandleFunc("/agent-request",
		util.RBAC(handler.ProcessAgentRequest, "EDIT_AGENT_REQUEST", false)).Methods("PUT") //frontend func
	router.HandleFunc("/test", handler.Test).Methods("GET")
	router.HandleFunc("/get-by-usernames", handler.GetProfileIdsByUsernames).Methods("POST")
	router.HandleFunc("/profile-interests/{id}",
		util.MSAuth(handler.GetProfileInterests, []string{"campaign"})).Methods("GET")
	fmt.Println("Starting server..")
	host, port := util.GetProfileHostAndPort()
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

func initSagaHandlers(handler *handler.Handler) []saga.ChannelHandler {
	temp := make([]saga.ChannelHandler, 0)
	temp = append(temp,
		saga.ChannelHandler{
			Channel: saga.ProfileChannel,
			Handler: handler.ChangePrivacyRollbackHandler,
		},
		saga.ChannelHandler{
			Channel: saga.ReplyChannel,
			Handler: func(client *redis.Client, message saga.Message) {
				fmt.Println("REPLY CHANNEL IN PROFILE")
			},
		})
	return temp
}

func main() {
	util.TracerInit("profile")
	db, client := initDBs()
	profileRepo := initProfileRepo(db, client)
	profileService := initService(profileRepo)
	profileHandler := initHandler(profileService)
	_ = util.SetupMSAuth("profile")
	saga.InitOrchestrator([]string{saga.ProfileChannel, saga.AuthChannel, saga.ConnectionChannel, saga.PostChannel, saga.ReplyChannel})
	go saga.Orch.Start()
	go saga.SubscribeAndRunPubSubHandlers(saga.Orch.Client, initSagaHandlers(profileHandler)...)
	handleFunc(profileHandler)
}
