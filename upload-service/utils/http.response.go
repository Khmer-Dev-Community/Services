package utils

import (
	"encoding/json"
	"math"
	"net/http"

	"github.com/sirupsen/logrus"
)

type HttpMeta struct {
	Total    int `json:"total"`
	Page     int `json:"page"`
	LastPage int `json:"lastPage"`
}

type HttpPaginationResponse struct {
	Data    interface{} `json:"data"`
	Meta    Meta        `json:"meta"`
	Status  uint        `json:"statusCode"`
	Message string      `json:"message"`
}
type HttpReponseJson struct {
	Data    interface{} `json:"data"`
	Status  uint        `json:"statusCode"`
	Message string      `json:"message"`
}

func NewHttpPaginationResponse(w http.ResponseWriter, data interface{}, total, page, limit, status int, message string) HttpPaginationResponse {
	// Calculate last page
	lastPage := int(math.Ceil(float64(total) / float64(limit)))

	// Create response object
	response := HttpPaginationResponse{
		Data: data,
		Meta: Meta{
			Total:    total,
			Page:     page,
			LastPage: lastPage,
		},
		Status:  uint(status),
		Message: message,
	}

	// Set Content-Type header to application/json
	w.Header().Set("Content-Type", "application/json")

	// Set the status code based on the `status` argument
	w.WriteHeader(status)

	// Convert response to JSON and write it to the response body
	json.NewEncoder(w).Encode(response)

	// Return the response in case the caller needs it
	return response
}

func HttpErrorResponse(w http.ResponseWriter, errorMessage string) HttpPaginationResponse {

	return HttpPaginationResponse{
		Data: nil,
		Meta: Meta{
			Total:    0,
			Page:     0,
			LastPage: 0,
		},
		Status:  uint(http.StatusInternalServerError),
		Message: errorMessage,
	}
}

func HttpSuccessResponse(w http.ResponseWriter, data interface{}, status int, message string) HttpReponseJson {

	response := HttpReponseJson{
		Data:    data,
		Status:  uint(status),
		Message: message,
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(response)
	switch status {
	case http.StatusOK: // 200
		Logger(response, status, message, logrus.InfoLevel)
	case http.StatusBadRequest: // 400
		Logger(response, status, message, logrus.WarnLevel)
	case http.StatusInternalServerError: // 500
		Logger(response, status, message, logrus.ErrorLevel)
	default:
		Logger(response, status, message, logrus.DebugLevel)
	}
	return response
}
