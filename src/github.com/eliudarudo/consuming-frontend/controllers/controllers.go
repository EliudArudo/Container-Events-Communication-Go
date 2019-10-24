package controllers

import (
	"encoding/json"
	"net/http"

	"github.com/eliudarudo/consuming-frontend/dockerapi"
	"github.com/eliudarudo/consuming-frontend/logs"
	"github.com/eliudarudo/consuming-frontend/tasks"
)

var filename = "controllers/controllers.go"

// IndexController controller returns a test response on root route
func IndexController(w http.ResponseWriter, r *http.Request) {
	respondJSON(w, http.StatusOK, "OK")
}

// RequestRouteController receives requests and routes them to TaskController which determines the task
// and sends it to through redis pubsub
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

// RouterHandler404 covers all undefined routes
func RouterHandler404(w http.ResponseWriter, r *http.Request) {
	respondJSON(w, http.StatusOK, "Error")
}
