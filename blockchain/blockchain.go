package blockchain

import (
	"bytes"
	"crypto/ecdsa"
	"encoding/hex"
	"errors"
	"fmt"
	"log"
	"os"

	"github.com/boltdb/bolt"

	"blockchain_from_scratch/blockchain/transaction"
	"blockchain_from_scratch/wallet"
)

const dbFile = "db/blockchain.db"
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

func (bc *Blockchain) DB() *bolt.DB {
	return bc.db
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
			coinbaseTx := transaction.NewCoinbaseTX(benefician, "")
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
func (bc *Blockchain) FindUnspentTransactions(pubHash []byte) []transaction.Transaction {
	var unspentTxs []transaction.Transaction
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
				if out.IsLockedWithKey(pubHash) {
					unspentTxs = append(unspentTxs, *tx)
				}
			}

			if !tx.IsCoinBase() {
				for _, in := range tx.Vin {
					if in.UsesKey(pubHash) {
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

// FindUTXO finds all unspent tx outputs and returns Transactions with spent outputs removed
func (bc *Blockchain) FindUTXO() map[string][]transaction.TXOutput {
	UTXO := make(map[string][]transaction.TXOutput)
	spentTXOs := make(map[string][]int)
	bci := bc.Iterator()

	for {
		block := bci.Next()

		for _, tx := range block.Transactions {
			txID := hex.EncodeToString(tx.ID)

		Outputs:
			for outIdx, out := range tx.Vout {
				if spentTXOs[txID] != nil {
					for _, spentOutIdx := range spentTXOs[txID] {
						if spentOutIdx == outIdx {
							continue Outputs
						}
					}
				}
				outs := UTXO[txID]
				outs = append(outs, out)
			}

			if !tx.IsCoinBase() {
				for _, in := range tx.Vin {
					txID := hex.EncodeToString(in.Txid)
					spentTXOs[txID] = append(spentTXOs[txID], in.Vout)
				}
			}
		}
		if len(block.PrevBlockHash) == 0 {
			break
		}
	}
	return UTXO
}

// MineBlock fn
func (bc *Blockchain) MineBlock(txs []*transaction.Transaction) {
	var lastHash []byte

	err := bc.db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(blocksBuket))
		lastHash = b.Get([]byte("l")) // Get lash block hash

		return nil
	})

	if err != nil {
		log.Panic(err)
	}

	for _, tx := range txs {
		if !bc.VerifyTransaction(tx) {
			log.Panic("Error: Invalid Transaction")
		}
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

func (bc *Blockchain) FindTransaction(ID []byte) (transaction.Transaction, error) {
	bci := bc.Iterator()

	for {
		block := bci.Next()
		for _, tx := range block.Transactions {
			if bytes.Equal(tx.ID, ID) /*equally*/ {
				return *tx, nil
			}
		}

		if len(block.PrevBlockHash) == 0 {
			break
		}
	}
	return transaction.Transaction{}, errors.New("Transaction is not found")
}

func (bc *Blockchain) SignTransaction(tx *transaction.Transaction, privKey ecdsa.PrivateKey) {
	prevTXs := make(map[string]transaction.Transaction)

	for _, vin := range tx.Vin {
		prevTX, err := bc.FindTransaction(vin.Txid)
		if err != nil {
			log.Panic(err)
		}
		prevTXs[hex.EncodeToString(prevTX.ID)] = prevTX
	}
	tx.Sign(privKey, prevTXs)
}

func (bc *Blockchain) VerifyTransaction(tx *transaction.Transaction) bool {
	prevTXs := make(map[string]transaction.Transaction)

	for _, vin := range tx.Vin {
		prevTX, err := bc.FindTransaction(vin.Txid)
		if err != nil {
			log.Panic(err)
		}
		prevTXs[hex.EncodeToString(prevTX.ID)] = prevTX
	}
	return tx.Verify(prevTXs)
}

// NewUTXOTransaction -
func NewUTXOTransaction(from, to string, amount int, u *UTXOSet) *transaction.Transaction {
	var inputs []transaction.TXInput
	var outputs []transaction.TXOutput

	wallets, err := wallet.NewWallets()
	if err != nil {
		log.Panic(err)
	}

	fromWallet := wallets.GetWallet(from)
	pubKeyHash := wallet.HashPubKey(fromWallet.PublicKey)
	accumulated, validOutputs := u.FindSpendableTransactions(pubKeyHash, amount)

	if accumulated < amount {
		log.Panic("ERROR: Not enough funds")
	}

	//Build a list of inputs
	for txid, outs := range validOutputs {
		txID, _ := hex.DecodeString(txid)

		for _, out := range outs {
			inputs = append(inputs, transaction.TXInput{Txid: txID, Vout: out, Signature: nil, Pubkey: fromWallet.PublicKey})
		}
	}

	//Build a list of output
	outputs = append(outputs, *transaction.NewTxOutput(amount, to))
	if accumulated > amount {
		outputs = append(outputs, *transaction.NewTxOutput(accumulated-amount, from)) //Refund
	}

	tx := transaction.Transaction{ID: nil, Vin: inputs, Vout: outputs}
	tx.SetID()
	u.Bc.SignTransaction(&tx, fromWallet.PrivateKey)

	return &tx
}
