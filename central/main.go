package main

import (
	"log"
	"context"
	pb "github.com/Sistemas-Distribuidos-2023-02/Grupo15-Laboratorio-1/tree/main/central/proto"
)

type server struct {
	pb.UnimplementedCentralServer
}

func (s *server) Keygen(ctx context.Context, message *pb.Message) (*pb.Message, error) {
	log.Printf("Received message body from client: %s", message.Body)
	return &pb.Message{Body: "0"}, nil
}