package main

import (
	"flag"
	"log"
	"strconv"

	"github.com/Mesit-Rathnayake/nexus/internal/network"
	"github.com/Mesit-Rathnayake/nexus/internal/node"
)

func main() {
	port := flag.Int("port", 8001, "port for the Nexus node")

	flag.Parse()

	nodeID := "nexus-node"

	log.Printf(
		"Starting Nexus node %s on port %d",
		nodeID,
		*port,
	)

	n := node.NewNode(
		nodeID,
		"127.0.0.1:"+strconv.Itoa(*port),
	)

	log.Printf(
		"Node created with blockchain containing %d blocks",
		len(n.Blockchain.Blocks),
	)

	server := network.NewServer(*port)

	log.Printf(
		"Nexus node %s is listening for P2P connections",
		nodeID,
	)

	if err := server.Start(); err != nil {
		log.Fatalf("Network server stopped: %v", err)
	}
}
