package main

import (
	"context"
	"fmt"
	"log"
	"strconv"

	redis "github.com/redis/go-redis/v9"
)

var MainCtx context.Context
var RedisClient *redis.Client

func InitRedis(config ServerConfig) {
	MainCtx = context.Background()
	RedisClient = redis.NewClient(&redis.Options{Addr: config.RedisUrl})
	_, err := RedisClient.Ping(MainCtx).Result()
	if err != nil {
		log.Fatal(err)
	}

	log.Println("Redis: inited")
}

func GetState(userID int64) int16 {
	value, err := RedisClient.Get(MainCtx, fmt.Sprint(userID)).Result()
	if err != nil {
		log.Fatal(err)
	}

	state, _ := strconv.Atoi(value)
	return int16(state)
}

func SetState(userID int64, state int16) {
	err := RedisClient.Set(MainCtx, fmt.Sprint(userID), fmt.Sprint(state), 0).Err()
	if err != nil {
		log.Fatal(err)
	}
}
