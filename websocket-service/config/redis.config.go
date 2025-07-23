package config

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/go-redis/redis/v8"
)

var RedisClient *redis.Client
var ctx = context.Background()

func InitRedis(address string, port string, password string) {
	// Initialize Redis client
	RedisClient = redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%s", address, port), // Using provided address and port
		Password: password,                            // Using provided password
		DB:       1,                                   // Default DB
	})

	// Ping Redis to test connection
	ctx := RedisClient.Context()
	pong, err := RedisClient.Ping(ctx).Result()
	if err != nil {
		fmt.Println("Failed to connect to Redis:", err)
		return
	}
	fmt.Println("Connected to Redis:", pong)
}

func Set(key, value string, expiration time.Duration) error {
	// Set key-value pair in Redis with expiration time
	ctx := RedisClient.Context()
	err := RedisClient.Set(ctx, key, value, expiration).Err()
	if err != nil {
		return err
	}
	return nil
}

func Get(key string) (string, error) {
	// Get value from Redis for the given key
	ctx := RedisClient.Context()
	val, err := RedisClient.Get(ctx, key).Result()
	if err != nil {
		return "", err
	}
	return val, nil
}
func SetWithExpiration(key string, value interface{}, expiration time.Duration) error {
	jsonData, err := json.Marshal(value)
	if err != nil {
		return err
	}

	return RedisClient.Set(ctx, key, jsonData, expiration).Err()
}

func UpdateExpiration(key string, expiration time.Duration) error {
	// Update expiration time of key in Redis
	ctx := RedisClient.Context()
	return RedisClient.Expire(ctx, key, expiration).Err()
}
func RemoveRedisKey(key string) error {
	_, err := RedisClient.Del(ctx, key).Result()
	if err != nil {
		return fmt.Errorf("could not delete key %s: %v", key, err)
	}
	return nil
}
