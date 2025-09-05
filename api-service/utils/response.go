package utils

import (
	"encoding/json"
	"net/http"

	"github.com/gin-gonic/gin"
)

// Response represents the structure of the response
type Response struct {
	Data    interface{} `json:"data"`
	Meta    Meta        `json:"meta"`
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
		Data: data,

		Status:  status,
		Message: "Success",
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(response)
}

type ResponseData struct {
	Data    interface{} `json:"data"`
	Message string      `json:"message"`
	Status  uint        `json:"statusCode"`
	Meta    Meta        `json:"meta"`
}

// SuccessResponse sends a standardized success JSON response.
func SuccessResponse(c *gin.Context, status int, data interface{}, message string) {
	c.JSON(status, ResponseData{
		Data:    data,
		Message: message,
		Status:  uint(status),
	})
}
func SuccessResponsePage(c *gin.Context, status int, data interface{}, total int, page int, lastPage int, message string) {
	c.JSON(status, ResponseData{
		Data: data,
		Meta: Meta{
			Total:    total,
			Page:     page,
			LastPage: lastPage,
		},
		Message: message,
		Status:  uint(status),
	})
}
func ErrorResponse(c *gin.Context, status int, message string, details ...interface{}) {

	response := ResponseData{
		Data:    nil, // Data is typically null for error responses
		Message: message,
		Status:  uint(status),
	}
	if len(details) > 0 {
		response.Data = gin.H{"details": details[0]} // Pass details if available
	}
	c.JSON(status, response)
}
