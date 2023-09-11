package main

import (
	"context"
	"fmt"
	"io/ioutil"
	"math/rand"
	"net"
	"strconv"
	// "strings"
	"time"
	// "os"
	// "bufio"

	"encoding/json"

	"github.com/Sistemas-Distribuidos-2023-02/Grupo15-Laboratorio-1/proto/betakeys"
	"github.com/streadway/amqp"
	"google.golang.org/grpc"
	// "google.golang.org/grpc/reflection"
	//"google.golang.org/protobuf/types/known/emptypb"
	//"google.golang.org/protobuf"
)

type server struct{}

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

	filePath := "./parametros_de_inicio.txt"
	parametroInicio, err := obtenerParametroInicio(filePath)
	if err != nil {
		fmt.Printf("Error reading startup_parameters: %v\n", err)
		return
	}

	cantidadUsuarios := CantidadUsuarios(parametroInicio)


	fmt.Printf("Cantidad de usuarios es %d\n", cantidadUsuarios)


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


	// // clientesss
	// conn, err := grpc.Dial("localhost:50051", grpc.WithInsecure())

	// if err != nil {
	// 	fmt.Printf("Error al conectar : %v\n", err)
	// }

	// defer conn.Close()

	// // crear cliente
	// client := betakeys.NewBetakeysServiceClient(conn)

	// // Crea una instancia del mensaje KeyNotification
    // notificacion := &betakeys.KeyNotification{
    //     KeygenNumber: "12345", // Establece el valor del campo
    // }
	
	// // Llama al método remoto NotifyRegionalServers de manera síncrona
    // respuesta, err := client.NotifyRegionalServers(context.Background(), notificacion)
    // if err != nil {
    //     fmt.Printf("Error al llamar al metodo: %v\n", err)
    // }

    // fmt.Printf("Respuesta del servidor: %v\n", respuesta)




	// // Create and set up gRPC server
	// grpcServer := grpc.NewServer()

	// betakeys.RegisterBetakeysServiceServer(grpcServer, &server{})

	// reflection.Register(grpcServer)

	// listener, err := net.Listen("tcp", ":50051")
	// if err != nil {
	// 	fmt.Printf("Failed to listen: %v\n", err)
	// 	return
	// }

	// // Start gRPC server
	// fmt.Println("Starting gRPC server on port: 50051")
	// go func() {
	// 	if err := server.Serve(listener); err != nil {
	// 		log.Fatalf("Error al servir: %v", err)
	// 	}
	// }()



	// Start gRPC server
	// fmt.Println("Starting gRPC server on port: 50051")
	// if err := grpcServer.Serve(listener); err != nil {
	// 	fmt.Printf("Failed to serve: %v\n", err)
	// 	return
	// }

}