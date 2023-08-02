package main

import (
	"log"
	"net/http"
	"server/api/chatting"

	"github.com/gorilla/mux"
	_ "github.com/jinzhu/gorm/dialects/postgres"
)

func main() {

	r := mux.NewRouter()
	r.HandleFunc("/chat", chatting.Handler).Methods("GET")
	log.Fatal(http.ListenAndServe(":18000", r))
}
