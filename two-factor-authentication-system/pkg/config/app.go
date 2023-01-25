package config

import (
	"context"

	"github.com/redis/go-redis/v9"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var db *gorm.DB
var redisClient *redis.Client
var ctx = context.Background()

func Connect() {
	dsn := "host=localhost user=postgres password=root dbname=tfas port=5433 sslmode=disable TimeZone=Asia/Dhaka"
	d, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})

	if err != nil {
		panic(d)
	}

	db = d

	client := redis.NewClient(&redis.Options{
		Addr:     "3.36.21.109:6379",
		Password: "123456",
		DB:       0,
	})
	redisClient = client
}

func GetDB() *gorm.DB {
	return db
}

func GetRedisClient() *redis.Client {
	return redisClient
}

func GetContext() context.Context{
	return ctx
}
