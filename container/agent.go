package main

import (
	"github.com/ii-varsha-ii/CloudLab-RDMARoCE-profile/container/api_handler"
	"github.com/ii-varsha-ii/CloudLab-RDMARoCE-profile/container/redis_handler"
	"github.com/ii-varsha-ii/CloudLab-RDMARoCE-profile/container/rpc_handler"
	"github.com/ii-varsha-ii/CloudLab-RDMARoCE-profile/container/sql_handler"
	"github.com/ii-varsha-ii/CloudLab-RDMARoCE-profile/container/utils"
	log "github.com/sirupsen/logrus"
	"golang.org/x/net/context"
)

var (
	redisMasterIP       string
	redisMasterPort     string
	redisMasterPassword string
	redisWriteKey       string
	redisListenKey      string
	sqlDBUser           string
	sqlDBPassword       string
	sqlDBName           string
	sqlDBHost           string
	sqlDBPort           string
	ctx                 context.Context
	cancel              context.CancelFunc
)

func initialize() {
	redisMasterIP = utils.GetEnv("REDIS_MASTER_IP", true)
	redisMasterPort = utils.GetEnv("REDIS_MASTER_PORT", true)
	redisMasterPassword = utils.GetEnv("REDIS_MASTER_PASSWORD", false)
	redisWriteKey = utils.GetEnv("REDIS_WRITE_KEY", true)
	redisListenKey = utils.GetEnv("REDIS_LISTEN_KEY", true)

	sqlDBUser = utils.GetEnv("SQL_DB_USER", true)
	sqlDBPassword = utils.GetEnv("SQL_DB_PASSWORD", true)
	sqlDBName = utils.GetEnv("SQL_DB_NAME", true)
	sqlDBHost = utils.GetEnv("SQL_DB_HOST", true)
	sqlDBPort = utils.GetEnv("SQL_DB_PORT", true)

	ctx = context.Background()

	if err := redis_handler.InitializeClient(ctx, redisMasterIP, redisMasterPort, redisMasterPassword, redisWriteKey, redisListenKey, 0); err != nil {
		log.Fatalf("Exception while initializing Redis client: %v", err)
	}

	if err := sql_handler.InitializeSQLClient(sqlDBUser, sqlDBPassword, sqlDBHost, sqlDBPort, sqlDBName); err != nil {
		log.Fatalf("Exception while initializing SQL client: %v", err)
	}

	if err := sql_handler.CreateTable(); err != nil {
		log.Fatalf("Exception while creating Table with SQL client: %v", err)
	}

	api_handler.InitializeAPIServer()

	go redis_handler.ReadFromRedisContinuously(ctx)
	go rpc_handler.StartRPCServer()
}

func close() {
	log.Infof("Closing Clients")
	//cancel()
	redis_handler.Close()
	sql_handler.Close()
}

func main() {
	initialize()
	defer close()
	api_handler.StartAPIServer()
}
