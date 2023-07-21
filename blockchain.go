package main

import (
	"log"

	"github.com/boltdb/bolt"
)

const dbFile = "blockchain.db"
const blocksBuket = "blocks"

// Blockchain struct
type Blockchain struct {
	tip []byte
	db  *bolt.DB
}

// NewBlockchain fn
func NewBlockchain() *Blockchain {
	var tip []byte
	db, err := bolt.Open(dbFile, 0600, nil)

	err = db.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(blocksBuket))

		if err != nil {
			log.Panic(err)
		}

		if b != nil {
			tip = b.Get([]byte("l")) //Last block hash
		} else {
			genesis := NewGenesisBlock()
			b, err := tx.CreateBucket([]byte(blocksBuket))
			if err != nil {
				log.Panic(err)
			}
			err = b.Put(genesis.Hash, genesis.Serialize())
			err = b.Put([]byte("l"), genesis.Hash)
			tip = genesis.Hash
		}
		return nil
	})

	bc := Blockchain{tip, db}
	return &bc
}

// AddBlock fn
func (bc *Blockchain) AddBlock(data string) {
	var lastHash []byte

	err := bc.db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(blocksBuket))
		lastHash = b.Get([]byte("l")) // Get lash block hash

		return nil
	})

	if err != nil {
		log.Panic(err)
	}

	newBlock := NewBlock(data, lastHash)

	err = bc.db.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(blocksBuket))
		err := b.Put(newBlock.Hash, newBlock.Serialize())

		if err != nil {
			log.Panic(err)
		}

		err = b.Put([]byte("l"), newBlock.Hash)
		bc.tip = newBlock.Hash

		return nil
	})
}
