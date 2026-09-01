package network

import (
	"fmt"
	"sync"
)

type PeerManager struct {
	mu    sync.RWMutex
	peers map[string]*Peer
}

func NewPeerManager() *PeerManager {
	return &PeerManager{
		peers: make(map[string]*Peer),
	}
}

func (pm *PeerManager) Add(id string, peer *Peer) error {
	if pm == nil {
		return fmt.Errorf("peer manager is nil")
	}

	if peer == nil {
		return fmt.Errorf("peer is nil")
	}

	if id == "" {
		return fmt.Errorf("peer ID is empty")
	}

	pm.mu.Lock()
	defer pm.mu.Unlock()

	if _, exists := pm.peers[id]; exists {
		return fmt.Errorf("peer %s already exists", id)
	}

	pm.peers[id] = peer

	return nil
}

func (pm *PeerManager) Get(id string) (*Peer, bool) {
	if pm == nil {
		return nil, false
	}

	pm.mu.RLock()
	defer pm.mu.RUnlock()

	peer, exists := pm.peers[id]

	return peer, exists
}

func (pm *PeerManager) Remove(id string) {
	if pm == nil {
		return
	}

	pm.mu.Lock()
	defer pm.mu.Unlock()

	if peer, exists := pm.peers[id]; exists {
		_ = peer.Close()
	}

	delete(pm.peers, id)
}

func (pm *PeerManager) Count() int {
	if pm == nil {
		return 0
	}

	pm.mu.RLock()
	defer pm.mu.RUnlock()

	return len(pm.peers)
}

func (pm *PeerManager) All() map[string]*Peer {
	if pm == nil {
		return nil
	}

	pm.mu.RLock()
	defer pm.mu.RUnlock()

	result := make(map[string]*Peer, len(pm.peers))

	for id, peer := range pm.peers {
		result[id] = peer
	}

	return result
}

func (pm *PeerManager) Broadcast(
	messageType MessageType,
	payload any,
) {
	if pm == nil {
		return
	}

	pm.mu.RLock()
	defer pm.mu.RUnlock()

	for _, peer := range pm.peers {
		if err := peer.Send(messageType, payload); err != nil {
			continue
		}
	}
}
