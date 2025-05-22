package main

import (
	"fmt"
	"log"
	"net/http"
	
	"urlshortner/database"
	"urlshortner/handlers"

	"github.com/gorilla/mux"
)

func main() {
	fmt.Println("Starting server...")
	database.InitDB("urlshortener.db")

	r := mux.NewRouter()

	r.HandleFunc("/shorten", handlers.CreateShortURL).Methods("POST")
	r.HandleFunc("/u/{code}", handlers.GetOriginalURL).Methods("GET")
	r.HandleFunc("/u/{code}", handlers.UpdateShortURL).Methods("PUT")
	r.HandleFunc("/u/{code}", handlers.DeleteShortURL).Methods("DELETE")
	r.HandleFunc("/stats/{code}", handlers.GetStats).Methods("GET")

	log.Fatal(http.ListenAndServe(":8080", r))
}
