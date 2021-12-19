package db

import (
	"fmt"
	"time"

	"github.com/go-redis/redis"
)

type DBService struct {
	redisClient *redis.Client
}

var (
	dbService = &DBService{}
)

func InitializeStore() *DBService {
	redisClient := redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "",
		DB:       0,
	})

	pong, err := redisClient.Ping().Result()
	if err != nil {
		panic(fmt.Sprintf("Error init Redis: %v", err))
	}

	fmt.Printf("\nRedis started successfully: pong message = {%s}", pong)
	dbService.redisClient = redisClient
	return dbService
}

func SaveUrlMapping(shortUrl string, urlMapping UrlMapping, ExpireTime time.Duration) error {
	if dbService.redisClient == nil {
		return fmt.Errorf("redis client is not initialized")
	}
	return dbService.redisClient.Set(shortUrl, urlMapping, ExpireTime).Err()
}

func RetrieveInitialUrl(shortUrl string) string {
	result, err := dbService.redisClient.Get(shortUrl).Result()
	if err != nil {
		panic(err)
	}
	urlMapping := &UrlMapping{}
	urlMapping.UnmarshalBinary([]byte(result))
	if err != nil {
		panic(fmt.Sprintf("Failed RetrieveInitialUrl url | Error: %v - shortUrl: %s\n", err, shortUrl))
	}
	urlMapping.Count++
	dbService.redisClient.Set(shortUrl, urlMapping, 0)
	return urlMapping.Original_url
}