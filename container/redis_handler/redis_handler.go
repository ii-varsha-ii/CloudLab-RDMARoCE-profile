package redis_handler

import (
	"context"
	"fmt"
	"github.com/ii-varsha-ii/CloudLab-RDMARoCE-profile/container/data"
	"github.com/ii-varsha-ii/CloudLab-RDMARoCE-profile/container/sql_handler"
	"github.com/ii-varsha-ii/CloudLab-RDMARoCE-profile/container/utils"
	"strings"

	"github.com/redis/go-redis/v9"
	log "github.com/sirupsen/logrus"
)

var (
	redisClient    *redis.Client
	redisWriteKey  string
	redisListenKey string
	err            error
)

func InitializeClient(ctx context.Context, redisMasterIP, redisMasterPort, redisMasterPassword, writeKey, listenKey string, dbName int) error {
	redisClient = redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%s", redisMasterIP, redisMasterPort),
		Password: redisMasterPassword, // No password by default
		DB:       dbName,              // Default DB
	})
	if err := redisClient.Ping(ctx).Err(); err != nil {
		log.Errorf("initializeClient: Exception while checking Redis connection: %v", err)
		return err
	}
	redisWriteKey = writeKey
	redisListenKey = listenKey
	return nil
}

func Close() error {
	return redisClient.Close()
}

func WriteToRedis(value string) error {
	err := redisClient.Set(context.Background(), redisWriteKey, value, 0).Err()
	if err != nil {
		log.Errorf("writeToRedis: Exception while writing to key '%s' with value '%s'. %v\n", redisWriteKey, value, err)
		return err
	}

	log.Infof("writeToRedis: Successfully wrote to key '%s' with value '%s'\n", redisWriteKey, value)
	return nil
}

func ReadFromRedisContinuously(ctx context.Context) {

	// Infinite loop to continuously listen for changes. Only handle new messages.
	lastMsg, _ := redisClient.Get(ctx, redisListenKey).Result()
	for {
		val, err := redisClient.Get(ctx, redisListenKey).Result()
		if err == nil && val != lastMsg {
			lastMsg = val
			currentTime := utils.GetCurrentTime()

			redisMsg := parseRedisMessage(val)

			msg := &data.Message{
				SourceType:      data.REDIS,
				Message:         redisMsg.Message,
				MessageSizeInKB: utils.GetMessageSizeInKB(redisMsg.Message),
				WriteTime:       redisMsg.WriteTime,
				ReadTime:        currentTime,
				DiffInMs:        currentTime.Sub(redisMsg.WriteTime).Milliseconds(),
			}
			log.Infof("ReadFromRedisContinuously: Received message: %s\n", msg.String())
			if err := sql_handler.RecordToDatabase(msg); err != nil {
				log.Errorf("ReadFromRedisContinuously: Exceotion while writing Message: %s to SQL DB: %v", msg.String(), err)
			}
			log.Infof("ReadFromRedisContinuously: Successfully wrote message: %s to SQL DB\n", msg.String())
		}
	}
}

func HandleBulkWrite(writeMessage data.WriteMessage) error {
	msg := utils.GenerateString(writeMessage.MessageSizeInKB)
	for i := 0; i < writeMessage.MessageCount; i++ {
		redisMsg := fmt.Sprintf("%s:%s", utils.GetCurrentTimeStr(), msg)
		if err := WriteToRedis(redisMsg); err != nil {
			log.Errorf("HandleBulkWrite: Exception while writing message to redis: %v", err)
			return err
		}
	}
	return nil
}

func parseRedisMessage(message string) data.RedisMessage {
	vals := strings.Split(message, ":")
	redisMsg := data.RedisMessage{
		Message:   vals[1],
		WriteTime: utils.ParseStrTime(vals[0]),
	}
	return redisMsg
}
