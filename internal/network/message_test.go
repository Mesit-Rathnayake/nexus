package network

import (
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
