package main

import (
	"fmt"
	"os"
	"time"

	"github.com/redis/go-redis/v9"
	log "github.com/sirupsen/logrus"
	"golang.org/x/net/context"
)

func getEnv(name string, panicIfEmpty bool) string {
	val := os.Getenv(name)
	if val == "" && panicIfEmpty {
		log.Fatalf("Env Var: %s is not set", name)
	}
	return val
}

func readFromRedis(ctx context.Context, client *redis.Client, key string) {

	// Infinite loop to continuously listen for changes
	lastMsg := ""
	for {
		val, err := client.Get(ctx, key).Result()
		if err == nil && val != lastMsg {
			lastMsg = val
			currentTime := time.Now()
			receivedTime, _ := time.Parse(time.RFC3339, val)
			diffInMs := currentTime.Sub(receivedTime)
			log.Infof("Received message: %s. Time taken: %vMs\n", val, diffInMs.Milliseconds())
		}
	}

}

func main() {
	redisMasterIP := getEnv("REDIS_MASTER_IP", true)
	redisMasterPort := getEnv("REDIS_MASTER_PORT", true)
	redisMasterPassword := getEnv("REDIS_MASTER_PASSWORD", false)
	redisKey := getEnv("REDIS_MASTER_KEY", true)
	ctx := context.Background()

	client := redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%s", redisMasterIP, redisMasterPort),
		Password: redisMasterPassword, // No password by default
		DB:       0,                   // Default DB
	})

	defer client.Close()

	log.Infof("Listening for key: %s", redisKey)
	readFromRedis(ctx, client, redisKey)
}
