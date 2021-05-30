package main

import (
	"./model"
	"fmt"
	"github.com/gorilla/mux"
	"net/http"
)

func test(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Test")
}

func main(){
	a := model.Profile{
		Biography: "kurac",
	}
	fmt.Println(a)
	router := mux.NewRouter().StrictSlash(true)
	router.HandleFunc("/", test).Methods("GET")

	http.ListenAndServe(":8081", router)
}
