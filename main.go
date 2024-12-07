package main

import (
	"log"
	"my-blockchain/api"
	"my-blockchain/blockchain"
	"net/http"
)

func main() {
	// Step 1: Initialize IPFS
	log.Println("Connecting to IPFS...")
	blockchain.ConnectToIPFS()
	log.Println("Connected to IPFS!")

	// Step 2: Connect to Ganache
	log.Println("Connecting to Ganache...")
	ethClient := blockchain.ConnectToGanache()
	defer ethClient.Close() // Ensure the connection is closed when the application exits
	log.Println("Connected to Ganache!")

	// Step 3: Initialize the API router
	log.Println("Initializing API router...")
	router := api.NewRouter()

	// Step 4: Start the HTTP server
	serverAddress := ":8080"
	log.Printf("Starting server on %s...\n", serverAddress)
	log.Fatal(http.ListenAndServe(serverAddress, router))
}
