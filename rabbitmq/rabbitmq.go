package main

import (
	"fmt"
	"strconv"
	"strings"
	"context"
	"encoding/json"
	"log"

	"google.golang.org/protobuf/types/known/emptypb"
	"google.golang.org/grpc"

	amqp "github.com/rabbitmq/amqp091-go"
)

func main() {
	conn, err := amqp.Dial("amqp://guest:guest@localhost:5672/")
	if err != nil {
		return nil, fmt.Errorf("failed to declare rabbitmq queue: %v", err)
	}

	return rabbitChannel, nil
}

func messageProcessing(numUsers int, numKeys *int) (numRegistered, numIgnored int) {
	if numUsers > *numKeys {
		numRegistered = *numKeys
		numIgnored = numUsers - *numKeys
	} else {
		numRegistered = numUsers
		numIgnored = 0
	}

	return numRegistered, numIgnored
}

func sendResultsToRegionalServer(serverName string, numRegistered, numIgnored int32) error {
	// Connect to gRPC server
	conn, err := grpc.Dial("localhost:50051", grpc.WithInsecure())
	if err != nil {
		log.Fatalf("Fallo en conectar gRPC server: %v", err)

	}

	ch, err := conn.Channel()
	if err != nil {
		log.Fatalf("Failed to open a RabbitMQ channel: %v", err)
		return
	}
	defer ch.Close()

	q, err := ch.QueueDeclare(
		"hello", // name
		false,   // durable
		false,   // delete when unused
		true,    // exclusive
		false,   // no-wait
		nil,     // arguments
	)
	if err != nil {
		log.Fatalf("Failed to declare a RabbitMQ queue %v: %v", q, err)
		return
	}

	// Message handling loop
	for msg := range msgs {
		var message regionalMessage
		if err := json.Unmarshal(msg.Body, &message); err != nil {
			fmt.Printf("failed to unmarshal message: %v\n", err)
			msg.Ack(false)
			continue
		}

		regionalServerName := strings.TrimSpace(message.ServerName)
		numUsers, err := strconv.Atoi(strings.TrimSpace(message.Content))
		if err != nil {
			fmt.Printf("invalid message format from regional servers: %v\n", err)
			return
		}

		fmt.Printf("Mensaje asincrono de servidor %v recibido. Numero de usuarios recibido: %v", regionalServerName, numUsers)

		// Process message
		numRegistered, numIgnored := messageProcessing(numUsers, numKeys)

		// Send results to regional server
		err = sendResultsToRegionalServer(regionalServerName, int32(numRegistered), int32(numIgnored))
		if err != nil {
			fmt.Printf("failed to send results to regional server %v: %v\n", regionalServerName, err)
			return
		}

		// Acknowledge message
		msg.Ack(false)

	}

}