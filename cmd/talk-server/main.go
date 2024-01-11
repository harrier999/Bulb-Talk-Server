package main

import (
	"log"
	"net/http"
	"server/internal/handler/chatting"

	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}


	r := mux.NewRouter()
	r.HandleFunc("/chat", chatting.Handler).Methods("GET")
	log.Fatal(http.ListenAndServe(":18000", r))
}
