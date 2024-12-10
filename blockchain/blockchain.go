package blockchain

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"
)

var BC *Blockchain

var DoctorProfiles = make(map[string]DoctorProfile)

// PatientProfiles maintains a map of all patient profiles in the system
var PatientProfiles = make(map[string]PatientProfile)

func InitBlockchain() {
	BC = NewBlockchain() // Initialize the global blockchain instance
	ConnectToIPFS()      // Ensure IPFS is connected
}

// NewBlockchain creates a new blockchain with the genesis block
func NewBlockchain() *Blockchain {
	genesisBlock := NewBlock("Genesis Block", []byte{}, "", "", "", []string{}, []string{})
	return &Blockchain{
		Blocks: []*Block{genesisBlock},
	}
}

// AddBlock adds a new block to the blockchain
func (blkch *Blockchain) AddBlock(AllData string, IPFSHash string) error {
	previousBlock := blkch.Blocks[len(blkch.Blocks)-1]
	newBlock := NewBlock(AllData, previousBlock.MyBlockHash, IPFSHash, "", "", []string{}, []string{})

	log.Printf("Adding block with Data: %s, IPFSHash: %s", AllData, IPFSHash)
	newBlock.SetHash()

	if len(newBlock.MyBlockHash) == 0 {
		return fmt.Errorf("block hash generation failed")
	}

	blkch.Blocks = append(blkch.Blocks, newBlock)
	return nil
}

func (blkch *Blockchain) AddBlockWithMetadata(AllData, IPFSHash, Owner, RecordID string, DoctorsWithPermission []string, Interactions []string, PatientID string) error {
	previousBlock := blkch.Blocks[len(blkch.Blocks)-1]
	newBlock := &Block{
		Timestamp:             time.Now().Unix(),
		PreviousHash:          previousBlock.MyBlockHash,
		AllData:               []byte(AllData),
		IPFSHash:              IPFSHash,
		TransactionID:         RecordID,
		Owner:                 Owner,
		DoctorsWithPermission: DoctorsWithPermission,
		Interactions:          Interactions,
		PatientProfiles:       make(map[string]PatientProfile),
	}

	patientProfile := PatientProfile{
		PatientID:          PatientID,
		AcceptedRecords:    []string{},
		RejectedRecords:    []string{},
		InteractionHistory: []string{},
	}

	newBlock.PatientProfiles[PatientID] = patientProfile
	newBlock.SetHash()

	if len(newBlock.MyBlockHash) == 0 {
		return fmt.Errorf("block hash generation failed")
	}

	blkch.Blocks = append(blkch.Blocks, newBlock)
	return nil
}

// GetAllDoctorPermissions retrieves a map of all doctor permissions across the blockchain
func (blkch *Blockchain) GetAllDoctorPermissions() map[string][]string {
	permissions := make(map[string][]string)

	// Iterate through all blocks and store doctor permissions
	for _, block := range blkch.Blocks {
		for _, doctor := range block.DoctorsWithPermission {
			permissions[doctor] = append(permissions[doctor], block.TransactionID)
		}
	}

	return permissions
}

// ValidateChain checks the integrity of the blockchain
func (blkch *Blockchain) ValidateChain() error {
	for i := 1; i < len(blkch.Blocks); i++ {
		currentBlock := blkch.Blocks[i]
		previousBlock := blkch.Blocks[i-1]

		// Check hash consistency between blocks
		if string(currentBlock.PreviousHash) != string(previousBlock.MyBlockHash) {
			return fmt.Errorf("blockchain validation failed: block %d has an invalid previous hash", i)
		}

		// Recalculate the hash and compare it to the stored hash
		expectedHash := currentBlock.MyBlockHash
		currentBlock.SetHash() // Recalculate hash
		if string(currentBlock.MyBlockHash) != string(expectedHash) {
			return fmt.Errorf("blockchain validation failed: block %d has been tampered with", i)
		}
	}
	return nil
}

func GrantPermission(w http.ResponseWriter, r *http.Request) {
	var request struct {
		RecordID  string `json:"record_id"`
		DoctorID  string `json:"doctor_id"`
		PatientID string `json:"patient_id"`
	}

	err := json.NewDecoder(r.Body).Decode(&request)
	if err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	for _, block := range BC.Blocks {
		if block.TransactionID == request.RecordID {
			block.DoctorsWithPermission = append(block.DoctorsWithPermission, request.DoctorID)
			w.WriteHeader(http.StatusOK)
			json.NewEncoder(w).Encode(map[string]string{
				"message":   "Permission granted successfully",
				"recordID":  request.RecordID,
				"doctorID":  request.DoctorID,
				"patientID": request.PatientID,
			})
			return
		}
	}
	http.Error(w, "Record not found", http.StatusNotFound)
}
