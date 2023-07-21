package main

import (
	"log"

	"github.com/boltdb/bolt"
)

// BlockchainIterator - Iterator implemetation to the blockchain
type BlockchainIterator struct {
	currentHash []byte
	db          *bolt.DB
}

// Iterator - Constructor
func (bc *Blockchain) Iterator() *BlockchainIterator {
	bci := &BlockchainIterator{bc.tip, bc.db}
	return bci
}

// Next function return the next block
func (i *BlockchainIterator) Next() *Block {
	if len(i.currentHash) == 0 {
		return nil
	}
	var block *Block

	err := i.db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(blocksBuket))
		encodedBlock := b.Get(i.currentHash)
		block = Deserialize(encodedBlock)

		return nil
	})

	if err != nil {
		log.Panic(err)
	}

	i.currentHash = block.PrevBlockHash
	return block
}
