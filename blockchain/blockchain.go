package blockchain

import (
	"fmt"
	"log"
	"time"
)

func NewBlockchain() *Blockchain {
	genesisBlock := NewBlock("Genesis Block", []byte{}, "", "", "")
	return &Blockchain{
		Blocks: []*Block{genesisBlock},
	}
}

func (blkch *Blockchain) AddBlock(AllData string, IPFSHash string) {
	previousBlock := blkch.Blocks[len(blkch.Blocks)-1]
	newBlock := NewBlock(AllData, previousBlock.MyBlockHash, IPFSHash, "", "")
	blkch.Blocks = append(blkch.Blocks, newBlock)
}

func (blkch *Blockchain) AddBlockWithMetadata(AllData, IPFSHash, Owner, RecordID string) error {
	previousBlock := blkch.Blocks[len(blkch.Blocks)-1]
	newBlock := &Block{
		Timestamp:     time.Now().Unix(),
		PreviousHash:  previousBlock.MyBlockHash,
		AllData:       []byte(AllData),
		IPFSHash:      IPFSHash,
		TransactionID: RecordID,
		Owner:         Owner,
	}

	log.Printf("Adding block with Data: %s, IPFSHash: %s, RecordID: %s, Owner: %s", AllData, IPFSHash, RecordID, Owner)
	newBlock.SetHash()

	if len(newBlock.MyBlockHash) == 0 {
		return fmt.Errorf("block hash generation failed")
	}

	blkch.Blocks = append(blkch.Blocks, newBlock)
	return nil
}
