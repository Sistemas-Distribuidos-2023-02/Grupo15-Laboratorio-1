package main

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"math/rand"

	//"net"
	"os"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/Sistemas-Distribuidos-2023-02/Grupo15-Laboratorio-1/proto/betakeys"
	amqp "github.com/rabbitmq/amqp091-go"
	"google.golang.org/grpc"
	"google.golang.org/protobuf/types/known/emptypb"
)

type regionalMessage struct {
	NombreServidor string `json:"nombre_servidor"`
	Usuarios       int    `json:"usuarios"`
}

func NotifyRegionalServers(ctx context.Context, request *betakeys.KeyNotification) (*emptypb.Empty, error) {
	keygenNumber := request.KeygenNumber
	log.Println("Notificacion recibida:", keygenNumber, "llaves generadas")
	return &emptypb.Empty{}, nil
}

func SendResponseToRegionalServer(ctx context.Context, request *betakeys.ResponseToRegionalServer) (*emptypb.Empty, error) {
	accepted := request.Accepted
	denied := request.Denied
	targetServerName := request.TargetServerName
	log.Println("Se inscribieron cupos en el servidor", targetServerName, ":", accepted, "inscritos,", denied, "denegados")
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
	return rand.Intn(maxKey-minKey+1) + minKey
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
	host := ""
	switch serverName {
	case "america":
		host = "dist057:50051"
	case "asia":
		host = "dist058:50052"
	case "europa":
		host = "dist059:50053"
	case "oceania":
		host = "dist060:50054"
	}

	conn, err := grpc.Dial(host, grpc.WithInsecure())
	if err != nil {
		log.Fatalf("Fallo en conectar gRPC server: %v", err)
	}

	client := betakeys.NewBetakeysServiceClient(conn)
	response := &betakeys.ResponseToRegionalServer{
		TargetServerName: serverName,
		Accepted:         numRegistered,
		Denied:           numIgnored,
	}

	_, err = client.SendResponseToRegionalServer(context.Background(), response)
	if err != nil {
		return fmt.Errorf("failed to send response to regional server while sending results: %v", err)
	}

	return nil
}

func rabbitMQMessageHandler(rabbitChannel *amqp.Channel, queueName string, numKeys *int, logFile *os.File) {
	var wg sync.WaitGroup
	// Consumer
	msgs, err := rabbitChannel.Consume(
		queueName,
		"",    // consumer
		false, // auto-ack
		false, // exclusive
		false, // no-local
		false, // no-wait
		nil,   // arguments
	)
	if err != nil {
		fmt.Printf("failed to register a RabbitMQ consumer: %v", err)
		return
	}

	messageChannel := make(chan regionalMessage, 4)

	// Message handling loop
	for msg := range msgs {
		// Check if queue is done
		if len(msgs) == 0 {
			fmt.Printf("No hay mas mensajes en la cola.\n")
			break
		}

		var message regionalMessage
		if err := json.Unmarshal(msg.Body, &message); err != nil {
			fmt.Printf("failed to unmarshal message: %v\n", err)
			msg.Ack(false)
			continue
		}

		regionalServerName := strings.TrimSpace(message.NombreServidor)
		numUsers := message.Usuarios
		if err != nil {
			fmt.Printf("invalid message format from regional servers: %v\n", err)
			return
		}

		fmt.Printf("Mensaje asincrono de servidor %v recibido.\n\n", regionalServerName)

		// Process message
		numRegistered, numIgnored := messageProcessing(numUsers, numKeys)

		fmt.Printf("Se han registrado %v usuarios, %v solicitados, %v denegados.\n", numRegistered, numUsers, numIgnored)

		// Send results to regional server
		wg.Add(1)
		go func(regionalServerName string, numRegistered, numIgnored int32) {
			defer wg.Done()

			err = sendResultsToRegionalServer(regionalServerName, int32(numRegistered), int32(numIgnored))
			if err != nil {
				fmt.Printf("failed to send results to regional server %v: %v\n", regionalServerName, err)
				return
			} else {
				log.Printf("Resultados enviados a servidor %v.\nUsuarios registrados: %v\nUsuarios rechazados: %v\n\n", regionalServerName, numRegistered, numIgnored)
			}
		}(regionalServerName, int32(numRegistered), int32(numIgnored))

		// Acknowledge message
		msg.Ack(false)
	}
	wg.Wait()
	close(messageChannel)
}

func main() {

	// 50051: AMERICA, 50052: ASIA, 50053: EUROPA, 50054: OCEANIA,
	serverAddresses := []string{"dist057:50051", "dist058:50052", "dist059:50053", "dist060:50054"}
	clients := []betakeys.BetakeysServiceClient{}

	for _, serverAddress := range serverAddresses {
		conn, err := grpc.Dial(serverAddress, grpc.WithInsecure())
		if err != nil {
			log.Fatalf("Fallo en conectar gRPC server en %s: %v", serverAddress, err)
		}
		defer conn.Close()

		// Crea un cliente para el servicio gRPC en cada servidor
		client := betakeys.NewBetakeysServiceClient(conn)
		clients = append(clients, client)

	}
	// Read start up parameters
	filePath := "./parametros_de_inicio.txt"

	minKey, maxKey, ite, err := startupParameters(filePath) // the "ite" variable is not used yet in this version of the code, it will be used in later implementation
	if err != nil {
		log.Fatalf("Error al leer archivo parametros: %v", err)
	}

	// Connect to RabbitMQ server
	const rabbitmqURL = "amqp://guest:guest@dist059:5673/"
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
	for count != ite {

		if ite == -1 {
			fmt.Println("Generacion", count+1, "/infinito")
		}
		if count < ite-1 {
			fmt.Println("Generacion", count+1, "/", ite)
		}

		keys := keygen(minKey, maxKey)
		log.Printf("%v llaves generadas\n", keys)

		// Send notification to regional servers
		goResults := make(chan error, 4)

		notification := &betakeys.KeyNotification{
			KeygenNumber: int32(keys),
		}

		var wg sync.WaitGroup

		sendKeygenNotification := func(client betakeys.BetakeysServiceClient) {
			_, err := client.NotifyRegionalServers(context.Background(), notification)
			if err != nil {
				goResults <- err
			} else {
				goResults <- nil
			}
			defer wg.Done()
		}

		wg.Add(4)
		for _, client := range clients {
			go sendKeygenNotification(client)
		}
		wg.Wait()

		close(goResults)

		for result := range goResults {
			if result != nil {
				fmt.Printf("failed to send keygen notification to regional server: %v\n", result)

			}
		}

		// Start RabbitMQ message handler
		// The rabbitMQMessageHandler function is the one that will go through the messages from the regional servers waiting in the RabbitMQ queue
		// The message handler will then process the messages and send the results to the regional servers
		rabbitMQMessageHandler(rabbitChannel, "keyVolunteers", &keys, logFile)
		fmt.Println("Message handler done.")

		count++
		time.Sleep(time.Second * 3)
	}
}
