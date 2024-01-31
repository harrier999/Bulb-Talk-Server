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



	r.HandleFunc("/signup", user.SignUp).Methods("POST")
	r.HandleFunc("/login", user.Login).Methods("POST")
	r.HandleFunc("/authenticate", user.RequestAuthNumber).Methods("POST")
	r.HandleFunc("/checkauth", user.CheckAuthNumber).Methods("POST")

	authroizedRouter := r.PathPrefix("/auth").Subrouter()
	authroizedRouter.Use(authenticator.JWTMiddleware)
	authroizedRouter.HandleFunc("/getfriends", friends.GetFriendList).Methods("POST")
	authroizedRouter.HandleFunc("/addfriends", friends.AddFriend).Methods("POST")
	authroizedRouter.HandleFunc("/rooms", room.GetRoomListHandler).Methods("POST")
	authroizedRouter.HandleFunc("/createrooms", room.CreateRoomHandler).Methods("POST")

	port := ":18000"
	
	log.Println("Server is successfully running on port " + port)
	log.Fatal(http.ListenAndServe(port, authroizedRouter))
	
}
