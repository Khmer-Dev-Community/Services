package utils

import (
	"encoding/json"
	"log"
	"math"
	"net/http"

	"github.com/streadway/amqp"
)

// SendExchageResponse publishes a response message to RabbitMQ.
func SendExchnageResponse(response map[string]interface{}, replyTo string, channel *amqp.Channel) {
	body, err := json.Marshal(response)
	if err != nil {
		log.Printf("Failed to marshal response: %v", err)
		return
	}

	err = channel.Publish(
		"",      // default exchange
		replyTo, // routing key (reply-to queue)
		false,   // mandatory
		false,   // immediate
		amqp.Publishing{
			ContentType: "application/json",
			Body:        body,
		},
	)
	if err != nil {
		log.Printf("Failed to publish response: %v", err)
	}
}

func SendExchangeResponseMaping(response interface{}, replyTo string, channel *amqp.Channel) {
	body, err := json.Marshal(response)
	if err != nil {
		log.Printf("Failed to marshal response: %v", err)
		return
	}

	err = channel.Publish(
		"",      // default exchange
		replyTo, // routing key (reply-to queue)
		false,   // mandatory
		false,   // immediate
		amqp.Publishing{
			ContentType: "application/json",
			Body:        body,
		},
	)
	if err != nil {
		log.Printf("Failed to publish response: %v", err)
	}
}
func SendResponseExchangeMappingDataPagination(data interface{}, total, page, limit, status int, message string, replyTo string, channel *amqp.Channel) {
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

	SendExchangeResponseMaping(response, replyTo, channel)

}
func SendExchangeSuccessResponse(data interface{}, replyTo string, channel *amqp.Channel) {
	response := map[string]interface{}{
		"data":       data,
		"statusCode": http.StatusOK,
	}
	SendExchnageResponse(response, replyTo, channel)
}
func SendExchageResponseQuene(data interface{}, status int, message string, replyTo string, channel *amqp.Channel) {
	response := map[string]interface{}{
		"data":       data,
		"statusCode": status,
		"message":    message,
	}
	SendExchnageResponse(response, replyTo, channel)
}

func SendExchnageErrorResponse(errorMessage string, replyTo string, channel *amqp.Channel) {
	response := map[string]interface{}{
		"error":      errorMessage,
		"statusCode": http.StatusInternalServerError,
	}
	SendExchnageResponse(response, replyTo, channel)
}
