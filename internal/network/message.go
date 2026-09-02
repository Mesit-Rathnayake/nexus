package network

import (
	"encoding/json"
	"fmt"
)

type MessageType string

const (
	MessageTypePing           MessageType = "PING"
	MessageTypePong           MessageType = "PONG"
	MessageTypeHello          MessageType = "HELLO"
	MessageTypeNewTransaction MessageType = "NEW_TRANSACTION"
	MessageTypeNewBlock       MessageType = "NEW_BLOCK"
	MessageTypeGetChain       MessageType = "GET_CHAIN"
	MessageTypeChainResponse  MessageType = "CHAIN_RESPONSE"
)

type Message struct {
	Type    MessageType `json:"type"`
	Payload []byte      `json:"payload"`
}

func NewMessage(messageType MessageType, payload any) ([]byte, error) {
	data, err := json.Marshal(payload)
	if err != nil {
		return nil, fmt.Errorf("marshal payload: %w", err)
	}

	message := Message{
		Type:    messageType,
		Payload: data,
	}

	return json.Marshal(message)
}

func DecodeMessage(data []byte) (*Message, error) {
	var message Message

	if err := json.Unmarshal(data, &message); err != nil {
		return nil, fmt.Errorf("decode message: %w", err)
	}

	if message.Type == "" {
		return nil, fmt.Errorf("message type is empty")
	}

	return &message, nil
}
