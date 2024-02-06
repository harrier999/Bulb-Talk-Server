package main

import (
	"log"
	"server/pkg/authenticator"

	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
	"server/internal/handler/room"
	"server/internal/handler/user"
	"server/internal/handler/friends"
	"net/http"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Println(err.Error())
		log.Fatal("Error loading .env file")
		
	}
	r := mux.NewRouter()
	r.Use(authenticator.JWTMiddleware)

	r.HandleFunc("/signup", user.SignUp).Methods("POST")
	r.HandleFunc("/login", user.Login).Methods("POST")
	r.HandleFunc("/authenticate", user.RequestAuthNumber).Methods("POST")
	r.HandleFunc("/checkauth", user.CheckAuthNumber).Methods("POST")

	r.HandleFunc("/friends", friends.GetFriendList).Methods("GET")
	r.HandleFunc("/friends", friends.AddFriend).Methods("POST")

	r.HandleFunc("/rooms", room.GetRoomListHandler).Methods("GET")
	r.HandleFunc("/rooms", room.CreateRoomHandler).Methods("POST")

	log.Fatal(http.ListenAndServe(":18002", r))
	
}
