package main

import (
	"log"

	"github.com/streadway/amqp"
)

func main() {
	// Establish a connection to RabbitMQ
	conn, err := amqp.Dial("amqp://localhost:5673")
	if err != nil {
		log.Fatalf("Failed to connect to RabbitMQ: %v", err)
	}
	defer conn.Close()

	// Create a channel
	channel, err := conn.Channel()
	if err != nil {
		log.Fatalf("Failed to open a channel: %v", err)
	}
	defer channel.Close()

	// Declare a queue
	queueName := "my-queue"
	_, err = channel.QueueDeclare(
		queueName, // Name of the queue
		true,      // Durable
		false,     // Delete when unused
		false,     // Exclusive
		false,     // No-wait
		nil,       // Arguments
	)
	if err != nil {
		log.Fatalf("Failed to declare a queue: %v", err)
	}

	log.Printf("Queue %s declared successfully", queueName)

	// You can now use the 'channel' for other operations or continue with your application logic.
}
