package blockchain

import (
	"log"

	"github.com/ethereum/go-ethereum/ethclient"
)

func ConnectToGanache() *ethclient.Client {
	client, err := ethclient.Dial("http://127.0.0.1:7545")
	if err != nil {
		log.Fatalf("Failed to connect to Ganache: %v", err)
	}
	log.Println("Connected to Ganache!")
	return client
}
