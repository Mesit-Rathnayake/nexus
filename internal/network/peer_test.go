package network

import (
	"net"
	"testing"
)

func TestPeerSendReceive(t *testing.T) {
	server, client := net.Pipe()

	defer server.Close()
	defer client.Close()

	peer := &Peer{
		Address: "pipe",
		Conn:    client,
	}

	done := make(chan error)

	go func() {
		message, err := Receive(server)

		if err != nil {
			done <- err
			return
		}

		if message.Type != MessageTypePing {
			done <- &testError{
				message: "unexpected message type",
			}
			return
		}

		done <- nil
	}()

	err := peer.Send(
		MessageTypePing,
		map[string]string{
			"message": "hello",
		},
	)

	if err != nil {
		t.Fatalf("Send() failed: %v", err)
	}

	if err := <-done; err != nil {
		t.Fatalf("receiver failed: %v", err)
	}
}

type testError struct {
	message string
}

func (e *testError) Error() string {
	return e.message
}
