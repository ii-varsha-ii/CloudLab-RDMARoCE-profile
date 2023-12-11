package utils

import (
	"math/rand"
	"os"
	"time"

	log "github.com/sirupsen/logrus"
)

const letterBytes = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"

func GetEnv(name string, panicIfEmpty bool) string {
	val := os.Getenv(name)
	if val == "" && panicIfEmpty {
		log.Fatalf("Env Var: %s is not set", name)
	}
	return val
}

func GetCurrentTime() time.Time {
	return time.Now().UTC()
}

func GetCurrentTimeStr() string {
	return time.Now().UTC().Format(time.RFC3339)
}

func ParseStrTime(strTime string) time.Time {
	timeObj, _ := time.Parse(time.RFC3339, strTime)
	return timeObj
}

func GenerateString(sizeInKB int) string {
	r := rand.New(rand.NewSource(time.Now().UnixNano()))

	totalBytes := sizeInKB * 1024

	randomBytes := make([]byte, totalBytes)

	for i := 0; i < totalBytes; i++ {
		randomBytes[i] = letterBytes[r.Intn(len(letterBytes))]
	}

	return string(randomBytes)
}

func GetMessageSizeInKB(message string) int {
	sizeInBytes := len(message)
	sizeInKB := sizeInBytes / 1024

	return sizeInKB
}
