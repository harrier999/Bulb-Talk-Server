package redis_db

import (
	"context"
	"log/slog"
	"os"
	"server/pkg/log"
	"sync"

	redis "github.com/redis/go-redis/v9"
)

type RedisInstance struct {
	Client *redis.Client
	Ctx    context.Context
	logger *slog.Logger
}

var (
	instance      *RedisInstance
	instanceMutex sync.Mutex
)

func GetChattingHistoryClient() *redis.Client {
	if instance == nil {
		instanceMutex.Lock()
		defer instanceMutex.Unlock()

		if instance == nil {
			instance = &RedisInstance{
				Ctx:    context.Background(),
				logger: log.NewColorLog(),
			}
			instance.connect()
		}
	}

	return instance.Client
}

func GetContext() context.Context {
	if instance == nil {
		GetChattingHistoryClient()
	}
	return instance.Ctx
}

func (r *RedisInstance) connect() {
	r.logger.Info("Connecting to Redis...")

	redisAddr := os.Getenv("REDIS_ADDR")
	redisPassword := os.Getenv("REDIS_PASS")

	r.Client = redis.NewClient(&redis.Options{
		Addr:     redisAddr,
		Password: redisPassword,
		DB:       0,
	})

	_, err := r.Client.Ping(r.Ctx).Result()
	if err != nil {
		r.logger.Error("Error connecting to Redis", "error", err.Error())
		os.Exit(1)
	}

	r.logger.Info("Successfully connected to Redis")
}

func CloseConnection() error {
	if instance != nil && instance.Client != nil {
		return instance.Client.Close()
	}
	return nil
}
