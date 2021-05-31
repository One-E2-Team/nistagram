package main

import (
	"fmt"
	"github.com/gorilla/mux"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"net/http"
	"nistagram/profile/model"
)

func test(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Test")
}

func initDB(){
	db, err := gorm.Open(mysql.Open("root:root@tcp(127.0.0.1:3306)/profile?charset=utf8mb4&parseTime=True&loc=Local"))

	if err != nil{
		fmt.Printf("Cannot connect to database!")
		return
	}

	db.AutoMigrate(&model.Category{})
	db.AutoMigrate(&model.Interest{})
	db.AutoMigrate(&model.PersonalData{})
	db.AutoMigrate(&model.ProfileSettings{})
	db.AutoMigrate(&model.VerificationRequest{})
	db.AutoMigrate(&model.Profile{})

}

func main() {

	initDB()

	fmt.Printf("Starting server..")

	router := mux.NewRouter().StrictSlash(true)
	router.HandleFunc("/", test).Methods("GET")

	http.ListenAndServe(":8081", router)
}
