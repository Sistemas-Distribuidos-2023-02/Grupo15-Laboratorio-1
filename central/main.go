package main

import (
	"context"
	"fmt"
	"io/ioutil"
	"math/rand"
	"net"
	"strconv"
	"strings"
	"time"
	"os"
	"os/signal"
	"syscall"
	"sync"

	"github.com/Sistemas-Distribuidos-2023-02/Grupo15-Laboratorio-1/proto/betakeys"
	"github.com/streadway/amqp"
	"google.golang.org/protobuf/proto"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
	"google.golang.org/protobuf/types/known/emptypb"
)

type server struct {}

func (s *server) NotifyRegionalServers(ctx context.Context, request *betakeys.KeyNotification) (*emptypb.Empty, error) {
	keygenNumber := request.KeygenNumber
	fmt.Printf("Received notification: %v keys generated\n", keygenNumber)
	return &emptypb.Empty{}, nil
}

func startupParameters(filePath string) (minKey, maxKey int, err error) {
	content, err := ioutil.ReadFile(filePath)
	if err != nil {
		return 0, 0, err
	}

	cotas := strings.Split(string(content), "-")
	if len(cotas) != 2 {
		return 0, 0, fmt.Errorf("invalid file format in parametros_de_inicio.txt")
	}

	minKey, err = strconv.Atoi(strings.TrimSpace(cotas[0]))
	if err != nil {
		return 0, 0, err
	}

	maxKey, err = strconv.Atoi(strings.TrimSpace(cotas[1]))
	if err != nil {
		return 0, 0, err
	}

	return minKey, maxKey, nil
}

func keygen(minKey, maxKey int) []int {
	rand.Seed(time.Now().UnixNano())

	n := rand.Intn(maxKey - minKey + 1) + minKey

	keys := make([]int, n)
	for i := 0; i < n; i++ {
		keys[i] = minKey + rand.Intn(n)
	}
	return keys
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
	if err := grpcServer.Serve(listener); err != nil {
		fmt.Printf("Failed to serve: %v\n", err)
		return
	}

	// Generate keys
	filePath := "./parametros_de_inicio.txt"

	minKey, maxKey, err := startupParameters(filePath)
	if err != nil {
		fmt.Printf("Error reading startup_parameters: %v\n", err)
		return
	}

	keys := keygen(minKey, maxKey)
	fmt.Printf("Keys: %v\n", keys)

	// Send notification to regional servers
	conn, err := grpc.Dial("localhost:50051", grpc.WithInsecure())
	if err != nil {
		fmt.Printf("Failed to connect to gRPC server: %v\n", err)
		return
	}
	defer conn.Close()

	client := betakeys.NewBetakeysServiceClient(conn)
	notification := &betakeys.KeyNotification{
		KeygenNumber: string(rune(len(keys))),
	}

	_, err = client.NotifyRegionalServers(context.Background(), notification)
	if err != nil {
		fmt.Printf("Failed to send notification: %v\n", err)
		return
	}

	//RabbitMQ server connection
	const rabbitmqURL = "amqp://guest:guest@localhost:25672/"

	rabbitConn, err := amqp.Dial(rabbitmqURL)
	if err != nil {
		fmt.Printf("Failed to connect to RabbitMQ server: %v\n", err)
		return
	}
	defer rabbitConn.Close()

	rabbitChannel, err := rabbitConn.Channel()
	if err != nil {
		fmt.Printf("Failed to open a RabbitMQ channel: %v\n", err)
		return
	}
	defer rabbitChannel.Close()

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
		fmt.Printf("Failed to declare RabbitMQ queue: %v\n", err)
		return
	}

	RabbitShutdownChannel := make(chan os.Signal, 1)
	signal.Notify(RabbitShutdownChannel, syscall.SIGINT, syscall.SIGTERM)
	var wg sync.WaitGroup

	// Asynchronous message reception
	wg.Add(1)
	go func ()  {
		// Consuming message
		msgs, err := rabbitChannel.Consume(
			queueName,
			"",		// consumer
			true,	// auto-ack
			false,	// exclusive
			false,	// no-local
			false,	// no-wait
			nil,	// arguments
		)
		if err != nil {
			fmt.Printf("Failed to register a regional server user: %v\n", err)
			return
		}

		for {
			select{
				case <- RabbitShutdownChannel:
					fmt.Printf("RabbitMQ server shutting down gracefully...\n")
					wg.Done()
					return
				case msg := <- msgs:
					// Unmarshalling message
					var reply betakeys.RegionalKeyReply

					err := proto.Unmarshal(msg.Body, &reply)
					if err != nil {
						fmt.Printf("Failed to unmarshal message: %v\n", err)
						return
					}
			
					serverName := reply.ServerName
					content := reply.Content
			
					fmt.Printf("Mensaje asincrono de servidor %v leido: %v\n", serverName, content)
					msg.Ack(false)
					wg.Done()
			}

		}
		
	} ()
	
	// Waiting for shutdown signal
	<- RabbitShutdownChannel

	// Graceful shutdown
	wg.Wait()
	rabbitChannel.Close()
	
}
