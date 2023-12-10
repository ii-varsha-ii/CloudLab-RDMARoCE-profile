package api_handler

import (
	"context"
	"fmt"
	"github.com/ii-varsha-ii/CloudLab-RDMARoCE-profile/container/rpc_handler"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/ii-varsha-ii/CloudLab-RDMARoCE-profile/container/data"
	"github.com/ii-varsha-ii/CloudLab-RDMARoCE-profile/container/redis_handler"
	"github.com/ii-varsha-ii/CloudLab-RDMARoCE-profile/container/sql_handler"
	"github.com/ii-varsha-ii/CloudLab-RDMARoCE-profile/container/utils"
	log "github.com/sirupsen/logrus"
)

const (
	API_ADDRESS = "0.0.0.0"
	API_PORT    = "8000"
)

var (
	apiRouter *gin.Engine
)

// API Handlers
func getAllDataHandler(c *gin.Context) {
	messages, err := sql_handler.ReadAllData()
	if err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}
	c.IndentedJSON(http.StatusOK, messages)
}

// postAlbums adds an album from JSON received in the request body.
func writeDataHandler(c *gin.Context) {
	var writeMessage data.WriteMessage

	if err := c.BindJSON(&writeMessage); err != nil {
		log.Errorf("writeDataHandler: Exception while parsing request body: %v", err)
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}
	switch writeMessage.SourceType {
	case data.REDIS:
		{
			if err := redis_handler.RedisHandleBulkWrite(writeMessage); err != nil {
				log.Errorf("writeDataHandler: Exception while sending bulk Redis Msg: %v", err)
				c.AbortWithError(http.StatusInternalServerError, err)
				return
			}
		}
	case data.RDMA:
		{
			if err := redis_handler.RedisHandleBulkWrite(writeMessage); err != nil {
				log.Errorf("writeDataHandler: Exception while sending bulk Redis Msg: %v", err)
				c.AbortWithError(http.StatusInternalServerError, err)
				return
			}
		}
	case data.HTTP:
		{
			if err := APIHandleBulkWrite(writeMessage); err != nil {
				log.Errorf("writeDataHandler: Exception while sending bulk API Msg: %v", err)
				c.AbortWithError(http.StatusInternalServerError, err)
				return
			}
		}
	case data.RPC:
		{
			if err := rpc_handler.RPCHandleBulkWrite(context.Background(), writeMessage); err != nil {
				log.Errorf("writeDataHandler: Exception while sending bulk API Msg: %v", err)
				c.AbortWithError(http.StatusInternalServerError, err)
				return
			}
		}

	default:
		c.AbortWithError(http.StatusBadRequest, fmt.Errorf("unknown Source type: %d", writeMessage.SourceType))
	}

	c.IndentedJSON(http.StatusAccepted, writeMessage)
}

func receiveMessageHandler(c *gin.Context) {
	var apiMessage data.APIMessage

	if err := c.BindJSON(&apiMessage); err != nil {
		log.Errorf("receiveMessageHandler: Exception while parsing request body: %v", err)
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	currentTime := utils.GetCurrentTime()

	msg := &data.Message{
		SourceType:      data.HTTP,
		Message:         apiMessage.Message,
		MessageSizeInKB: utils.GetMessageSizeInKB(apiMessage.Message),
		WriteTime:       apiMessage.WriteTime,
		ReadTime:        currentTime,
		DiffInMs:        currentTime.Sub(apiMessage.WriteTime).Milliseconds(),
	}
	log.Infof("receiveMessageHandler: Received message: %s\n", msg.String())
	if err := sql_handler.RecordToDatabase(msg); err != nil {
		log.Errorf("receiveMessageHandler: Exceotion while writing Message: %s to SQL DB: %v", msg.String(), err)
		c.AbortWithError(http.StatusInternalServerError, err)
	}

	c.IndentedJSON(http.StatusAccepted, apiMessage)
}

func InitializeAPIServer() {
	apiRouter = gin.Default()
	// Register the API endpoint handler
	apiRouter.GET("/getAllData", getAllDataHandler)
	apiRouter.POST("/writeData", writeDataHandler)
	apiRouter.POST("/receiveMessage", receiveMessageHandler)
}

func StartAPIServer() {
	// Start the HTTP server
	log.Infof("StartAPIServer: API server listening on %s:%s", API_ADDRESS, API_PORT)
	err := apiRouter.Run(fmt.Sprintf("%s:%s", API_ADDRESS, API_PORT))
	if err != nil {
		log.Fatalf("Exception while launching API server: %v", err)
	}
}
