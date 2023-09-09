package main

import (
	"context"
	"log"
	"net"

	pb "github.com/Sistemas-Distribuidos-2023-02/Grupo15-Laboratorio-1/tree/main/central/proto"
	"google.golang.org/grpc"
)

type server struct {
	pb.UnimplementedCentralServer
}

func (s *server) Keygen(ctx context.Context, message *pb.Message) (*pb.Message, error) {
	log.Printf("Received message body from client: %s", message.Body)
	return &pb.Message{Body: "0"}, nil
}

func main() {
	listener, err := net.Listen("tcp", ":8080")

	if err != nil {
		panic("cannot create tcp server" + err.Error())
	}

	serv := grpc.NewServer()

	pb.RegisterCentralServer(serv, &server{})

	if err = serv.Serve(listener); err != nil {
		panic("cannot initialize grpc server" + err.Error())
	}
}
