package rabbitmq

import (
	"fmt"
	"log"

	"github.com/sirupsen/logrus"
	"github.com/streadway/amqp"
)

// RabbitMQ struct to hold connection and channel
type RabbitMQ struct {
	Connection *amqp.Connection
	Channel    *amqp.Channel
}

var RMQ *RabbitMQ // Package-level variable to hold RabbitMQ connection
type Controller interface {
	HandleMessage(msg amqp.Delivery, channel *amqp.Channel)
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
	fmt.Println(" RabbitMQ Connected")
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
		true,  // auto-ack
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
	if rmq.Connection != nil {
		if err := rmq.Connection.Close(); err != nil {
			return err
		}
	}
	if rmq.Channel != nil {
		if err := rmq.Channel.Close(); err != nil {
			return err
		}
	}
	return nil
}

// DeclareExchange declares an exchange in RabbitMQ.
func (rmq *RabbitMQ) DeclareExchange_(exchangeName, exchangeType string) error {
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

func (rmq *RabbitMQ) DeclareExchange(name, kind string) error {
	err := rmq.Channel.ExchangeDeclare(
		name,  // name of the exchange
		kind,  // type of exchange (e.g., "direct", "fanout", "topic", "headers")
		true,  // durable
		false, // auto-delete
		false, // internal
		false, // no-wait
		nil,   // arguments
	)
	if err != nil {
		log.Printf("Failed to declare an exchange: %v", err)
		return err
	}
	return nil
}
func ExchnageDeclareAndHandleQueues(
	channel *amqp.Channel,
	exchangeName string,
	exchangeType string,
	queueNames []string,
	controller Controller,
) {
	logger := logrus.New()
	logger.SetFormatter(&logrus.JSONFormatter{})

	// Declare the exchange
	err := channel.ExchangeDeclare(
		exchangeName, // name
		exchangeType, // type (change to your required exchange type, e.g., "direct", "topic", "fanout")
		true,         // durable
		false,        // auto-delete
		false,        // internal
		false,        // no-wait
		nil,          // arguments
	)
	if err != nil {
		logger.Fatalf("Failed to declare exchange %s: %v", exchangeName, err)
	}
	logger.Infof("Exchange declared: %s", exchangeName)

	for _, queueName := range queueNames {
		// Declare the queue
		err := declareQueue(channel, queueName)
		if err != nil {
			logger.Fatalf("Failed to declare queue %s: %v", queueName, err)
		}
		logger.Infof("Queue declared: %s", queueName)

		// Bind the queue to the exchange with a routing key
		err = channel.QueueBind(
			queueName,    // queue name
			queueName,    // routing key (assuming one-to-one mapping; change if needed)
			exchangeName, // exchange name
			false,        // no-wait
			nil,          // arguments
		)
		if err != nil {
			logger.Fatalf("Failed to bind queue %s to exchange %s: %v", queueName, exchangeName, err)
		}
		logger.Infof("Queue %s bound to exchange %s with routing key %s", queueName, exchangeName, queueName)

		// Start consuming messages from the queue
		msgs, err := channel.Consume(
			queueName, // queue
			"",        // consumer tag
			true,      // auto-ack
			false,     // exclusive
			false,     // no-local
			false,     // no-wait
			nil,       // args
		)
		if err != nil {
			logger.Fatalf("Failed to register a consumer for queue %s: %v", queueName, err)
		}
		logger.Infof("Consuming messages from queue: %s", queueName)

		// Handle messages
		go func(queueName string) {
			for msg := range msgs {
				controller.HandleMessage(msg, channel)
			}
		}(queueName)
	}
}

// declareQueue declares a queue in RabbitMQ.
func declareQueue(channel *amqp.Channel, queueName string) error {
	_, err := channel.QueueDeclare(
		queueName, // name
		false,     // durable
		false,     // delete when unused
		false,     // exclusive
		false,     // no-wait
		nil,       // arguments
	)
	return err
}
