package blockchain

import (
	"encoding/hex"
	"log"

	"github.com/boltdb/bolt"

	"blockchain_from_scratch/blockchain/transaction"
)

const utxoSetBucket = "chainstate"

type UTXOSet struct {
	Bc *Blockchain
}

func (u UTXOSet) Reindex() {
	db := u.Bc.db
	bucket := []byte(utxoSetBucket)

	err := db.Update(func(tx *bolt.Tx) error {
		err := tx.DeleteBucket(bucket)
		if err != nil {
			log.Panic(err)
		}
		_, err = tx.CreateBucket(bucket)
		if err != nil {
			return err
		}
		return nil
	})

	if err != nil {
		log.Panic(err)
	}

	allUTXOs := u.Bc.FindUTXO()

	err = db.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket(bucket)

		for txId, outputs := range allUTXOs {
			key, err := hex.DecodeString(txId)
			if err != nil {
				log.Panic(err)
			}
			err = b.Put(key, transaction.SerializeTXOutputs(outputs))

			if err != nil {
				log.Panic(err)
			}
		}
		return nil
	})

	if err != nil {
		log.Panic(err)
	}
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

	if err != nil {
		log.Panic(err)
	}
	return accumulated, txs
}

// FindUTXO
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

	if err != nil {
		log.Panic(err)
	}
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

					for outId, out := range outs {
						if outId != vin.Vout {
							updateOuts = append(updateOuts, out)
						}
					}
					if len(updateOuts) == 0 {
						err := b.Delete(vin.Txid)

						if err != nil {
							log.Panic(err)
						}
					} else {
						err := b.Put(vin.Txid, transaction.SerializeTXOutputs(updateOuts))

						if err != nil {
							log.Panic(err)
						}
					}
				}
			}

			var newOutputs []transaction.TXOutput
			newOutputs = append(newOutputs, tx.Vout...)
			err := b.Put(tx.ID, transaction.SerializeTXOutputs(newOutputs))

			if err != nil {
				log.Panic(err)
			}
		}
		return nil
	})

	if err != nil {
		log.Panic(err)
	}
}
