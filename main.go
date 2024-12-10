package main

import (
	"context"
	"log"
	"my-blockchain/api"
	"my-blockchain/blockchain"
	"net/http"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	// Initialize Blockchain
	log.Println("Initializing Blockchain...")
	blockchain.InitBlockchain() // Initializes the global BC variable

	// IPFS and Ganache connection code...
	// Connect to IPFS and Ganache as before

	// Step 3: Initialize the API router
	log.Println("Initializing API router...")
	router := api.NewRouter() // Router does not need to pass blockchainInstance

	// Step 4: Setup graceful shutdown
	serverAddress := ":8080"
	srv := &http.Server{
		Addr:    serverAddress,
		Handler: router,
	}

	// Go routine for handling shutdown signals
	go func() {
		sigint := make(chan os.Signal, 1)
		signal.Notify(sigint, syscall.SIGINT, syscall.SIGTERM)
		<-sigint
		log.Println("Shutting down gracefully...")
		if err := srv.Shutdown(context.Background()); err != nil {
			log.Fatalf("Server Shutdown Failed: %v", err)
		}
		log.Println("Server shutdown successfully")
	}()

	// Step 5: Start the HTTP server
	log.Printf("Starting server on %s...\n", serverAddress)
	if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		log.Fatalf("Server failed: %v", err)
	}
}
