package main

import (
	"log"
	"net/http"
	"os"
	"server/internal/db/postgres_db"
	"server/internal/handler/chatting"
	"server/internal/handler/friends"
	"server/internal/handler/room"
	"server/internal/handler/user"
	"server/internal/repository/postgres"
	redisRepo "server/internal/repository/redis"
	"server/internal/service"
	"server/pkg/authenticator"

	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
	"github.com/redis/go-redis/v9"
)

func main() {

	err := godotenv.Load()
	if err != nil {
		log.Println(err.Error())
		log.Fatal("Error loading .env file")
	}

	postgresDB := postgres_db.GetPostgresClient()
	redisClient := getRedisClient()

	userRepo := postgres.NewPostgresUserRepository(postgresDB)
	friendRepo := postgres.NewPostgresFriendRepository(postgresDB)
	roomRepo := postgres.NewPostgresRoomRepository(postgresDB)
	messageRepo := redisRepo.NewRedisMessageRepository(redisClient)

	userService := service.NewUserService(userRepo)
	authService := service.NewAuthService(nil)
	friendService := service.NewFriendService(friendRepo, userRepo)
	roomService := service.NewRoomService(roomRepo)
	chatService := service.NewChatService(messageRepo)

	userHandler := user.NewHandler(userService, authService)
	friendHandler := friends.NewHandler(friendService)
	roomHandler := room.NewHandler(roomService)
	chatHandler := chatting.NewChatHandler(chatService)

	r := mux.NewRouter()

	r.HandleFunc("/signup", userHandler.SignUp).Methods("POST")
	r.HandleFunc("/login", userHandler.Login).Methods("POST")
	r.HandleFunc("/authenticate", userHandler.RequestAuthNumber).Methods("POST")
	r.HandleFunc("/checkauth", userHandler.CheckAuthNumber).Methods("POST")
	r.HandleFunc("/chat", chatHandler.HandleWebSocket).Methods("GET")
	r.HandleFunc("/messages", chatHandler.GetMessages).Methods("GET")

	authorizedRouter := r.PathPrefix("/auth").Subrouter()
	authorizedRouter.Use(authenticator.JWTMiddleware)
	authorizedRouter.HandleFunc("/getfriends", friendHandler.GetFriendList).Methods("POST")
	authorizedRouter.HandleFunc("/addfriends", friendHandler.AddFriend).Methods("POST")
	authorizedRouter.HandleFunc("/rooms", roomHandler.GetRoomList).Methods("POST")
	authorizedRouter.HandleFunc("/createrooms", roomHandler.CreateRoom).Methods("POST")

	authorizedRouter.HandleFunc("/blockfriend", friendHandler.BlockFriend).Methods("POST")
	authorizedRouter.HandleFunc("/unblockfriend", friendHandler.UnblockFriend).Methods("POST")
	authorizedRouter.HandleFunc("/adduser", roomHandler.AddUser).Methods("POST")
	authorizedRouter.HandleFunc("/removeuser", roomHandler.RemoveUser).Methods("POST")

	port := ":18000"
	log.Println("Server is successfully running on port " + port)
	log.Fatal(http.ListenAndServe(port, r))
}

func getRedisClient() *redis.Client {

	redisAddr := os.Getenv("REDIS_ADDR")
	if redisAddr == "" {
		redisAddr = "localhost:6379"
	}

	redisPassword := os.Getenv("REDIS_PASSWORD")
	redisDB := 0

	return redis.NewClient(&redis.Options{
		Addr:     redisAddr,
		Password: redisPassword,
		DB:       redisDB,
	})
}
