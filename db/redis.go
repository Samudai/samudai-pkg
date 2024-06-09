package db

import (
	"context"
	"os"

	"github.com/Samudai/samudai-pkg/logger"
	"github.com/redis/go-redis/v9"
)

var rdb *redis.Client

func InitRedis() {
	redisURL := os.Getenv("REDIS_URL")

	rdb = redis.NewClient(&redis.Options{
		Addr: redisURL,
	})

	_, err := rdb.Ping(context.Background()).Result()
	if err != nil {
		logger.LogMessage("error", "error connecting to Redis: %v", err)
		panic(err)
	}

	logger.LogMessage("info", "Redis connected")
}

func GetRedis() *redis.Client {
	return rdb
}

func CloseRedis() {
	err := rdb.Close()
	if err != nil {
		logger.LogMessage("error", "error closing Redis connection: %v", err)
	}
}
