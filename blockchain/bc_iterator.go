package blockchain

import (
	"log"

	"github.com/boltdb/bolt"
)

// Iterator - Iterator implemetation to the blockchain
type Iterator struct {
	currentHash []byte
	db          *bolt.DB
}

// Iterator - Constructor
func (bc *Blockchain) Iterator() *Iterator {
	bci := &Iterator{bc.tip, bc.db}
	return bci
}

// Next function return the next block
func (i *Iterator) Next() *Block {
	if len(i.currentHash) == 0 {
		return nil
	}
	var block *Block

	err := i.db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(blocksBuket))
		encodedBlock := b.Get(i.currentHash)
		block = DeserializeBlock(encodedBlock)

		return nil
	})

	if err != nil {
		log.Panic(err)
	}

	i.currentHash = block.PrevBlockHash
	return block
}
