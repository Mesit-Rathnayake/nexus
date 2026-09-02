package network

import (
	"encoding/json"
	"fmt"
	"log"
)

type Network struct {
	NodeID      string
	NodeAddress string
	Peers       *PeerManager
}

func NewNetwork(
	nodeID string,
	nodeAddress string,
	peers *PeerManager,
) *Network {
	return &Network{
		NodeID:      nodeID,
		NodeAddress: nodeAddress,
		Peers:       peers,
	}
}

func (n *Network) Connect(address string) error {
	if n == nil {
		return fmt.Errorf("network is nil")
	}

	if address == "" {
		return fmt.Errorf("peer address is empty")
	}

	if n.Peers == nil {
		return fmt.Errorf("peer manager is nil")
	}

	log.Printf("Connecting to peer at %s", address)

	peer, err := NewPeer(address)
	if err != nil {
		return err
	}

	hello := HelloPayload{
		NodeID:  n.NodeID,
		Address: n.NodeAddress,
	}

	if err := peer.Send(MessageTypeHello, hello); err != nil {
		_ = peer.Close()

		return fmt.Errorf("failed to send HELLO: %w", err)
	}

	log.Printf("HELLO sent to %s", address)

	message, err := peer.Receive()
	if err != nil {
		_ = peer.Close()

		return fmt.Errorf("failed to receive HELLO response: %w", err)
	}

	if message.Type != MessageTypeHello {
		_ = peer.Close()

		return fmt.Errorf(
			"expected HELLO response, got %s",
			message.Type,
		)
	}

	var response HelloPayload

	if err := decodeHello(message, &response); err != nil {
		_ = peer.Close()

		return err
	}

	if response.NodeID == "" {
		_ = peer.Close()

		return fmt.Errorf("peer returned empty node ID")
	}

	if err := n.Peers.Add(response.NodeID, peer); err != nil {
		_ = peer.Close()

		return fmt.Errorf(
			"failed to register peer %s: %w",
			response.NodeID,
			err,
		)
	}

	log.Printf(
		"Handshake completed with %s at %s",
		response.NodeID,
		response.Address,
	)

	return nil
}

func decodeHello(
	message *Message,
	payload *HelloPayload,
) error {
	if message == nil {
		return fmt.Errorf("message is nil")
	}

	if payload == nil {
		return fmt.Errorf("payload is nil")
	}

	if err := json.Unmarshal(message.Payload, payload); err != nil {
		return fmt.Errorf(
			"failed to decode HELLO payload: %w",
			err,
		)
	}

	return nil
}
