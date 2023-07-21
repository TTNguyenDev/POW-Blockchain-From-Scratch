package main

import (
	"fmt"
	"log"
	"os"

	"github.com/boltdb/bolt"
)

const dbFile = "blockchain.db"
const blocksBuket = "blocks"

// Blockchain struct
type Blockchain struct {
	tip []byte
	db  *bolt.DB
}

func dbExists() bool {
	if _, err := os.Stat(dbFile); os.IsNotExist(err) {
		return false
	}
	return true
}

func BCInstance() *Blockchain {
	if !dbExists() {
		fmt.Println("No existing blockchain found. Create a new one!")
		os.Exit(1)
	}

	var tip []byte
	db, err := bolt.Open(dbFile, 0600, nil)

	err = db.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(blocksBuket))
		tip = b.Get([]byte("l"))

		if tip == nil {
			log.Panic("Latest blockhash doesn't exist")
		}

		return nil
	})

	if err != nil {
		log.Panic(err)
	}

	bc := Blockchain{tip, db}
	return &bc
}

// NewBlockchain fn
func NewBlockchain(benefician string) *Blockchain {
	if dbExists() {
		fmt.Println("Blockchain already exists.")
		os.Exit(1)
	}
	var tip []byte
	db, err := bolt.Open(dbFile, 0600, nil)

	err = db.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(blocksBuket))

		if b != nil {
			tip = b.Get([]byte("l")) //Last block hash
		} else {
			coinbaseTx := NewCoinbaseTX(benefician, "")
			genesis := NewGenesisBlock(coinbaseTx)
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

	if err != nil {
		log.Panic(err)
	}

	bc := Blockchain{tip, db}
	return &bc
}

// AddBlock fn
func (bc *Blockchain) AddBlock(txs []*Transaction, data string) {
	var lastHash []byte

	err := bc.db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(blocksBuket))
		lastHash = b.Get([]byte("l")) // Get lash block hash

		return nil
	})

	if err != nil {
		log.Panic(err)
	}

	newBlock := NewBlock(txs, data, lastHash)

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
