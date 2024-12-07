package blockchain

type Block struct {
	Timestamp     int64  `json:"timestamp"`
	PreviousHash  []byte `json:"prevBlockHash"`
	MyBlockHash   []byte `json:"hash"`
	AllData       []byte `json:"medicalData"`
	IPFSHash      string `json:"ipfsHash"`
	TransactionID string `json:"recordId"`
	Owner         string `json:"owner"`
}

type Blockchain struct {
	Blocks []*Block `json:"blocks"`
}
