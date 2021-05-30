package main

import (
	"fmt"
	"net/http"
	"github.com/gorilla/mux"
)

func test(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Test")
}

func main(){
	router := mux.NewRouter().StrictSlash(true)
	router.HandleFunc("/", test).Methods("GET")

	http.ListenAndServe(":8081", router)
}
