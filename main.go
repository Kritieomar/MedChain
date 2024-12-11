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

	log.Println("Initializing Blockchain...")
	blockchain.InitBlockchain()

	log.Println("Initializing API router...")
	router := api.NewRouter()

	serverAddress := ":8080"
	srv := &http.Server{
		Addr:    serverAddress,
		Handler: router,
	}

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

	log.Printf("Starting server on %s...\n", serverAddress)
	if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		log.Fatalf("Server failed: %v", err)
	}
}
