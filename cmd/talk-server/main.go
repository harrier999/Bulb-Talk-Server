package main

import (
	"log"
	"net/http"
	"net/url"
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
	"strings"

	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
	"github.com/redis/go-redis/v9"
)

// 환경 변수에서 허용된 오리진 목록을 가져옵니다
func getAllowedOrigins() []string {
	allowedOriginsStr := os.Getenv("ALLOWED_ORIGINS")
	if allowedOriginsStr == "" {
		return []string{} // 기본값 없음
	}
	return strings.Split(allowedOriginsStr, ",")
}

// CORS 미들웨어 함수
func corsMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		origin := r.Header.Get("Origin")
		if origin == "" {
			next.ServeHTTP(w, r)
			return
		}

		// 허용할 도메인 확인
		allowedOrigins := getAllowedOrigins()
		allowed := false

		// '*'가 허용된 오리진 목록에 있으면 모든 오리진 허용
		if contains(allowedOrigins, "*") {
			allowed = true
		} else {
			// 로컬호스트 확인 (디버깅 용도)
			parsedOrigin, err := url.Parse(origin)
			if err == nil {
				host := parsedOrigin.Host
				if strings.HasPrefix(host, "localhost") ||
					strings.HasPrefix(host, "127.0.0.1") ||
					strings.HasSuffix(host, ".localhost") {
					allowed = true
				}
			}

			// 허용된 도메인 목록 확인
			if !allowed {
				for _, allowedOrigin := range allowedOrigins {
					// 와일드카드 도메인 처리 (*.example.com)
					if strings.HasPrefix(allowedOrigin, "*.") {
						suffix := allowedOrigin[1:] // "*."를 제거
						if strings.HasSuffix(origin, suffix) || strings.Contains(origin, suffix) {
							allowed = true
							break
						}
					} else if strings.Contains(origin, allowedOrigin) {
						allowed = true
						break
					}
				}
			}
		}

		if allowed {
			// CORS 헤더 설정
			w.Header().Set("Access-Control-Allow-Origin", origin)
			w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS, PATCH")
			w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization, X-Requested-With, Accept")
			w.Header().Set("Access-Control-Allow-Credentials", "true")
			w.Header().Set("Access-Control-Max-Age", "3600")
		}

		// OPTIONS 요청 처리 (preflight 요청)
		if r.Method == "OPTIONS" {
			w.WriteHeader(http.StatusOK)
			return
		}

		next.ServeHTTP(w, r)
	})
}

// 문자열 슬라이스에 특정 문자열이 포함되어 있는지 확인하는 헬퍼 함수
func contains(slice []string, item string) bool {
	for _, s := range slice {
		if strings.TrimSpace(s) == item {
			return true
		}
	}
	return false
}

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

	// CORS 미들웨어 적용
	r.Use(corsMiddleware)

	// OPTIONS 메서드를 모든 경로에 대해 허용
	r.Methods("OPTIONS").HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})

	r.HandleFunc("/signup", userHandler.SignUp).Methods("POST", "OPTIONS")
	r.HandleFunc("/login", userHandler.Login).Methods("POST", "OPTIONS")
	r.HandleFunc("/authenticate", userHandler.RequestAuthNumber).Methods("POST", "OPTIONS")
	r.HandleFunc("/checkauth", userHandler.CheckAuthNumber).Methods("POST", "OPTIONS")
	r.HandleFunc("/chat", chatHandler.HandleWebSocket).Methods("GET", "OPTIONS")
	r.HandleFunc("/messages", chatHandler.GetMessages).Methods("GET", "OPTIONS")

	authorizedRouter := r.PathPrefix("/auth").Subrouter()
	authorizedRouter.Use(authenticator.JWTMiddleware)

	// 인증된 라우터에도 OPTIONS 메서드 허용
	authorizedRouter.Methods("OPTIONS").HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})

	// 친구 관련 RESTful API 엔드포인트
	authorizedRouter.HandleFunc("/friends", friendHandler.GetFriendList).Methods("GET", "OPTIONS")
	authorizedRouter.HandleFunc("/friends", friendHandler.AddFriend).Methods("POST", "OPTIONS")
	authorizedRouter.HandleFunc("/friends/{friendId}/block", friendHandler.BlockFriend).Methods("PUT", "OPTIONS")
	authorizedRouter.HandleFunc("/friends/{friendId}/unblock", friendHandler.UnblockFriend).Methods("PUT", "OPTIONS")

	// 채팅방 관련 RESTful API 엔드포인트
	authorizedRouter.HandleFunc("/rooms", roomHandler.GetRoomList).Methods("GET", "OPTIONS")
	authorizedRouter.HandleFunc("/rooms", roomHandler.CreateRoom).Methods("POST", "OPTIONS")
	authorizedRouter.HandleFunc("/rooms/{roomId}/users", roomHandler.AddUser).Methods("POST", "OPTIONS")
	authorizedRouter.HandleFunc("/rooms/{roomId}/users/{userId}", roomHandler.RemoveUser).Methods("DELETE", "OPTIONS")

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
