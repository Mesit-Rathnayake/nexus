package main

import (
	"flag"
	"log"
	"strconv"

	"github.com/Mesit-Rathnayake/nexus/internal/network"
	"github.com/Mesit-Rathnayake/nexus/internal/node"
)

func main() {
	port := flag.Int(
		"port",
		8001,
		"port for the Nexus node",
	)

	nodeID := flag.String(
		"id",
		"nexus-node",
		"unique ID for the Nexus node",
	)

	flag.Parse()

	address := "127.0.0.1:" + strconv.Itoa(*port)

	log.Printf(
		"Starting Nexus node %s on %s",
		*nodeID,
		address,
	)

	n := node.NewNode(
		*nodeID,
		address,
	)

	log.Printf(
		"Node created with blockchain containing %d blocks",
		len(n.Blockchain.Blocks),
	)

	log.Printf(
		"Peer manager initialized with %d peers",
		n.Peers.Count(),
	)

	server := network.NewServer(
		*port,
		*nodeID,
		address,
		n.Peers,
	)

	log.Printf(
		"Nexus node %s is listening for P2P connections",
		*nodeID,
	)

	if err := server.Start(); err != nil {
		log.Fatalf("Network server stopped: %v", err)
	}
}
