package db

import (
	"context"
	"log"
	"os"
	"sync"

	redis "github.com/redis/go-redis/v9"
)

var ctx = context.Background()
var redisDB *redis.Client
var once sync.Once

func GetChattingHistoryClient() *redis.Client {
	once.Do(func() {
		log.Println("Connecting to redis...")
		redis_addr := os.Getenv("REDIS_ADDR")
		redis_password := os.Getenv("REDIS_PASS")
		redisDB = redis.NewClient(&redis.Options{
			Addr:     redis_addr,
			Password: redis_password,
			DB:       0,
		})
		_, err := redisDB.Ping(ctx).Result()
		if err != nil {
			log.Fatalln("Error connecting to redis. Error: ", err.Error())
			panic(err.Error())
		}
	})

	return redisDB
}

func GetContext() context.Context {
	return ctx
}
