package config

import (
	"context"
	"github.com/go-redis/redis/v8"
	"log"
)

type RedisConfig struct {
	Host     string `env:"HOST" envDefault:"localhost"`
	Port     string `env:"PORT" envDefault:"6379"`
	Password string `env:"PASSWORD" envDefault:""`
	Db       int    `env:"DB" envDefault:"0"`
}

var RedisClient *redis.Client

func InitRedis(addr, password string, db int, ctx context.Context) *redis.Client {
	client := redis.NewClient(&redis.Options{
		Addr:     addr,
		Password: password,
		DB:       db,
	})

	_, err := client.Ping(ctx).Result()
	if err != nil {
		log.Fatalf("Failed to connect to Redis: %v", err)
	}

	RedisClient = client
	return client
}
