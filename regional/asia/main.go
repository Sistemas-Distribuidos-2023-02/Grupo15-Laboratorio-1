package main

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"math/rand"
	"net"
	"strconv"
	"time"
	"log"

	"github.com/Sistemas-Distribuidos-2023-02/Grupo15-Laboratorio-1/proto/betakeys"
	"github.com/streadway/amqp"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
	"google.golang.org/protobuf/types/known/emptypb"
)


type server struct {
	keys int
	serverName string
	betakeys.UnimplementedBetakeysServiceServer
}


func (s *server) NotifyRegionalServers (ctx context.Context, request *betakeys.KeyNotification) (*emptypb.Empty, error) {

	keygenNumber := request.KeygenNumber
	ServerName := s.serverName
	log.Println("Notificacion recibida:", keygenNumber, "llaves generadas por central")

	if err := enviarUsuariosAQueue(int(keygenNumber), ServerName); err != nil {
		log.Fatalf("Error al comunicar con cola Rabbit: %v", err)
	}

	return &emptypb.Empty{}, nil
}

func (s *server) SendResponseToRegionalServer (ctx context.Context, request *betakeys.ResponseToRegionalServer) (*emptypb.Empty, error) {
	
	accepted := request.Accepted
	denied := request.Denied
	targetServerName := request.TargetServerName
	log.Println("Se inscribieron cupos en el servidor" ,targetServerName, ":" ,accepted, "inscritos," ,denied, "denegados")
	
	if denied == 0 {
		return &emptypb.Empty{}, nil
	} else {
		keygenNumber := denied
		if err := enviarUsuariosAQueue(int(keygenNumber), targetServerName); err != nil {
			log.Fatalf("Error al comunicar con cola Rabbit: %v", err)
		}
	}

	return &emptypb.Empty{}, nil
}

func obtenerParametroInicio(nombreArchivo string) (parametro int, err error) {
    contenido, err := ioutil.ReadFile(nombreArchivo)
	if err != nil {
		return 0, err
	}
	parametroInicio, err := strconv.Atoi(string(contenido))
	if err != nil {
		return 0, err
	}
    return parametroInicio, nil
}

func generarValorAleatorio(min, max int) int {
	rand.Seed(time.Now().UnixNano())
	return rand.Intn(max-min+1) + min
}

func CantidadUsuarios(parametrosInicio int) int {
	aux := float64(parametrosInicio)/2
	min := int(aux - 0.2*aux)
	max := int(aux + 0.2*aux)
	return generarValorAleatorio(min, max)
}

type MensajeRegistro struct {
	NombreServidor string `json:"nombre_servidor"`
	Usuarios       int    `json:"usuarios"`
}

func enviarUsuariosAQueue(cantidad int, servidor string) error {
	// Conectar a RabbitMQ
	conn, err := amqp.Dial("amqp://usuario:contraseña@localhost:5672/")
	if err != nil {
		return err
	}
	defer conn.Close()

	// Abrir un canal
	ch, err := conn.Channel()
	if err != nil {
		return err
	}
	defer ch.Close()

	// Declarar una cola
	q, err := ch.QueueDeclare(
		"usuarios_registrados", // Nombre de la cola
		false,                  // Durable
		false,                  // Auto-borrado
		false,                  // Exclusivo
		false,                  // No esperar a la confirmación
		nil,                    // Argumentos adicionales
	)
	if err != nil {
		return err
	}

	// Crear un mensaje en formato JSON
	mensaje := MensajeRegistro{
		NombreServidor: servidor,
		Usuarios:       cantidad,
	}

	// Serializar el mensaje en JSON
	mensajeJSON, err := json.Marshal(mensaje)
	if err != nil {
		return err
	}

	// Publicar el mensaje en la cola
	err = ch.Publish(
		"",     // Intercambio (exchange) por defecto
		q.Name, // Cola a la que se envía el mensaje
		false,  // No esperar a la confirmación
		false,  // No requerir confirmación para la entrega
		amqp.Publishing{
			ContentType: "application/json",
			Body:        mensajeJSON,
		})
	if err != nil {
		return err
	}

	return nil
}

func main() {

	filePath := "regional/asia/parametros_de_inicio.txt"
	parametroInicio, err := obtenerParametroInicio(filePath)
	if err != nil {
		log.Fatalf("Error al leer archivo parametros: %v", err)
	}

	// Generate keys 
	cantidadUsuarios := CantidadUsuarios(parametroInicio)
	log.Println("Cantidad de usuarios es" ,cantidadUsuarios)

	// Receive notification from central server	aa
	grpcServer := grpc.NewServer()

	keyServer := &server{
		keys: cantidadUsuarios,
		serverName: "asia",
	}

	betakeys.RegisterBetakeysServiceServer(grpcServer, keyServer)

	reflection.Register(grpcServer)

	listener, err := net.Listen("tcp", ":50051")
	if err != nil {
		log.Fatalf("Error al escuchar tcp: %v", err)
	}

	// Start gRPC server
	fmt.Println("Starting gRPC server on port: 50051")
	if err := grpcServer.Serve(listener); err != nil {
		log.Fatalf("Error al Serve: %v", err)
	}

}