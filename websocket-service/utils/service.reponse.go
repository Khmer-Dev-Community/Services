package utils

import (
	"math"

	"github.com/sirupsen/logrus"
)

type PaginationRequestDTO struct {
	Page  int    `json:"page"`
	Limit int    `json:"limit"`
	Query string `json:"query"`
}
type ServiceReponseObject struct {
	Data    interface{} `json:"data"`
	Status  uint        `json:"statusCode"`
	Message string      `json:"message"`
}

type ServicePaginationResponse struct {
	Data    interface{} `json:"data"`
	Meta    Meta        `json:"meta"`
	Status  uint        `json:"statusCode"`
	Message string      `json:"message"`
}

func NewServicePaginationResponse(data interface{}, total, page, limit int, status int, message string, logLevel logrus.Level, service string) ServicePaginationResponse {
	lastPage := int(math.Ceil(float64(total) / float64(limit)))
	LoggerService(data, "Response", logLevel)

	return ServicePaginationResponse{
		Data: data,
		Meta: Meta{
			Total:    total,
			Page:     page,
			LastPage: lastPage,
		},
		Status:  uint(status),
		Message: message,
	}
}
func ServiceResponse(data interface{}, status int, message string, logLevel logrus.Level) ServiceReponseObject {
	// Create a structured log entry
	logEntry := logrus.WithFields(logrus.Fields{
		"data":    data,
		"status":  status,
		"message": message,
	})
	switch logLevel {
	case logrus.ErrorLevel:
		logEntry.Error(message) // Log the message parameter dynamically for errors
	case logrus.InfoLevel:
		logEntry.Info(message) // Log the message parameter dynamically for info
	default:
		logEntry.Info(message) // Default log level will be Info, with dynamic message
	}

	// Return the service response object
	return ServiceReponseObject{
		Data:    data,
		Status:  uint(status),
		Message: message,
	}
}
