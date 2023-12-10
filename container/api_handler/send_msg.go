package api_handler

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/ii-varsha-ii/CloudLab-RDMARoCE-profile/container/data"
	"github.com/ii-varsha-ii/CloudLab-RDMARoCE-profile/container/utils"
	log "github.com/sirupsen/logrus"
)

func APIHandleBulkWrite(writeMessage data.WriteMessage) error {
	serverHost := writeMessage.APIHostName
	serverPort := writeMessage.APIPort
	apiEndpoint := "/receiveMessage"

	msg := utils.GenerateString(writeMessage.MessageSizeInKB)

	apiURL := fmt.Sprintf("http://%s:%s%s", serverHost, serverPort, apiEndpoint)

	for i := 0; i < writeMessage.MessageCount; i++ {
		requestPayload := data.APIMessage{
			Message:   msg,
			WriteTime: utils.GetCurrentTime(),
		}

		requestBody, err := json.Marshal(requestPayload)
		if err != nil {
			log.Errorf("APIHandleBulkWrite: exception while marshalling json: %v", err)
			return err
		}
		response, err := http.Post(apiURL, "application/json", bytes.NewBuffer(requestBody))
		if err != nil {
			log.Errorf("APIHandleBulkWrite: exception while sending POST request to %s: %v", apiURL, err)
			return err
		}

		if response.StatusCode == http.StatusAccepted {
			log.Infof("APIHandleBulkWrite: Request successful. Status code: %d", response.StatusCode)
			response.Body.Close()
		} else {
			err := fmt.Errorf("exception in POST request's response to %s. Status code: %d", apiURL, response.StatusCode)
			log.Errorf("APIHandleBulkWrite: %v", err)
			response.Body.Close()
			return err
		}
		time.Sleep(1 * time.Second)
	}
	return nil
}
