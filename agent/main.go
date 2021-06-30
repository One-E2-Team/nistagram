package main

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"nistagram/agent/handler"
	"nistagram/agent/model"
	"nistagram/agent/repository"
	"nistagram/agent/service"
	"nistagram/agent/util"
	"os"
	"time"

	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
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

func initPostService() *service.PostService {
	return &service.PostService{}
}

func initPostHandler(postService *service.PostService) *handler.PostHandler {
	return &handler.PostHandler{PostService: postService}
}

func initCampaignService() *service.CampaignService {
	return &service.CampaignService{}
}

func initCampaignHandler(campaignService *service.CampaignService) *handler.CampaignHandler {
	return &handler.CampaignHandler{CampaignService: campaignService}
}

func initConnectionService() *service.ConnectionService {
	return &service.ConnectionService{}
}

func initConnectionHandler(connectionService *service.ConnectionService) *handler.ConnectionHandler {
	return &handler.ConnectionHandler{ConnectionService: connectionService}
}

func handlerFunc(authHandler *handler.AuthHandler, productHandler *handler.ProductHandler,
	postHandler *handler.PostHandler, campaignHandler *handler.CampaignHandler, connectionHandler *handler.ConnectionHandler) {
	fmt.Println("Agent application started...")
	router := mux.NewRouter().StrictSlash(true)
	router.HandleFunc("/register", authHandler.Register).Methods("POST")
	router.HandleFunc("/login", authHandler.LogIn).Methods("POST")
	router.HandleFunc("/validate/{id}/{uuid}", authHandler.ValidateUser).Methods("GET")
	router.HandleFunc("/product",
		authHandler.AuthService.RBAC(productHandler.CreateProduct, "CREATE_PRODUCT", false)).Methods("POST")
	router.HandleFunc("/product", productHandler.GetAllProducts).Methods("GET")
	router.HandleFunc("/product/{id}",
		authHandler.AuthService.RBAC(productHandler.DeleteProduct, "DELETE_PRODUCT", false)).Methods("DELETE")
	router.HandleFunc("/product",
		authHandler.AuthService.RBAC(productHandler.UpdateProduct, "UPDATE_PRODUCT", false)).Methods("PUT")
	router.HandleFunc("/order",
		authHandler.AuthService.RBAC(productHandler.CreateOrder, "CREATE_ORDER", false)).Methods("POST")
	router.HandleFunc("/api-token",
		authHandler.AuthService.RBAC(authHandler.CreateAPIToken, "CREATE_TOKEN", false)).Methods("POST")
	router.HandleFunc("/my-posts",
		authHandler.AuthService.RBAC(postHandler.GetMyPosts, "READ_POSTS", true)).Methods("GET")
	router.HandleFunc("/report/campaign/{id}", campaignHandler.SaveCampaignReport).Methods("POST")
	router.HandleFunc("/my-campaigns",
		authHandler.AuthService.RBAC(campaignHandler.GetMyCampaigns, "READ_CAMPAIGNS", true)).Methods("GET")
	router.HandleFunc("/campaign",
		authHandler.AuthService.RBAC(campaignHandler.CreateCampaign, "CREATE_CAMPAIGN", false)).Methods("POST")
	router.HandleFunc("/interests",
		authHandler.AuthService.RBAC(campaignHandler.GetInterests, "CREATE_CAMPAIGN", true)).Methods("GET")
	router.HandleFunc("/followed-profiles",
		authHandler.AuthService.RBAC(connectionHandler.GetMyFollowedProfiles, "CREATE_CAMPAIGN", true)).Methods("GET")
	router.HandleFunc("/campaign/{id}",
		authHandler.AuthService.RBAC(campaignHandler.EditCampaign, "EDIT_CAMPAIGN", false)).Methods("PUT")

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
	if err != nil {
		fmt.Println(err)
		return
	}
}

func main() {
	//ProofOfConceptXMLDatabaseExistDB()
	db := initDB()
	authRepo := initAuthRepo(db)
	authService := initAuthService(authRepo)
	authHandler := initAuthHandler(authService)
	productRepo := initProductRepo(db)
	productService := initProductService(productRepo)
	productHandler := initProductHandler(productService)
	postService := initPostService()
	postHandler := initPostHandler(postService)
	campaignService := initCampaignService()
	campaignHandler := initCampaignHandler(campaignService)
	connectionService := initConnectionService()
	connectionHandler := initConnectionHandler(connectionService)
	handlerFunc(authHandler, productHandler, postHandler, campaignHandler, connectionHandler)
}

func ProofOfConceptXMLDatabaseExistDB() {
	resp, err := util.FuckingExistDBRequest(http.MethodPut, "/exist/rest/collection/document.xml", []byte("<root>\n  <fuel>-176916974.97817087</fuel>\n  <they as=\"radio\">801833596.2181821</they>\n  <taken partly=\"tip\">-1771145293.1620789</taken>\n  <afraid>involved</afraid>\n  <moon>\n    <window where=\"life\">340945322</window>\n    <list>percent</list>\n    <open wrong=\"club\">\n      <when radio=\"date\">growth</when>\n      <family>1174444784.414483</family>\n      <largest finger=\"answer\">throat</largest>\n      <strength do=\"him\">coast</strength>\n      <total>gate</total>\n      <tropical>although</tropical>\n    </open>\n    <minerals>1464602803.2550406</minerals>\n    <in>\n      <weight frog=\"powder\">science</weight>\n      <construction pull=\"shells\">-1501341086</construction>\n      <consider theory=\"science\">smoke</consider>\n      <actual perfect=\"nearby\">420013962</actual>\n      <help cold=\"student\">home</help>\n      <combine lack=\"flow\">1762325905.5516014</combine>\n    </in>\n    <single alone=\"actually\">569743812</single>\n  </moon>\n  <curve adjective=\"carbon\">1868724166.2018332</curve>\n</root>"), map[string]string{})
	if err != nil {
		fmt.Println(err)
	} else {
		fmt.Println(resp.StatusCode)
		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			fmt.Println(err)
		}
		fmt.Println(string(body))
	}
	fmt.Println()
	resp, err = util.FuckingExistDBRequest(http.MethodGet, "/exist/rest/collection/document.xml", []byte(""), map[string]string{})
	if err != nil {
		fmt.Println(err)
	} else {
		fmt.Println(resp.StatusCode)
		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			fmt.Println(err)
		}
		fmt.Println(string(body))
	}
	fmt.Println()
	var xParhAndOrQueryEtc = map[string]string{"_query": "//single[@alone='actually']"}
	resp, err = util.FuckingExistDBRequest(http.MethodGet, "/exist/rest/collection/document.xml"+util.GenerateFuckingExistDBHTTPRequestParametersQuery(xParhAndOrQueryEtc), []byte(""), map[string]string{})
	if err != nil {
		fmt.Println(err)
	} else {
		fmt.Println(resp.StatusCode)
		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			fmt.Println(err)
		}
		fmt.Println(string(body))
	}
	fmt.Println()
}
