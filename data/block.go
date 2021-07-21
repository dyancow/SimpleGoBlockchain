package data

import (
	"crypto/sha256"
	"encoding/hex"
)

type Block struct {
	Index     int
	Timestamp string
	BPM       int
	Hash      string
	PrevHash  string
}

// CalculateHash takes current block data and calculate its hash
func (block Block) CalculateHash() string {
	record := string(block.Index) + block.Timestamp + string(block.BPM) + block.PrevHash
	h := sha256.New()
	h.Write([]byte(record))
	hashed := h.Sum(nil)
	return hex.EncodeToString(hashed)
}

//IsBlockValid validation for blocks
func (currentBlock Block) IsBlockValid(prevBlock Block) bool {
	//check indices
	if currentBlock.Index != prevBlock.Index+1 {
		return false
	}
	//check hashes
	if currentBlock.PrevHash != prevBlock.Hash {
		return false
	}
	//check current block's hash again
	if currentBlock.Hash != currentBlock.CalculateHash(){
		return false
	}

	return true
}
