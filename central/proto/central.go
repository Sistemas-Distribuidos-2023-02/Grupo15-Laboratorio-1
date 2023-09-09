package central

import (
	"log"
	"golang.org/x/net/context"
)

type Server struct{}

func (s *Server) Keygen(ctx context.Context, message *Message) (*Message, error) {
	log.Printf("Received message body from client: %s", message.Body)
	return &Message{keyAmount: "0"}, nil
}