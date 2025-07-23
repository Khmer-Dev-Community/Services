package utils

import (
	"encoding/json"
	"net/http"
)

// Response represents the structure of the response
type Response struct {
	Data    interface{} `json:"data"`
	Status  int         `json:"statusCode"`
	Message string      `json:"message"`
}

// RespondWithError responds with an error message
func RespondWithError(w http.ResponseWriter, status int, message string) {
	response := Response{
		Data:    nil,
		Status:  status,
		Message: message,
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(response)
}

// RespondWithJSON responds with JSON data
func RespondWithJSON(w http.ResponseWriter, status int, data interface{}) {
	response := Response{
		Data:    data,
		Status:  status,
		Message: "Success",
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(response)
}
