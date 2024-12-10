package blockchain

import (
	"bytes"
	"crypto/sha256"
	"strconv"
	"time"
)

// SetHash sets the hash of the block
func (block *Block) SetHash() {
	timestamp := []byte(strconv.FormatInt(block.Timestamp, 10))
	headers := bytes.Join([][]byte{
		timestamp,
		block.PreviousHash,
		block.AllData,
		[]byte(block.IPFSHash),
		[]byte(block.TransactionID),
		[]byte(block.Owner),
	}, []byte{})
	hash := sha256.Sum256(headers)
	block.MyBlockHash = hash[:]
}

// NewBlock creates a new block with the provided data
func NewBlock(AllData string, PreviousHash []byte, IPFSHash, Owner, RecordID string, DoctorsWithPermission []string, Interactions []string) *Block {
	block := &Block{
		Timestamp:             time.Now().Unix(),
		PreviousHash:          PreviousHash,
		AllData:               []byte(AllData),
		IPFSHash:              IPFSHash,
		TransactionID:         RecordID,
		Owner:                 Owner,
		DoctorsWithPermission: DoctorsWithPermission,
		Interactions:          Interactions,
	}
	block.SetHash()
	return block
}

// AddDoctorPermission adds a doctor to the permission list
func (blk *Block) AddDoctorPermission(doctorID string) {
	// Avoid duplicates
	for _, doctor := range blk.DoctorsWithPermission {
		if doctor == doctorID {
			return // Doctor is already in the list
		}
	}
	// Add the doctor to the permission list
	blk.DoctorsWithPermission = append(blk.DoctorsWithPermission, doctorID)
}

// RemoveDoctorPermission removes a doctor from the permission list
func (blk *Block) RemoveDoctorPermission(doctorID string) {
	for i, doctor := range blk.DoctorsWithPermission {
		if doctor == doctorID {
			// Remove the doctor from the slice
			blk.DoctorsWithPermission = append(blk.DoctorsWithPermission[:i], blk.DoctorsWithPermission[i+1:]...)
			return
		}
	}
}

// LogInteraction logs an interaction with the block (e.g., doctor requesting access)
func (blk *Block) LogInteraction(interaction string) {
	blk.Interactions = append(blk.Interactions, interaction)
}

// HasPermission checks if a doctor has permission to access the block
func (blk *Block) HasPermission(doctorID string) bool {
	for _, doctor := range blk.DoctorsWithPermission {
		if doctor == doctorID {
			return true
		}
	}
	return false
}

// GetInteractions returns all interactions logged for the block
func (blk *Block) GetInteractions() []string {
	return blk.Interactions
}
