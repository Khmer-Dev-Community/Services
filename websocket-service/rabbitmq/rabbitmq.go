package rabbitmq

import (
	"fmt"
	"log"

	"github.com/streadway/amqp"
)

// RabbitMQ struct to hold connection and channel
type RabbitMQ struct {
	Connection *amqp.Connection
	Channel    *amqp.Channel
}

var RMQ *RabbitMQ // Package-level variable to hold RabbitMQ connection

// Controller interface for handling messages
type Controller interface {
	HandleMessage(msg amqp.Delivery, channel *amqp.Channel) error
}

// InitializeRabbitMQ initializes RabbitMQ connection
func InitializeRabbitMQ(url string) error {
	conn, err := amqp.Dial(url)
	if err != nil {
		return err
	}

	channel, err := conn.Channel()
	if err != nil {
		conn.Close()
		return err
	}

	RMQ = &RabbitMQ{
		Connection: conn,
		Channel:    channel,
	}
	fmt.Println("RabbitMQ Connected")
	return nil
}

// DeclareQueue declares a queue
func (rmq *RabbitMQ) DeclareQueue(queueName string) error {
	_, err := rmq.Channel.QueueDeclare(
		queueName, // name
		false,     // durable
		false,     // delete when unused
		false,     // exclusive
		false,     // no-wait
		nil,       // arguments
	)
	return err
}

// CreateReplyQueue creates a reply queue
func (rmq *RabbitMQ) CreateReplyQueue() (amqp.Queue, error) {
	return rmq.Channel.QueueDeclare(
		"",    // name
		false, // durable
		false, // delete when unused
		true,  // exclusive
		false, // no-wait
		nil,   // arguments
	)
}

// PublishMessage publishes a message to RabbitMQ
func (rmq *RabbitMQ) PublishMessage(exchange, routingKey string, msg amqp.Publishing) error {
	return rmq.Channel.Publish(
		exchange,   // exchange
		routingKey, // routing key
		false,      // mandatory
		false,      // immediate
		msg,        // message
	)
}

// ConsumeMessages consumes messages from a queue
func (rmq *RabbitMQ) ConsumeMessages(queue string) (<-chan amqp.Delivery, error) {
	return rmq.Channel.Consume(
		queue, // queue
		"",    // consumer
		false, // auto-ack
		false, // exclusive
		false, // no-local
		false, // no-wait
		nil,   // args
	)
}

// DeleteQueue deletes a queue
func (rmq *RabbitMQ) DeleteQueue(queueName string) error {
	_, err := rmq.Channel.QueueDelete(
		queueName, // name
		false,     // ifUnused
		false,     // ifEmpty
		false,     // noWait
	)
	return err
}

// Close closes the RabbitMQ connection and channel
func (rmq *RabbitMQ) Close() error {
	if rmq.Channel != nil {
		if err := rmq.Channel.Close(); err != nil {
			return err
		}
	}
	if rmq.Connection != nil {
		if err := rmq.Connection.Close(); err != nil {
			return err
		}
	}
	return nil
}

// DeclareExchange declares an exchange in RabbitMQ.
func (rmq *RabbitMQ) DeclareExchange(exchangeName, exchangeType string) error {
	err := rmq.Channel.ExchangeDeclare(
		exchangeName, // name of the exchange
		exchangeType, // type of the exchange (e.g., "direct", "fanout", "topic")
		true,         // durable
		false,        // auto-deleted
		false,        // internal
		false,        // no-wait
		nil,          // arguments
	)
	if err != nil {
		return fmt.Errorf("failed to declare exchange: %w", err)
	}
	return nil
}

// StartConsumer starts consuming messages from the specified queue
func (rmq *RabbitMQ) StartConsumer(queueName string, controller Controller) error {
	msgs, err := rmq.ConsumeMessages(queueName)
	if err != nil {
		return fmt.Errorf("failed to consume messages: %w", err)
	}

	go func() {
		for msg := range msgs {
			if err := controller.HandleMessage(msg, rmq.Channel); err != nil {
				log.Printf("Error handling message: %s, error: %v", msg.Body, err)
			} else {
				msg.Ack(false) // Acknowledge the message after successful processing
			}
		}
	}()

	return nil
}

// ExchangeDeclareAndHandleQueues declares an exchange and handles multiple queues.
func ExchangeDeclareAndHandleQueues(channel *amqp.Channel, exchangeName, exchangeType string, queueNames []string, controller Controller) error {
	// Declare the exchange
	err := channel.ExchangeDeclare(
		exchangeName, // name
		exchangeType, // type
		true,         // durable
		false,        // auto-deleted
		false,        // internal
		false,        // no-wait
		nil,          // arguments
	)
	if err != nil {
		return fmt.Errorf("failed to declare exchange: %w", err)
	}

	for _, queueName := range queueNames {
		// Declare each queue
		_, err := channel.QueueDeclare(
			queueName, // name
			false,     // durable
			false,     // delete when unused
			false,     // exclusive
			false,     // no-wait
			nil,       // arguments
		)
		if err != nil {
			return fmt.Errorf("failed to declare queue %s: %w", queueName, err)
		}

		// Bind the queue to the exchange with its own routing key
		err = channel.QueueBind(
			queueName,    // queue name
			queueName,    // routing key
			exchangeName, // exchange name
			false,        // no-wait
			nil,          // arguments
		)
		if err != nil {
			return fmt.Errorf("failed to bind queue %s: %w", queueName, err)
		}

		// Start consuming messages from the queue
		msgs, err := channel.Consume(
			queueName, // queue
			"",        // consumer
			true,      // auto-ack
			false,     // exclusive
			false,     // no-local
			false,     // no-wait
			nil,       // args
		)
		if err != nil {
			return fmt.Errorf("failed to register consumer for queue %s: %w", queueName, err)
		}

		// Handle messages in a separate goroutine
		go func(qName string) {
			for msg := range msgs {
				if err := controller.HandleMessage(msg, channel); err != nil {
					// Log the error if the message handling fails
					fmt.Printf("Error handling message from queue %s: %v\n", qName, err)
				}
			}
		}(queueName)
	}

	return nil
}
