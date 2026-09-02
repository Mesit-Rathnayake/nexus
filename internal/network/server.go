package network

import (
	"bufio"
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
	log.Printf(
		"Peer connected: %s",
		conn.RemoteAddr(),
	)

	reader := bufio.NewReader(conn)

	message, err := Receive(reader)
	if err != nil {
		log.Printf(
			"failed to receive initial message: %v",
			err,
		)

		_ = conn.Close()
		return
	}

	if message.Type != MessageTypeHello {
		log.Printf(
			"expected HELLO, got %s",
			message.Type,
		)

		_ = conn.Close()
		return
	}

	peerID, err := s.handleHello(conn, message)
	if err != nil {
		log.Printf(
			"handshake failed: %v",
			err,
		)

		_ = conn.Close()
		return
	}

	peer := &Peer{
		Address: conn.RemoteAddr().String(),
		Conn:    conn,
		Reader:  reader,
	}

	if err := s.Peers.Replace(peerID, peer); err != nil {
		log.Printf(
			"failed to register persistent peer: %v",
			err,
		)

		_ = conn.Close()
		return
	}

	log.Printf(
		"Starting message loop for peer %s",
		peerID,
	)

	s.messageLoop(peer)
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
) (string, error) {
	var hello HelloPayload

	if err := json.Unmarshal(
		message.Payload,
		&hello,
	); err != nil {
		return "", fmt.Errorf(
			"failed to decode HELLO payload: %w",
			err,
		)
	}

	if hello.NodeID == "" {
		return "", fmt.Errorf(
			"HELLO contains empty node ID",
		)
	}

	if hello.Address == "" {
		return "", fmt.Errorf(
			"HELLO contains empty node address",
		)
	}

	log.Printf(
		"Handshake from node %s at %s",
		hello.NodeID,
		hello.Address,
	)

	response := HelloPayload{
		NodeID:  s.NodeID,
		Address: s.NodeAddress,
	}

	if err := sendMessage(
		conn,
		MessageTypeHello,
		response,
	); err != nil {
		return "", fmt.Errorf(
			"failed to send HELLO response: %w",
			err,
		)
	}

	log.Printf(
		"Handshake completed with node %s",
		hello.NodeID,
	)

	return hello.NodeID, nil
}

func (s *Server) messageLoop(peer *Peer) {
	for {
		message, err := peer.Receive()
		if err != nil {
			log.Printf(
				"connection to %s closed: %v",
				peer.Address,
				err,
			)

			return
		}

		log.Printf(
			"Received %s from %s",
			message.Type,
			peer.Address,
		)

		s.handleMessage(peer, message)
	}
}

func (s *Server) handleMessage(
	peer *Peer,
	message *Message,
) {
	switch message.Type {

	case MessageTypePing:
		if err := peer.Send(
			MessageTypePong,
			map[string]string{
				"message": "hello from Nexus",
			},
		); err != nil {
			log.Printf(
				"failed to send PONG: %v",
				err,
			)
		}

	case MessageTypeHello:
		log.Printf(
			"unexpected HELLO from %s",
			peer.Address,
		)

	default:
		log.Printf(
			"Unhandled message type: %s",
			message.Type,
		)
	}
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
