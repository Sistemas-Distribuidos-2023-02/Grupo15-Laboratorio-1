package main

import (
	"log"
	"net"

	"github.com/Sistemas-Distribuidos-2023-02/Grupo15-Laboratorio-1/tree/main/central/proto"
	"google.golang.org/grpc"
)

func main() {
	lis, err := net.Listen("tcp", ":9000")
	if err != nil {
		log.Fatal("Failed to listen on port 9000: %v", err)
	}

	s := central.Server{}

	grpcServer := grpc.NewServer()

	central.RegisterCentralServer(grpcServer, &s)

	if err := grpcServer.Serve(lis); err != nil {
		log.Fatal("Failed to serve gRPC server over port 9000: %v", err)
	}
}
