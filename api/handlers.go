package api

import (
	"encoding/json"
	"log"
	"my-blockchain/blockchain"
	"net/http"

	"github.com/gorilla/mux"
)

var BC = blockchain.NewBlockchain()

func AddRecord(w http.ResponseWriter, r *http.Request) {
	// Decode the incoming JSON request
	var request struct {
		Data     string `json:"data"`
		Owner    string `json:"owner"`
		RecordID string `json:"record_id"`
	}
	err := json.NewDecoder(r.Body).Decode(&request)
	if err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	// Add the file to IPFS
	ipfsHash, err := blockchain.AddFileToIPFS(request.Data)
	if err != nil {
		log.Printf("Error uploading to IPFS: %v", err)
		http.Error(w, "Failed to upload to IPFS", http.StatusInternalServerError)
		return
	}

	// Add the record to the blockchain
	err = BC.AddBlockWithMetadata(request.Data, ipfsHash, request.Owner, request.RecordID)
	if err != nil {
		log.Printf("Error adding block: %v", err)
		http.Error(w, "Failed to add block to blockchain", http.StatusInternalServerError)
		return
	}

	// Respond with success
	w.WriteHeader(http.StatusCreated)
	response := map[string]string{
		"message":  "Record successfully added",
		"cid":      ipfsHash,
		"recordID": request.RecordID,
	}
	json.NewEncoder(w).Encode(response)
}

func GetFileFromIPFS(w http.ResponseWriter, r *http.Request) {
	// Extract CID from request
	vars := mux.Vars(r)
	cid := vars["cid"]
	if cid == "" {
		http.Error(w, "Missing CID parameter", http.StatusBadRequest)
		return
	}

	// Retrieve the file from IPFS
	data, err := blockchain.GetFileFromIPFS(cid)
	if err != nil {
		log.Printf("Error retrieving file from IPFS: %v", err)
		http.Error(w, "Failed to retrieve file from IPFS", http.StatusInternalServerError)
		return
	}

	// Respond with the file content
	response := map[string]string{"cid": cid, "data": data}
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response)
}

func GetBlockchain(w http.ResponseWriter, r *http.Request) {
	// Return the blockchain in JSON format
	json.NewEncoder(w).Encode(BC)
}
