package api

import (
	"encoding/json"
	"log"
	"my-blockchain/blockchain"
	"net/http"

	"github.com/gorilla/mux"
)

var BC = blockchain.NewBlockchain()

// AddRecord handles POST requests to add medical records
func AddRecord(w http.ResponseWriter, r *http.Request) {
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

	// Upload to IPFS
	ipfsHash, err := blockchain.AddFileToIPFS(request.Data)
	if err != nil {
		http.Error(w, "Failed to upload to IPFS", http.StatusInternalServerError)
		return
	}

	// Add block to the blockchain
	err = BC.AddBlockWithMetadata(request.Data, ipfsHash, request.Owner, request.RecordID)
	if err != nil {
		log.Printf("Error adding block to blockchain: %v", err)
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

// GetFileFromIPFS handles GET requests to retrieve a file from IPFS using its CID
func GetFileFromIPFS(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	cid := vars["cid"]

	if cid == "" {
		http.Error(w, "Missing CID parameter", http.StatusBadRequest)
		return
	}

	// Retrieve data from IPFS
	data, err := blockchain.GetFileFromIPFS(cid)
	if err != nil {
		log.Printf("Error retrieving file from IPFS for CID %s: %v", cid, err)
		http.Error(w, "Failed to retrieve file from IPFS", http.StatusInternalServerError)
		return
	}

	// Return data and CID as a response
	response := map[string]string{"cid": cid, "data": data}
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response)
}

// GetBlockchain handles GET requests to view the current state of the blockchain
func GetBlockchain(w http.ResponseWriter, r *http.Request) {
	// Return the blockchain
	json.NewEncoder(w).Encode(BC)
}
