package main

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"math/rand"
	"net"
	"strconv"
	"strings"
	"time"
	"os"

	"github.com/Sistemas-Distribuidos-2023-02/Grupo15-Laboratorio-1/proto/betakeys"
	amqp "github.com/rabbitmq/amqp091-go"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
	"google.golang.org/protobuf/types/known/emptypb"
)

type server struct {
	betakeys.UnimplementedBetakeysServiceServer
}

type regionalMessage struct {
	ServerName string `json:"NombreServidor"`
	Content string `json:"Usuarios"`
}

func (s *server) NotifyRegionalServers(ctx context.Context, request *betakeys.KeyNotification) (*emptypb.Empty, error) {
	keygenNumber := request.KeygenNumber
	fmt.Printf("Received notification: %v keys generated\n", keygenNumber)
	return &emptypb.Empty{}, nil
}

func (s *server)SendResponseToRegionalServer(ctx context.Context, request *betakeys.ResponseToRegionalServer) (*emptypb.Empty, error) {
	accepted := request.Accepted
	denied := request.Denied
	targetServerName := request.TargetServerName
	fmt.Printf("Se inscribieron cupos en el servidor %v: %v inscritos, %v denegados\n", targetServerName, accepted, denied)
	return &emptypb.Empty{}, nil
}

func startupParameters(filePath string) (minKey, maxKey, ite int, err error) {
	content, err := ioutil.ReadFile(filePath)
	if err != nil {
		return 0, 0, 0, err
	}

	lines := strings.Split(string(content), "\n")

	firstLine := strings.TrimSpace(lines[0])

	cotas := strings.Split(firstLine, "-")

	minKey, err = strconv.Atoi(cotas[0])
	if err != nil {
		return 0, 0, 0, err
	}

	maxKey, err = strconv.Atoi(cotas[1])
	if err != nil {
		return 0, 0, 0, err
	}

	ite, err = strconv.Atoi(strings.TrimSpace(lines[1]))
	if err != nil {
		return 0, 0, 0, err
	}

	return minKey, maxKey, ite, nil
}

func keygen(minKey, maxKey int) int {
	rand.Seed(time.Now().UnixNano())
	return rand.Intn(maxKey - minKey + 1) + minKey
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

func rabbitMQMessageHandler(rabbitChannel *amqp.Channel, queueName string, numKeys *int, logFile *os.File){
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
		
		log.Printf("Se han registrado %v usuarios, %v solicitados, %v denegados.\n", numRegistered, numUsers, numIgnored)

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

func main() {
	// Create and set up gRPC server
	grpcServer := grpc.NewServer()

	betakeys.RegisterBetakeysServiceServer(grpcServer, &server{})

	reflection.Register(grpcServer)

	listener, err := net.Listen("tcp", ":50051")
	if err != nil {
		fmt.Printf("Failed to listen: %v\n", err)
		return
	}

	// Start gRPC server
	fmt.Println("Starting gRPC server on port: 50051")
	go func() {
		if err := grpcServer.Serve(listener); err != nil {
			fmt.Printf("Failed to serve: %v\n", err)
			return
		}
	}()

	// Generate keys and read start up parameters
	filePath := "./parametros_de_inicio.txt"

	minKey, maxKey, ite, err := startupParameters(filePath) // the "ite" variable is not used yet in this version of the code, it will be used in later implementation
	if err != nil {
		fmt.Printf("Error reading startup_parameters: %v\n", err)
		return
	}

	// Connect to RabbitMQ server
	const rabbitmqURL = "amqp://guest:guest@localhost:5672/"
	rabbitConn, err := amqp.Dial(rabbitmqURL)
    if err != nil {
        log.Fatalf("Failed to connect to RabbitMQ server: %v", err)
    }
    defer rabbitConn.Close()

    rabbitChannel, err := rabbitConn.Channel()
    if err != nil {
        log.Fatalf("Failed to open a RabbitMQ channel: %v", err)
    }
    defer rabbitChannel.Close()


	// RabbitMQ queue declaration
	queueName := "keyVolunteers"

	_, err = rabbitChannel.QueueDeclare(
        queueName,
        false, // durable
        false, // delete when unused
        false, // exclusive
        false, // no-wait
        nil,   // arguments
    )
    if err != nil {
        log.Fatalf("Failed to declare RabbitMQ queue: %v", err)
    }

	// Prepare logging file
	logFile, err := os.OpenFile("central.log", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		fmt.Printf("Failed to open log file: %v\n", err)
		return
	}
	log.SetOutput(logFile)

	// Begin iterations
	var count int = 0
	for ite != 0 {
		count++
		if ite == -1{
			fmt.Printf("Generacion %v/%v\n", count, "infinito")
		}
		if ite > 0 {
			fmt.Printf("Generacion %v/%v\n", count, ite)
		}

		keys := keygen(minKey, maxKey)
		
		log.Printf("%v llaves generadas\n", keys)

		// Send notification to regional servers
		conn, err := grpc.Dial("localhost:50051", grpc.WithInsecure())
		if err != nil {
			fmt.Printf("Failed to connect to gRPC server: %v\n", err)
			return
		}
		defer conn.Close()

		client := betakeys.NewBetakeysServiceClient(conn)
		notification := &betakeys.KeyNotification{
			KeygenNumber: int32(keys),
		}

		_, err = client.NotifyRegionalServers(context.Background(), notification)
		if err != nil {
			fmt.Printf("Failed to send notification: %v\n", err)
			return
		}

		// Message handling
		rabbitMQMessageHandler(rabbitChannel, queueName, &keys, logFile)

		if ite > 0 {
			ite--
		}
	}
	
}