package blockchain

import (
	"encoding/hex"
	"log"

	"github.com/boltdb/bolt"

	"blockchain_from_scratch/blockchain/transaction"
	"blockchain_from_scratch/utils"
)

const utxoSetBucket = "chainstate"

// UTXOSet represent the UTXOSet model of Bitcoin
type UTXOSet struct {
	Bc *Blockchain
}

// Reindex checks all blocks in the blockchain staring from genesis block and finds all UTXO transactions
func (u UTXOSet) Reindex() {
	db := u.Bc.db
	bucket := []byte(utxoSetBucket)

	err := db.Update(func(tx *bolt.Tx) error {
		err := tx.DeleteBucket(bucket)
		if err != nil && err != bolt.ErrBucketNotFound {
			log.Panic(err)
		}
		_, err = tx.CreateBucket(bucket)
		utils.CheckError(err)
		return nil
	})
	utils.CheckError(err)

	allUTXOs := u.Bc.FindUTXO()

	err = db.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket(bucket)

		for txID, outputs := range allUTXOs {
			key, err := hex.DecodeString(txID)
			utils.CheckError(err)
			err = b.Put(key, utils.GobEncode(outputs))
			utils.CheckError(err)
		}
		return nil
	})
	utils.CheckError(err)
}

// FindSpendableTransactions -
func (u UTXOSet) FindSpendableTransactions(pubHash []byte, amount int) (int, map[string][]int) {
	txs := make(map[string][]int)
	// utxos := bc.FindUnspentTransactions(pubHash)
	accumulated := 0
	db := u.Bc.db

	err := db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(utxoSetBucket))
		c := b.Cursor()

		for k, v := c.First(); k != nil; k, v = c.Next() {
			txID := hex.EncodeToString(k)
			outputs := transaction.DeSerializeTXOuputs(v)

			for outIdx, out := range outputs {
				if out.IsLockedWithKey(pubHash) && accumulated < amount {
					accumulated += out.Value
					txs[txID] = append(txs[txID], outIdx)
				}
			}
		}
		return nil
	})
	utils.CheckError(err)
	return accumulated, txs
}

// FindUTXO ..
func (u UTXOSet) FindUTXO(pubHash []byte) []transaction.TXOutput {
	var result []transaction.TXOutput
	db := u.Bc.db

	err := db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(utxoSetBucket))
		c := b.Cursor()

		for k, v := c.First(); k != nil; k, v = c.Next() {
			outputs := transaction.DeSerializeTXOuputs(v)

			for _, out := range outputs {
				if out.IsLockedWithKey(pubHash) {
					result = append(result, out)
				}
			}
		}
		return nil
	})

	utils.CheckError(err)
	return result
}

// Update ..
func (u UTXOSet) Update(block *Block) {
	db := u.Bc.db

	err := db.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(utxoSetBucket))

		for _, tx := range block.Transactions {
			if !tx.IsCoinBase() {
				for _, vin := range tx.Vin {
					var updateOuts []transaction.TXOutput
					outBytes := b.Get(vin.Txid)
					outs := transaction.DeSerializeTXOuputs(outBytes)

					for outID, out := range outs {
						if outID != vin.Vout {
							updateOuts = append(updateOuts, out)
						}
					}
					if len(updateOuts) == 0 {
						err := b.Delete(vin.Txid)
						utils.CheckError(err)
					} else {
						err := b.Put(vin.Txid, utils.GobEncode(updateOuts))
						utils.CheckError(err)
					}
				}
			}

			var newOutputs []transaction.TXOutput
			newOutputs = append(newOutputs, tx.Vout...)
			err := b.Put(tx.ID, utils.GobEncode(newOutputs))
			utils.CheckError(err)
		}
		return nil
	})
	utils.CheckError(err)
}
