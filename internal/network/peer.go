package network

import (
	"bufio"
	"encoding/json"
	"fmt"
	"net"
	"sync"
)

type Peer struct {
	Address string
	Conn    net.Conn
	Reader  *bufio.Reader
	mu      sync.Mutex
}

func NewPeer(address string) (*Peer, error) {
	conn, err := net.Dial("tcp", address)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to peer %s: %w", address, err)
	}

	return &Peer{
		Address: address,
		Conn:    conn,
		Reader:  bufio.NewReader(conn),
	}, nil
}

func (p *Peer) Send(messageType MessageType, payload any) error {
	if p == nil || p.Conn == nil {
		return fmt.Errorf("peer connection is not available")
	}

	data, err := NewMessage(messageType, payload)
	if err != nil {
		return err
	}

	data = append(data, '\n')

	p.mu.Lock()
	defer p.mu.Unlock()

	_, err = p.Conn.Write(data)

	if err != nil {
		return fmt.Errorf("failed to send message: %w", err)
	}

	return nil
}

func Receive(reader *bufio.Reader) (*Message, error) {
	if reader == nil {
		return nil, fmt.Errorf("reader is nil")
	}

	data, err := reader.ReadBytes('\n')
	if err != nil {
		return nil, fmt.Errorf("failed to read message: %w", err)
	}

	var message Message

	if err := json.Unmarshal(data, &message); err != nil {
		return nil, fmt.Errorf("failed to decode message: %w", err)
	}

	if message.Type == "" {
		return nil, fmt.Errorf("message type is empty")
	}

	return &message, nil
}

func (p *Peer) Close() error {
	if p == nil || p.Conn == nil {
		return nil
	}

	return p.Conn.Close()
}

func (p *Peer) Receive() (*Message, error) {
	if p == nil || p.Conn == nil {
		return nil, fmt.Errorf("peer connection is not available")
	}

	if p.Reader == nil {
		p.Reader = bufio.NewReader(p.Conn)
	}

	return Receive(p.Reader)
}
