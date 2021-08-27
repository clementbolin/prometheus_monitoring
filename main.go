package main

import (
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

func main() {
	router := mux.NewRouter()

	router.PathPrefix("/").Handler(http.FileServer(http.Dir("./static/")))

	log.Println("Serving request on port 9000: localhost:9000")
	err := http.ListenAndServe(":9000", router)
	log.Fatalln(err)
}