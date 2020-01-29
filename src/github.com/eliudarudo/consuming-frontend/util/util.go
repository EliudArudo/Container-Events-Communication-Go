package util

import (
	"github.com/eliudarudo/consuming-frontend/interfaces"
)

var filename = "util/util.go"

// ResponseBuffer is an array of events that have arrived
var ResponseBuffer []string

var responses []interfaces.ReceivedEventInterface

// PushResponseToBuffers pushes response object and ID to ResponseBuffer and responses array respectively
func PushResponseToBuffers(response *(interfaces.ReceivedEventInterface)) {
	ResponseBuffer = append(ResponseBuffer, (*response).RequestID)
	responses = append(responses, *response)
}

func clearResponseFromBuffers(requestID string) {
	newResponseBuffer := make([]string, 0)
	newResponsesArray := make([]interfaces.ReceivedEventInterface, 0)
	for _, ID := range ResponseBuffer {
		isResponseID := ID != requestID
		if isResponseID {
			newResponseBuffer = append(newResponseBuffer, ID)
		}
	}

	for _, response := range responses {
		isResponse := response.RequestID != requestID
		if isResponse {
			newResponsesArray = append(newResponsesArray, response)
		}
	}

	ResponseBuffer = newResponseBuffer
	responses = newResponsesArray
}

// GetResponseFromBuffer checks and retrieves the response when it's delivered from redis pubsub
func GetResponseFromBuffer(requestID string) *(interfaces.ReceivedEventInterface) {
	var responseArrived bool
	var response interfaces.ReceivedEventInterface
	for _, ID := range ResponseBuffer {
		responseArrived = (ID == requestID)
	}

	if responseArrived {
		for _, bufferResponse := range responses {
			if bufferResponse.RequestID == requestID {
				response = bufferResponse
			}
		}
	}

	return &response
}
