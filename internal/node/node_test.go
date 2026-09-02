package node

import "testing"

func TestNewNode(t *testing.T) {
	node := NewNode("node-1", "127.0.0.1:8001")

	if node == nil {
		t.Fatal("NewNode() returned nil")
	}

	if node.ID != "node-1" {
		t.Fatalf(
			"node ID = %s, want node-1",
			node.ID,
		)
	}

	if node.Address != "127.0.0.1:8001" {
		t.Fatalf(
			"node address = %s, want 127.0.0.1:8001",
			node.Address,
		)
	}

	if node.Blockchain == nil {
		t.Fatal("node blockchain is nil")
	}

	if node.Mempool == nil {
		t.Fatal("node mempool is nil")
	}

	if node.Peers == nil {
		t.Fatal("node peers manager is nil")
	}

	if len(node.Blockchain.Blocks) != 1 {
		t.Fatalf(
			"blockchain contains %d blocks, want 1",
			len(node.Blockchain.Blocks),
		)
	}

	if node.Mempool.Size() != 0 {
		t.Fatalf(
			"mempool size = %d, want 0",
			node.Mempool.Size(),
		)
	}

	if node.Peers == nil {
		t.Fatal("expected peer manager to be initialized")
	}

	if node.Peers.Count() != 0 {
		t.Fatalf(
			"expected empty peer manager, got %d peers",
			node.Peers.Count(),
		)
	}
}
