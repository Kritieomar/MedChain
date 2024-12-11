package blockchain

import (
	"bytes"
	"crypto/sha256"
	"strconv"
	"time"
)

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

func (blk *Block) AddDoctorPermission(doctorID string) {

	for _, doctor := range blk.DoctorsWithPermission {
		if doctor == doctorID {
			return
		}
	}

	blk.DoctorsWithPermission = append(blk.DoctorsWithPermission, doctorID)
}

func (blk *Block) RemoveDoctorPermission(doctorID string) {
	for i, doctor := range blk.DoctorsWithPermission {
		if doctor == doctorID {

			blk.DoctorsWithPermission = append(blk.DoctorsWithPermission[:i], blk.DoctorsWithPermission[i+1:]...)
			return
		}
	}
}

func (blk *Block) LogInteraction(interaction string) {
	blk.Interactions = append(blk.Interactions, interaction)
}

func (blk *Block) HasPermission(doctorID string) bool {
	for _, doctor := range blk.DoctorsWithPermission {
		if doctor == doctorID {
			return true
		}
	}
	return false
}

func (blk *Block) GetInteractions() []string {
	return blk.Interactions
}
