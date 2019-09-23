package controllers

import (
	"encoding/json"
	"net/http"
)

// type indexStruct struct {
// 	Path  string `json:"path"`
// 	Route string `json:"route"`
// }

// HTTPIndex controller receives GET requests from '/' route
func HTTPIndex(w http.ResponseWriter, r *http.Request) {
	respondJSON(w, http.StatusOK, "OK")
}

// HTTPPostTaskHandler gets an Item from our store using route parameters
func HTTPPostTaskHandler(w http.ResponseWriter, r *http.Request) {
	respondJSON(w, http.StatusOK, "OK")

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

// RespondError to cover all unhandled routes
func RespondError(w http.ResponseWriter, r *http.Request) {
	respondJSON(w, http.StatusOK, "Error")
}
