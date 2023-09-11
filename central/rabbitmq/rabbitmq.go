package rabbitmq

import (
	"fmt"
	"strconv"
	"strings"
	"context"
	"encoding/json"

	"google.golang.org/protobuf/types/known/emptypb"
	"google.golang.org/grpc"
	amqp "github.com/rabbitmq/amqp091-go"
	"github.com/Sistemas-Distribuidos-2023-02/Grupo15-Laboratorio-1/proto/betakeys"
)

type server struct {
	betakeys.UnimplementedBetakeysServiceServer
}

type regionalMessage struct {
	ServerName string `json:"serverName"`
	Content string `json:"content"`
}

func (s *server)SendResponseToRegionalServer(ctx context.Context, request *betakeys.ResponseToRegionalServer) (*emptypb.Empty, error) {
	accepted := request.Accepted
	denied := request.Denied
	targetServerName := request.TargetServerName
	fmt.Printf("Se inscribieron cupos en el servidor %v: %v inscritos, %v denegados\n", targetServerName, accepted, denied)
	return &emptypb.Empty{}, nil
}

func setupRabbitMQ()(*amqp.Channel, error){
	//RabbitMQ server connection
	const rabbitmqURL = "amqp://guest:guest@localhost:25672/"

	rabbitConn, err := amqp.Dial(rabbitmqURL)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to rabbitmq server: %v", err)
	}

	rabbitChannel, err := rabbitConn.Channel()
	if err != nil {
		return nil, fmt.Errorf("failed to open a rabbitmq channel: %v", err)
	}

	// RabbitMQ queue declaration
	queueName := "keyVolunteers"

	_, err = rabbitChannel.QueueDeclare(
		queueName,
		false,	// durable
		false,	// delete when unused
		false,	// exclusive
		false,	// no-wait
		nil,	// arguments
	)
	if err != nil {
		return nil, fmt.Errorf("failed to declare rabbitmq queue: %v", err)
	}

	return rabbitChannel, nil
}

func messageProcessing(numUsers int, numKeys *int) (numRegistered, numIgnored int) {
	if numUsers > *numKeys {
		numIgnored = numUsers - *numKeys
		numUsers = *numKeys
	}

	*numKeys -= numUsers
	if *numKeys < 0 {
		*numKeys = 0
	}
	numRegistered = numUsers

	return numRegistered, numIgnored
}

func sendResultsToRegionalServer(serverName string, numRegistered, numIgnored int32) error {
	// Connect to gRPC server
	conn, err := grpc.Dial("localhost:50051", grpc.WithInsecure())
	if err != nil {
		return fmt.Errorf("failed to connect to gRPC server while sending results: %v", err)
	}
	defer conn.Close()

	client := betakeys.NewBetakeysServiceClient(conn)
	response := &betakeys.ResponseToRegionalServer{
		TargetServerName: serverName,
		Accepted: numRegistered,
		Denied: numIgnored,
	}

	_, err = client.SendResponseToRegionalServer(context.Background(), response)
	if err != nil {
		return fmt.Errorf("failed to send response to regional server while sending results: %v", err)
	}

	return nil
}

func rabbitMQMessageHandler(rabbitChannel *amqp.Channel, queueName string, numKeys *int){
	// Consumer
	msgs, err := rabbitChannel.Consume(
		queueName,
		"",		// consumer
		false,	// auto-ack
		false,	// exclusive
		false,	// no-local
		false,	// no-wait
		nil,	// arguments
	)
	if err != nil {
		fmt.Printf("failed to register a RabbitMQ consumer: %v", err)
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

		fmt.Printf("Mensaje asincrono de servidor %v recibido.", regionalServerName)

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