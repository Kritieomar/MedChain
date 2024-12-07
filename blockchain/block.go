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

func NewBlock(AllData string, PreviousHash []byte, IPFSHash, Owner, RecordID string) *Block {
	block := &Block{
		Timestamp:     time.Now().Unix(),
		PreviousHash:  PreviousHash,
		AllData:       []byte(AllData),
		IPFSHash:      IPFSHash,
		TransactionID: RecordID,
		Owner:         Owner,
	}
	block.SetHash()
	return block
}
