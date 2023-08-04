package db

import (
	redis "github.com/redis/go-redis/v9"
	"os"
	"context"
	"log"
)

var ctx = context.Background()

func GetRedisClient() *redis.Client {
	
	redis_addr := os.Getenv("REDIS_ADDR")
	redis_password := os.Getenv("REDIS_PASS")
	log.Println("redis addr is ", redis_addr)
	log.Println("redis password is ", redis_password)

	redisDB := redis.NewClient(&redis.Options{
		Addr:     redis_addr,
		Password: redis_password,
		DB:       0,
	})
	_, err := redisDB.Ping(ctx).Result()
	if err != nil {
		panic(err.Error())
	}

	return redisDB
}

func GetContext() context.Context {
	return ctx
}