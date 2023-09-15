package main

import (
	"context"
	"fmt"
	"io/ioutil"
	"math/rand"
	"strconv"
	"strings"
	"time"
	"log"
	"encoding/json"
	
	amqp "github.com/rabbitmq/amqp091-go"
	"github.com/Sistemas-Distribuidos-2023-02/Grupo15-Laboratorio-1/proto/betakeys"
	"google.golang.org/grpc"
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
	log.Println("Notificacion recibida:", keygenNumber, "llaves generadas")
	return &emptypb.Empty{}, nil
}

func (s *server)SendResponseToRegionalServer(ctx context.Context, request *betakeys.ResponseToRegionalServer) (*emptypb.Empty, error) {
	accepted := request.Accepted
	denied := request.Denied
	targetServerName := request.TargetServerName
	log.Println("Se inscribieron cupos en el servidor" ,targetServerName, ":" ,accepted, "inscritos," ,denied, "denegados")
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

func SetupRabbitMQ()(*amqp.Channel, error){
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

func RabbitMQMessageHandler(rabbitChannel *amqp.Channel, queueName string, numKeys *int){
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

func main() {
	
	// Generate keys and read start up parameters
	filePath := "central/parametros_de_inicio.txt"

	minKey, maxKey, ite, err := startupParameters(filePath) // the "ite" variable is not used yet in this version of the code, it will be used in later implementation
	if err != nil {
		log.Fatalf("Error al leer archivo parametros: %v", err)
	}

	keys := keygen(minKey, maxKey)
	log.Println("Se generaron" ,keys, "llaves")
	// Begin iterations
	var count int = 0
	for count != ite {
		if ite == -1{
			log.Println("Generacion",count+1,"/infinito")
		} else {
			log.Println("Generacion" ,count+1, "/" ,ite)
		}

		// Send notification to regional servers
		conn, err := grpc.Dial("localhost:50051", grpc.WithInsecure())
		if err != nil {
			log.Fatalf("Fallo en conectar gRPC server: %v", err)
		}

		// NotifyRegionalServers, enviar llaves disp
		client := betakeys.NewBetakeysServiceClient(conn)
		notification := &betakeys.KeyNotification{
			KeygenNumber: int32(keys),
		}
		_, err = client.NotifyRegionalServers(context.Background(), notification)
		if err != nil {
			log.Fatalf("Fallo en enviar notificacion: %v", err)
		}

		// Set up RabbitMQ
		rabbitChannel, err := SetupRabbitMQ() // queueName = "keyVolunteers"
		if err != nil {
			fmt.Printf("Failed to set up RabbitMQ: %v", err)
			return
		}
		defer rabbitChannel.Close()

		// Start RabbitMQ message handler
		// The rabbitMQMessageHandler function is the one that will go through the messages from the regional servers waiting in the RabbitMQ queue
		// The message handler will then process the messages and send the results to the regional servers
		RabbitMQMessageHandler(rabbitChannel, "keyVolunteers", &keys)

		count++
	}
	
}
