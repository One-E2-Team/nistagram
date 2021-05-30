package main

import (
	"fmt"
	"github.com/gorilla/mux"
	"net/http"
	"nistagram/auth/model"
)

func test(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Test")
}

func main() {
	p1 := model.Privilege{
		Name: "name1",
	}
	p2 := model.Privilege{
		Name: "name2",
	}
	r := model.Role{
		Name:       "test",
		Privileges: []model.Privilege{p1, p2},
	}
	fmt.Println(p1, p2, r)
	router := mux.NewRouter().StrictSlash(true)
	router.HandleFunc("/", test).Methods("GET")

	http.ListenAndServe(":8082", router)
}
