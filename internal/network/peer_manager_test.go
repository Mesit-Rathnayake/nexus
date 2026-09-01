package network

import (
	"net"
	"testing"
)

func TestPeerManager(t *testing.T) {
	manager := NewPeerManager()

	if manager == nil {
		t.Fatal("NewPeerManager() returned nil")
	}

	if manager.Count() != 0 {
		t.Fatalf(
			"initial peer count = %d, want 0",
			manager.Count(),
		)
	}

	server, client := net.Pipe()

	defer server.Close()
	defer client.Close()

	peer := &Peer{
		Address: "test-peer",
		Conn:    client,
	}

	if err := manager.Add("peer-1", peer); err != nil {
		t.Fatalf("Add() failed: %v", err)
	}

	if manager.Count() != 1 {
		t.Fatalf(
			"peer count = %d, want 1",
			manager.Count(),
		)
	}

	got, exists := manager.Get("peer-1")

	if !exists {
		t.Fatal("peer was not found")
	}

	if got != peer {
		t.Fatal("Get() returned wrong peer")
	}
}

func TestPeerManagerRejectsDuplicate(t *testing.T) {
	manager := NewPeerManager()

	server, client := net.Pipe()

	defer server.Close()
	defer client.Close()

	peer := &Peer{
		Address: "test-peer",
		Conn:    client,
	}

	if err := manager.Add("peer-1", peer); err != nil {
		t.Fatalf("first Add() failed: %v", err)
	}

	if err := manager.Add("peer-1", peer); err == nil {
		t.Fatal("manager accepted duplicate peer")
	}
}

func TestPeerManagerRemove(t *testing.T) {
	manager := NewPeerManager()

	server, client := net.Pipe()

	defer server.Close()

	peer := &Peer{
		Address: "test-peer",
		Conn:    client,
	}

	if err := manager.Add("peer-1", peer); err != nil {
		t.Fatalf("Add() failed: %v", err)
	}

	manager.Remove("peer-1")

	if manager.Count() != 0 {
		t.Fatalf(
			"peer count after removal = %d, want 0",
			manager.Count(),
		)
	}

	_, exists := manager.Get("peer-1")

	if exists {
		t.Fatal("removed peer still exists")
	}
}
