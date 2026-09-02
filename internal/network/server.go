package network

import (
	"encoding/json"
	"fmt"
	"log"
	"net"
)

type Server struct {
	Address     string
	NodeID      string
	NodeAddress string
	Peers       *PeerManager
}

func NewServer(
	port int,
	nodeID string,
	nodeAddress string,
	peers *PeerManager,
) *Server {
	return &Server{
		Address:     fmt.Sprintf(":%d", port),
		NodeID:      nodeID,
		NodeAddress: nodeAddress,
		Peers:       peers,
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

	switch message.Type {

	case MessageTypeHello:
		s.handleHello(conn, message)

	case MessageTypePing:
		s.handlePing(conn)

	default:
		log.Printf(
			"Unsupported message type: %s",
			message.Type,
		)
	}
}

func (s *Server) handlePing(conn net.Conn) {
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

func (s *Server) handleHello(
	conn net.Conn,
	message *Message,
) {
	var hello HelloPayload

	if err := json.Unmarshal(message.Payload, &hello); err != nil {
		log.Printf("failed to decode HELLO payload: %v", err)
		return
	}

	if hello.NodeID == "" {
		log.Printf("HELLO contains empty node ID")
		return
	}

	if hello.Address == "" {
		log.Printf("HELLO contains empty node address")
		return
	}

	log.Printf(
		"Handshake from node %s at %s",
		hello.NodeID,
		hello.Address,
	)

	peer := &Peer{
		Address: hello.Address,
		Conn:    conn,
	}

	if err := s.Peers.Add(hello.NodeID, peer); err != nil {
		log.Printf(
			"failed to register peer %s: %v",
			hello.NodeID,
			err,
		)
		return
	}

	response := HelloPayload{
		NodeID:  s.NodeID,
		Address: s.NodeAddress,
	}

	if err := sendMessage(
		conn,
		MessageTypeHello,
		response,
	); err != nil {
		log.Printf(
			"failed to send HELLO response: %v",
			err,
		)

		s.Peers.Remove(hello.NodeID)
		return
	}

	log.Printf(
		"Handshake completed with node %s",
		hello.NodeID,
	)
}

func sendMessage(
	conn net.Conn,
	messageType MessageType,
	payload any,
) error {
	data, err := NewMessage(messageType, payload)
	if err != nil {
		return err
	}

	data = append(data, '\n')

	_, err = conn.Write(data)

	return err
}
