package network

import (
	"encoding/json"
	"testing"
)

func TestNewMessage(t *testing.T) {
	payload := map[string]string{
		"message": "hello nexus",
	}

	data, err := NewMessage(MessageTypePing, payload)

	if err != nil {
		t.Fatalf("NewMessage() failed: %v", err)
	}

	if len(data) == 0 {
		t.Fatal("NewMessage() returned empty data")
	}

	message, err := DecodeMessage(data)

	if err != nil {
		t.Fatalf("DecodeMessage() failed: %v", err)
	}

	if message.Type != MessageTypePing {
		t.Fatalf(
			"message type = %s, want %s",
			message.Type,
			MessageTypePing,
		)
	}

	if len(message.Payload) == 0 {
		t.Fatal("message payload is empty")
	}
}

func TestDecodeMessageRejectsInvalidData(t *testing.T) {
	_, err := DecodeMessage([]byte("this isn't json"))

	if err == nil {
		t.Fatal("DecodeMessage() accepted invalid data")
	}
}

func TestDecodeMessageRejectsMissingType(t *testing.T) {
	data := []byte(`{"payload":"hello"}`)

	_, err := DecodeMessage(data)

	if err == nil {
		t.Fatal("DecodeMessage() accepted message without type")
	}
}

func TestMessagePayloadRoundTrip(t *testing.T) {
	type TestPayload struct {
		NodeID string `json:"node_id"`
		Port   int    `json:"port"`
	}

	original := TestPayload{
		NodeID: "node-1",
		Port:   8080,
	}

	data, err := NewMessage(MessageTypePing, original)

	if err != nil {
		t.Fatalf("NewMessage() failed: %v", err)
	}

	message, err := DecodeMessage(data)

	if err != nil {
		t.Fatalf("DecodeMessage() failed: %v", err)
	}

	var decoded TestPayload

	if err := json.Unmarshal(message.Payload, &decoded); err != nil {
		t.Fatalf("payload decoding failed: %v", err)
	}

	if decoded.NodeID != original.NodeID {
		t.Fatalf(
			"NodeID = %s, want %s",
			decoded.NodeID,
			original.NodeID,
		)
	}

	if decoded.Port != original.Port {
		t.Fatalf(
			"Port = %d, want %d",
			decoded.Port,
			original.Port,
		)
	}
}
