package blockchain

import (
	"fmt"
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
	newBlock.SetHash()

	// Validate the block hash before adding it
	if len(newBlock.MyBlockHash) == 0 {
		return fmt.Errorf("block hash generation failed")
	}

	blkch.Blocks = append(blkch.Blocks, newBlock)
	return nil
}
