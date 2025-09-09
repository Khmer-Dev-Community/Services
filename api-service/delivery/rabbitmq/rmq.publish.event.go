package rabbitmq

import (
	"encoding/json"
	"fmt"
	"log"
	"time"

	"github.com/google/uuid"
	"github.com/streadway/amqp"
)

// EventPublisher is a struct for publishing messages
type EventPublisher struct{}
type RabbitMQMessage struct {
	EventType string          `json:"event_type"`
	Payload   json.RawMessage `json:"payload"`
}

// NewEventPublisher creates a new publisher instance
func NewEventPublisher() *EventPublisher {
	return &EventPublisher{}
}
func (p *EventPublisher) PublishEvent(eventType string, exchange string, payload interface{}) error {
	if RMQ == nil {
		return fmt.Errorf("RabbitMQ connection is not initialized")
	}

	body, err := json.Marshal(payload)
	if err != nil {
		return fmt.Errorf("failed to marshal event to JSON: %w", err)
	}

	// Declare the exchange (safe to call multiple times)
	if err := RMQ.DeclareExchange(exchange, "direct"); err != nil {
		return fmt.Errorf("failed to declare exchange: %w", err)
	}

	// Publish the message to the exchange with the eventType as the routing key
	if err := RMQ.PublishMessage(exchange, eventType, amqp.Publishing{
		ContentType: "application/json",
		Body:        body,
	}); err != nil {
		return fmt.Errorf("failed to publish message: %w", err)
	}

	log.Printf("Successfully published event '%s' to exchange '%s'", eventType, exchange)
	return nil
}

// Corrected: The GetNewChannel method belongs to the RabbitMQ struct.
// In your rabbitmq package, in the file that contains the RabbitMQ struct.
func (rmq *RabbitMQ) GetNewChannel() (*amqp.Channel, error) {
	const maxRetries = 15 // Increased retries for resilience
	for i := 0; i < maxRetries; i++ {
		if rmq.Connection == nil || rmq.Connection.IsClosed() {
			log.Println("RabbitMQ connection is not open. Waiting for reconnect...")
			time.Sleep(1 * time.Second)
			continue
		}
		newChannel, err := rmq.Connection.Channel()
		if err == nil {
			return newChannel, nil // Success
		}
		log.Printf("Failed to get new channel: %v. Retrying...", err)
		time.Sleep(1 * time.Second)
	}
	return nil, fmt.Errorf("failed to get a new channel after %d retries", maxRetries)
}

func PublishBroadcastEvent(eventType string, exchangeName string, payload interface{}) {
	// 1. Get a new channel for this specific request
	tempChannel, err := RMQ.GetNewChannel()
	if err != nil {
		log.Printf("Failed to get a new channel: %v", err)
		return
	}
	defer tempChannel.Close()

	// 2. Marshal the payload to a JSON byte slice
	payloadBody, err := json.Marshal(payload)
	if err != nil {
		log.Printf("Failed to marshal payload to JSON: %v", err)
		return
	}

	// 3. Create the complete event message struct (what the consumer expects)
	eventMessage := RabbitMQMessage{
		EventType: eventType,
		Payload:   payloadBody,
	}

	// 4. Marshal the complete event message struct
	body, err := json.Marshal(eventMessage)
	if err != nil {
		log.Printf("Failed to marshal event message to JSON: %v", err)
		return
	}

	// 5. Declare the fanout exchange (safe to call multiple times)
	err = tempChannel.ExchangeDeclare(
		exchangeName,
		"fanout", // Use fanout for broadcasting
		true,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		log.Printf("Failed to declare an exchange: %v", err)
		return
	}

	// 6. Publish the message with an empty routing key
	err = tempChannel.Publish(
		exchangeName,
		"", // Routing key is ignored by fanout exchanges
		false,
		false,
		amqp.Publishing{
			ContentType: "application/json",
			Body:        body,
		},
	)
	if err != nil {
		log.Printf("Failed to publish a message: %v", err)
		return
	}

	log.Printf("Message sent to RabbitMQ for fanout exchange: %s", exchangeName)
}

func HandleRequestReply(exchangeName string, routingKey string, payload interface{}) ([]byte, error) {
	// 1. Get a new channel for this specific request
	tempChannel, err := RMQ.GetNewChannel()
	if err != nil {
		return nil, fmt.Errorf("failed to get a new channel: %w", err)
	}
	defer tempChannel.Close()

	// 2. Declare the exchange on the temporary channel
	err = tempChannel.ExchangeDeclare(
		exchangeName,
		"direct", // Use direct for request-reply
		true,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to declare an exchange: %w", err)
	}

	// 3. Declare the reply queue on the temporary channel
	replyQueue, err := tempChannel.QueueDeclare(
		"",
		false,
		false,
		true,
		false,
		nil,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to declare a reply queue: %w", err)
	}
	defer tempChannel.QueueDelete(replyQueue.Name, false, false, false)

	// 4. Marshal the payload to a JSON byte slice
	body, err := json.Marshal(payload)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal payload to JSON: %w", err)
	}

	// 5. Publish the message using the temporary channel
	// You must also include a CorrelationId for the reply
	corrID := uuid.New().String()
	err = tempChannel.Publish(
		exchangeName,
		routingKey,
		false,
		false,
		amqp.Publishing{
			ContentType:   "application/json",
			Body:          body,
			ReplyTo:       replyQueue.Name,
			CorrelationId: corrID,
		},
	)
	if err != nil {
		return nil, fmt.Errorf("failed to publish a message: %w", err)
	}

	// 6. Consume the reply
	msgs, err := tempChannel.Consume(
		replyQueue.Name,
		"",
		true,
		true,
		false,
		false,
		nil,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to consume messages: %w", err)
	}

	// 7. Wait for the reply with a timeout
	timeout := time.After(5 * time.Second)
	select {
	case msg := <-msgs:
		// Check if the reply is for this request
		if msg.CorrelationId == corrID {
			return msg.Body, nil
		}
	case <-timeout:
		return nil, fmt.Errorf("request timed out")
	}

	return nil, fmt.Errorf("no response received")
}
