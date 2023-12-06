package api_handler

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/ii-varsha-ii/CloudLab-RDMARoCE-profile/container/data"
	"github.com/ii-varsha-ii/CloudLab-RDMARoCE-profile/container/sql_handler"
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
		log.Errorf("sendDataHandler: Exception while parsing request body: %v", err)
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}
	switch writeMessage.SourceType {
	case data.REDIS:
		{
			if err := c.BindJSON(&writeMessage); err != nil {
				log.Errorf("sendDataHandler: Exception while sending buld Redis Msg: %v", err)
				c.AbortWithError(http.StatusInternalServerError, err)
				return
			}
		}
	default:
		c.AbortWithError(http.StatusBadRequest, fmt.Errorf("unknown Source type: %d", writeMessage.SourceType))
	}

	c.IndentedJSON(http.StatusAccepted, writeMessage)
}

func InitializeAPIServer() {
	apiRouter = gin.Default()
	// Register the API endpoint handler
	apiRouter.GET("/getAllData", getAllDataHandler)
	apiRouter.POST("/writeData", writeDataHandler)
}

func StartAPIServer() {
	// Start the HTTP server
	log.Infof("StartAPIServer: API server listening on %s:%s", API_ADDRESS, API_PORT)
	err := apiRouter.Run(fmt.Sprintf("%s:%s", API_ADDRESS, API_PORT))
	if err != nil {
		log.Fatalf("Exception while launching API server: %v", err)
	}
}
