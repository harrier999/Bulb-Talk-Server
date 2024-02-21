package main

import (
	"log"
	"net/http"

	"server/internal/handler/user"
	"server/internal/handler/friends"
	"server/internal/models/orm"
	"server/internal/db/postgres_db"
	"github.com/joho/godotenv"

	"github.com/gorilla/mux"
)

func main() {

	godotenv.Load(".env")

	r := mux.NewRouter()


	postgres_db.GetPostgresClient().AutoMigrate(&orm.User{}) 
	postgres_db.GetPostgresClient().AutoMigrate(&orm.Friend{})
	r.HandleFunc("/requestauth", user.RequestAuthNumber).Methods("POST")
	r.HandleFunc("/checkauth", user.CheckAuthNumber).Methods("POST")
	r.HandleFunc("/signUp", user.SignUp).Methods("POST")
	r.HandleFunc("/login", user.Login).Methods("POST")
	r.HandleFunc("/friends", friends.GetFriendList).Methods("POST")
	r.HandleFunc("/addfriend", friends.AddFriend).Methods("POST")
	
	log.Fatal(http.ListenAndServe(":2222", r))
}
