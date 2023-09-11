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
	
	mq "github.com/Sistemas-Distribuidos-2023-02/Grupo15-Laboratorio-1/rabbitmq"
	"github.com/Sistemas-Distribuidos-2023-02/Grupo15-Laboratorio-1/proto/betakeys"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
	"google.golang.org/protobuf/types/known/emptypb"
)

type server struct {
	betakeys.UnimplementedBetakeysServiceServer
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

		// Send notification to regional servers
		conn, err := grpc.Dial("localhost:50051", grpc.WithInsecure())
		if err != nil {
			fmt.Printf("Failed to connect to gRPC server: %v\n", err)
			return
		}
		defer conn.Close()

		client := betakeys.NewBetakeysServiceClient(conn)
		notification := &betakeys.KeyNotification{
			KeygenNumber: strconv.Itoa(keys),
		}

		_, err = client.NotifyRegionalServers(context.Background(), notification)
		if err != nil {
			fmt.Printf("Failed to send notification: %v\n", err)
			return
		}

		// Set up RabbitMQ
		rabbitChannel, err := mq.SetupRabbitMQ() // queueName = "keyVolunteers"
		if err != nil {
			fmt.Printf("Failed to set up RabbitMQ: %v\n", err)
			return
		}
		defer rabbitChannel.Close()

		// Start RabbitMQ message handler
		// The rabbitMQMessageHandler function is the one that will go through the messages from the regional servers waiting in the RabbitMQ queue
		// The message handler will then process the messages and send the results to the regional servers
		go mq.RabbitMQMessageHandler(rabbitChannel, "keyVolunteers", &keys)

		if ite > 0 {
			ite--
		}
	}
	
}
