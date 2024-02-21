package main

import (
	"log"
	"server/pkg/authenticator"

	"net/http"
	"server/internal/handler/chatting"
	"server/internal/handler/friends"
	"server/internal/handler/room"
	"server/internal/handler/user"

	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Println(err.Error())
		log.Fatal("Error loading .env file")
	}
	r := mux.NewRouter()

	r.HandleFunc("/signup", user.SignUp).Methods("POST")
	r.HandleFunc("/login", user.Login).Methods("POST")
	r.HandleFunc("/authenticate", user.RequestAuthNumber).Methods("POST")
	r.HandleFunc("/checkauth", user.CheckAuthNumber).Methods("POST")
	r.HandleFunc("/chat", chatting.Handler).Methods("GET")

	authroizedRouter := r.PathPrefix("/auth").Subrouter()
	authroizedRouter.Use(authenticator.JWTMiddleware)
	authroizedRouter.HandleFunc("/getfriends", friends.GetFriendList).Methods("POST")
	authroizedRouter.HandleFunc("/addfriends", friends.AddFriend).Methods("POST")
	authroizedRouter.HandleFunc("/rooms", room.GetRoomListHandler).Methods("POST")
	authroizedRouter.HandleFunc("/createrooms", room.CreateRoomHandler).Methods("POST")

	port := ":18000"

	log.Println("Server is successfully running on port " + port)
	log.Fatal(http.ListenAndServe(port, r))

}
