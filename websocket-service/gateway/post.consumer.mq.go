package gateway

import (
	"encoding/json"
	"fmt"
	"log"
	"websocket-service/rabbitmq"
	"websocket-service/utils"

	socketio "github.com/googollee/go-socket.io"
)

// RabbitMQMessage represents the data structure of a message from RabbitMQ.
type RabbitMQMessage struct {
	EventType string          `json:"event_type"`
	Payload   json.RawMessage `json:"payload"`
}

// StartRabbitMQConsumer starts a consumer that listens for messages from RabbitMQ and broadcasts them via Socket.IO.
// It uses the pre-existing global RabbitMQ connection.
func StartRabbitMQConsumer(s *socketio.Server) error {
	// Check if the global RabbitMQ connection is available.
	if rabbitmq.RMQ == nil || rabbitmq.RMQ.Channel == nil {
		return fmt.Errorf("RabbitMQ connection is not initialized")
	}
	// Use the existing channel from the global RMQ instance.
	ch := rabbitmq.RMQ.Channel

	// Declare the fanout exchange that your publisher uses
	err := ch.ExchangeDeclare("post_events", "fanout", true, false, false, false, nil)
	if err != nil {
		return fmt.Errorf("failed to declare an exchange: %w", err)
	}

	// Declare a temporary queue and bind it to the exchange
	q, err := ch.QueueDeclare("", false, true, true, false, nil)
	if err != nil {
		return fmt.Errorf("failed to declare a queue: %w", err)
	}

	err = ch.QueueBind(q.Name, "", "post_events", false, nil)
	if err != nil {
		return fmt.Errorf("failed to bind queue to exchange: %w", err)
	}

	msgs, err := ch.Consume(q.Name, "", true, false, false, false, nil)
	if err != nil {
		return fmt.Errorf("failed to register a consumer: %w", err)
	}

	go func() {
		log.Println("RabbitMQ consumer started, waiting for messages...")

		for d := range msgs {
			var event RabbitMQMessage
			if err := json.Unmarshal(d.Body, &event); err != nil {
				log.Printf("Error unmarshalling RabbitMQ message: %v", err)
				continue
			}
			utils.InfoLog(event, "Event data ")
			// Decide where to broadcast the event based on its type
			switch event.EventType {
			case "post_created":
				// Example payload: {"post_id": "p123", "title": "Hello World"}
				var payload map[string]interface{}
				if err := json.Unmarshal(event.Payload, &payload); err != nil {
					log.Printf("Error unmarshalling post_created payload: %v", err)
					continue
				}
				// Broadcast to a general room for all posts
				s.BroadcastToRoom("/", "all_posts", "post_created", payload)

			case "reaction_added":
				// Example payload: {"post_id": "p123", "user_id": "u456", "type": "like"}
				var payload map[string]interface{}
				if err := json.Unmarshal(event.Payload, &payload); err != nil {
					log.Printf("Error unmarshalling reaction_added payload: %v", err)
					continue
				}
				// Extract the post_id to broadcast to the specific post's room
				if postID, ok := payload["post_id"].(string); ok {
					s.BroadcastToRoom("/", postID, "reaction_added", payload)
				}

			default:
				log.Printf("Received unknown event type: %s", event.EventType)
			}
		}
	}()

	return nil
}
