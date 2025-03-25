package redis

import (
	"github.com/go-redis/redis/v8"
	"log"
	"time"
)

// SetCache записывает данные в кэш с указанным сроком жизни
func SetCache(key string, value interface{}, expiration time.Duration) error {
	err := RedisClient.Set(ctx, key, value, expiration).Err()
	if err != nil {
		log.Printf("Error setting cache in Redis: %v", err)
		return err
	}
	return nil
}

// GetCache получает данные из кэша по ключу
func GetCache(key string) (string, error) {
	val, err := RedisClient.Get(ctx, key).Result()
	if err == redis.Nil {
		log.Printf("Key does not exist in Redis: %s", key)
		return "", nil
	} else if err != nil {
		log.Printf("Error getting cache from Redis: %v", err)
		return "", err
	}
	return val, nil
}

// DeleteCache удаляет данные из кэша по ключу
func DeleteCache(key string) error {
	err := RedisClient.Del(ctx, key).Err()
	if err != nil {
		log.Printf("Error deleting cache from Redis: %v", err)
		return err
	}
	return nil
}
