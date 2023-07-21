package main

import (
	"encoding/hex"
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

// BCInstance returns a blockchain instance if any
func BCInstance() *Blockchain {
	if !dbExists() {
		fmt.Println("No existing blockchain found. Create a new one!")
		os.Exit(1)
	}

	var tip []byte
	db, err := bolt.Open(dbFile, 0600, nil)

	if err != nil {
		log.Panic(err)
	}

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

	if err != nil {
		log.Panic(err)
	}

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

			if err != nil {
				log.Panic(err)
			}
			err = b.Put([]byte("l"), genesis.Hash)

			if err != nil {
				log.Panic(err)
			}
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

// FindUnspentTransactions ..
func (bc *Blockchain) FindUnspentTransactions(address string) []Transaction {
	var unspentTxs []Transaction
	spentTXOs := make(map[string][]int)
	bci := bc.Iterator()

	for {
		b := bci.Next()

		for _, tx := range b.Transactions {
			txID := hex.EncodeToString(tx.ID)

		Outputs:
			for outIdx, out := range tx.Vout {
				if spentTXOs[txID] != nil {
					for _, spentOut := range spentTXOs[txID] {
						if spentOut == outIdx {
							continue Outputs
						}
					}
				}
				// Check the owner of this output value
				if out.CanBeUnlockedWith(address) {
					unspentTxs = append(unspentTxs, *tx)
				}
			}

			if !tx.IsCoinBase() {
				for _, in := range tx.Vin {
					if in.CanUnlockOutputWith(address) {
						inTxID := hex.EncodeToString(in.Txid)
						spentTXOs[inTxID] = append(spentTXOs[inTxID], in.Vout)
					}
				}
			}
		}
		if len(b.PrevBlockHash) == 0 {
			break
		}
	}
	return unspentTxs
}

// FindUTXO ..
func (bc *Blockchain) FindUTXO(address string) []TXOutput {
	var UTXOs []TXOutput
	unspentTxs := bc.FindUnspentTransactions(address)

	for _, tx := range unspentTxs {
		for _, out := range tx.Vout {
			if out.CanBeUnlockedWith(address) {
				UTXOs = append(UTXOs, out)
			}
		}
	}
	return UTXOs
}

// MineBlock fn
func (bc *Blockchain) MineBlock(txs []*Transaction) {
	var lastHash []byte

	err := bc.db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(blocksBuket))
		lastHash = b.Get([]byte("l")) // Get lash block hash

		return nil
	})

	if err != nil {
		log.Panic(err)
	}

	newBlock := NewBlock(txs, lastHash)

	err = bc.db.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(blocksBuket))
		err := b.Put(newBlock.Hash, newBlock.Serialize())

		if err != nil {
			log.Panic(err)
		}

		err = b.Put([]byte("l"), newBlock.Hash)

		if err != nil {
			log.Panic(err)
		}
		bc.tip = newBlock.Hash

		return nil
	})
	if err != nil {
		log.Panic(err)
	}
}
