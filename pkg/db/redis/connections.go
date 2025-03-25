package redis

import (
	"SSO/internal/config"
	"context"
	"fmt"
	"github.com/go-redis/redis/v8"
	"log"
)

var (
	RedisClient *redis.Client
	ctx         = context.Background()
)

// InitializeRedis инициализирует соединение с Redis
func InitializeRedis(redisParams config.RedisParams) error {
	RedisClient = redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%d", redisParams.Host, redisParams.Port), // адрес Redis-сервера
		Password: redisParams.Password,                                     // если пароль не установлен, оставьте пустым
		DB:       redisParams.DB,                                           // используемая база данных Redis
	})

	// Проверка соединения
	_, err := RedisClient.Ping(ctx).Result()
	if err != nil {
		log.Fatalf("Could not connect to Redis: %v", err)
		return err
	}

	return nil
}

func CloseRedisConnection() error {
	err := RedisClient.Close()
	if err != nil {
		return err
	}

	return nil
}
