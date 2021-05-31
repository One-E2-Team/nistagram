package main

import (
	"fmt"
	"github.com/gorilla/mux"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"net/http"
	"nistagram/auth/model"
)

func test(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Test")
}

func initDB() bool {
	db, err := gorm.Open(mysql.Open("root:root@tcp(127.0.0.1:3307)/auth?charset=utf8mb4&parseTime=True&loc=Local"))
	if err != nil {
		fmt.Printf("Cannot connect to database!")
		return false
	}
	db.AutoMigrate(&model.Privilege{})
	db.AutoMigrate(&model.Credentials{})
	db.AutoMigrate(&model.Role{})
	db.AutoMigrate(&model.User{})
	return true
}

func main() {
	if initDB() {
		fmt.Println("Auth server started...")
		router := mux.NewRouter().StrictSlash(true)
		router.HandleFunc("/", test).Methods("GET")
		http.ListenAndServe(":8082", router)
	}
}
