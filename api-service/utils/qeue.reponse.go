package utils

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/streadway/amqp"
)

// SendSuccessResponse sends a success response with the provided data.
func SendSuccessResponse(data interface{}, replyTo string, channel *amqp.Channel) {
	response := map[string]interface{}{
		"data":       data,
		"statusCode": http.StatusOK,
	}
	sendResponse(response, replyTo, channel)
}

// SendErrorResponse sends an error response with the provided error message.
func SendErrorResponse(errorMessage string, replyTo string, channel *amqp.Channel) {
	response := map[string]interface{}{
		"error":      errorMessage,
		"statusCode": "error",
	}
	sendResponse(response, replyTo, channel)
}

// SendPaginationResponse sends a response containing pagination data.
func SendPaginationResponse(data interface{}, total, page, limit int, replyTo string, channel *amqp.Channel) {
	meta := map[string]int{
		"total":    total,
		"page":     page,
		"lastPage": calculateLastPage(total, limit),
	}
	response := map[string]interface{}{
		"data":       data,
		"meta":       meta,
		"statusCode": http.StatusOK,
	}
	sendResponse(response, replyTo, channel)
}

// sendResponse sends the response message.
func sendResponse(response map[string]interface{}, replyTo string, channel *amqp.Channel) {
	jsonResponse, err := json.Marshal(response)
	if err != nil {
		log.Printf("Failed to marshal JSON: %v", err)
		return
	}
	err = channel.Publish(
		"",
		replyTo,
		false,
		false,
		amqp.Publishing{
			ContentType: "application/json",
			Body:        jsonResponse,
		},
	)
	if err != nil {
		log.Printf("Failed to publish a message: %v", err)
	}
}

// calculateLastPage calculates the last page based on total records and limit per page.
func calculateLastPage(total, limit int) int {
	if total == 0 {
		return 1
	}
	if total%limit == 0 {
		return total / limit
	}
	return (total / limit) + 1
}
