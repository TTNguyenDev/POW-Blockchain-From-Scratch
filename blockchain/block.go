package blockchain

import (
	"bytes"
	"encoding/gob"
	"log"
	"time"

	"blockchain_from_scratch/blockchain/transaction"
)

// Block - A definition for this simple implementation
type Block struct {
	//TODO: Use merkle tree to store transactions
	Transactions  []*transaction.Transaction
	TimeStamp     int64
	PrevBlockHash []byte
	Hash          []byte
	Nonce         int
}

// NewBlock creates a new block with given txs
func NewBlock(txs []*transaction.Transaction, prevBlockHash []byte) *Block {
	block := &Block{txs, time.Now().Unix(), prevBlockHash, []byte{}, 0}
	pow := NewProofOfWork(block)
	nonce, hash := pow.Run()
	block.Hash = hash
	block.Nonce = nonce
	return block
}

// NewGenesisBlock creates genesis data for the blockchain
func NewGenesisBlock(coinbase *transaction.Transaction) *Block {
	return NewBlock([]*transaction.Transaction{coinbase}, []byte{})
}

// Serialize data
func (b *Block) Serialize() []byte {
	var result bytes.Buffer
	encoder := gob.NewEncoder(&result)

	err := encoder.Encode(b)
	if err != nil {
		log.Panic(err)
	}

	return result.Bytes()
}

// Deserialize data
func Deserialize(d []byte) *Block {
	var block Block

	decoder := gob.NewDecoder(bytes.NewReader(d))
	err := decoder.Decode(&block)

	if err != nil {
		log.Panic(err)
	}

	return &block
}

// HashTransactions creates a single sha256 of all hash sha256 transactions
func (b *Block) HashTransactions() []byte {
	var txs [][]byte

	for _, tx := range b.Transactions {
		txs = append(txs, tx.Serialize())
	}

	mTree := transaction.NewMerkleTree(txs)
	return mTree.RootNode.Data
}
