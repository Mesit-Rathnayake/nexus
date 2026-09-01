package network

import (
	"fmt"
	"log"
	"net"
)

type Server struct {
	Address string
}

func NewServer(port int) *Server {
	return &Server{
		Address: fmt.Sprintf(":%d", port),
	}
}

func (s *Server) Start() error {
	listener, err := net.Listen("tcp", s.Address)
	if err != nil {
		return fmt.Errorf("failed to listen on %s: %w", s.Address, err)
	}

	log.Printf("Nexus P2P server listening on %s", s.Address)

	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Printf("failed to accept connection: %v", err)
			continue
		}

		go s.handleConnection(conn)
	}
}

func (s *Server) handleConnection(conn net.Conn) {
	defer conn.Close()

	log.Printf("Peer connected: %s", conn.RemoteAddr())

	message, err := Receive(conn)
	if err != nil {
		log.Printf("failed to receive message: %v", err)
		return
	}

	log.Printf(
		"Received message: %s",
		message.Type,
	)

	if message.Type == MessageTypePing {
		response, err := NewMessage(
			MessageTypePong,
			map[string]string{
				"message": "hello from Nexus",
			},
		)

		if err != nil {
			log.Printf("failed to create response: %v", err)
			return
		}

		response = append(response, '\n')

		if _, err := conn.Write(response); err != nil {
			log.Printf("failed to send response: %v", err)
		}
	}
}
