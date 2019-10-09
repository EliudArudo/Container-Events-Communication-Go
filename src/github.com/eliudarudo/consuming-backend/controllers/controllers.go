package controllers

import (
	"encoding/json"
	"net/http"

	"github.com/eliudarudo/consuming-backend/dockerapi"
	"github.com/eliudarudo/consuming-backend/logs"
	"github.com/eliudarudo/consuming-backend/tasks"
)

var filename = "controllers/controllers.go"

// IndexController controller receives GET requests from '/' route
func IndexController(w http.ResponseWriter, r *http.Request) {
	respondJSON(w, http.StatusOK, "OK")
}

// RequestRouteController gets an Item from our store using route parameters
func RequestRouteController(w http.ResponseWriter, r *http.Request) {

	myContainerInfo := dockerapi.GetMyOfflineContainerInfo()

	var decodedRequestBody map[string]interface{}
	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&decodedRequestBody); err != nil {
		logs.StatusFileMessageLogging("FAILURE", filename, "RequestRouteController", err.Error())
	}
	defer r.Body.Close()

	response, err := tasks.TaskController(decodedRequestBody, myContainerInfo)
	if err != nil {
		respondError(w, http.StatusBadGateway, "Server Error")
		return
	}

	respondJSON(w, http.StatusOK, response)
}

// Utility functions
func respondJSON(w http.ResponseWriter, status int, payload interface{}) {
	response, err := json.Marshal(payload)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	w.Write([]byte(response))
}

func respondError(w http.ResponseWriter, code int, message string) {
	respondJSON(w, code, map[string]string{"error": message})
}

// RouterHandler404 to cover all unhandled routes
func RouterHandler404(w http.ResponseWriter, r *http.Request) {
	respondJSON(w, http.StatusOK, "Error")
}

// // RedisController receives messages sent through redis lines
// func RedisController(event *redis.Message, containerInfo interfaces.ContainerInfoStruct) {
// 	logic.EventDeterminer(event.Payload, containerInfo)
// }
