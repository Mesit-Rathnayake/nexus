package network

import (
	"net"
	"testing"
	"time"
)

func TestServerPingPong(t *testing.T) {
	listener, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		t.Fatalf("failed to create listener: %v", err)
	}

	defer listener.Close()

	go func() {
		conn, err := listener.Accept()
		if err != nil {
			return
		}

		server := &Server{
			Address: listener.Addr().String(),
		}

		server.handleConnection(conn)
	}()

	peer, err := NewPeer(listener.Addr().String())
	if err != nil {
		t.Fatalf("failed to connect to server: %v", err)
	}

	defer peer.Close()

	if err := peer.Send(
		MessageTypePing,
		map[string]string{
			"message": "hello nexus",
		},
	); err != nil {
		t.Fatalf("failed to send PING: %v", err)
	}

	// Give the server a moment to respond.
	_ = peer.Conn.SetReadDeadline(
		time.Now().Add(2 * time.Second),
	)

	response, err := peer.Receive()
	if err != nil {
		t.Fatalf("failed to receive PONG: %v", err)
	}

	if response.Type != MessageTypePong {
		t.Fatalf(
			"response type = %s, want %s",
			response.Type,
			MessageTypePong,
		)
	}
}
