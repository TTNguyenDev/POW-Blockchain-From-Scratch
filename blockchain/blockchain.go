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
	"blockchain_from_scratch/utils"
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
	utils.CheckError(err)

	err = db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(blocksBuket))
		tip = b.Get([]byte("l"))

		if tip == nil {
			log.Panic("Latest blockhash doesn't exist")
		}

		return nil
	})
	utils.CheckError(err)

	bc := Blockchain{tip, db}
	return &bc
}

// DB returns the db attribute of Blockchain class
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
	utils.CheckError(err)

	coinbaseTx := transaction.NewCoinbaseTX(benefician, "")
	genesis := NewGenesisBlock(coinbaseTx)

	err = db.Update(func(tx *bolt.Tx) error {

		b, err := tx.CreateBucket([]byte(blocksBuket))
		utils.CheckError(err)
		err = b.Put(genesis.Hash, utils.GobEncode(genesis))

		utils.CheckError(err)
		err = b.Put([]byte("l"), genesis.Hash)

		utils.CheckError(err)
		tip = genesis.Hash
		return nil
	})
	utils.CheckError(err)

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
				UTXO[txID] = outs
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
func (bc *Blockchain) MineBlock(txs []*transaction.Transaction) *Block {
	var lastHash []byte
	var lastHeight int

	err := bc.db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(blocksBuket))
		lastHash = b.Get([]byte("l")) // Get lash block hash

		blockData := b.Get(lastHash)
		block := DeserializeBlock(blockData)
		lastHeight = block.Height

		return nil
	})
	utils.CheckError(err)

	for _, tx := range txs {
		if !bc.VerifyTransaction(tx) {
			log.Panic("Error: Invalid Transaction")
		}
	}
	newBlock := NewBlock(txs, lastHash, lastHeight+1)

	err = bc.db.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(blocksBuket))
		err := b.Put(newBlock.Hash, utils.GobEncode(newBlock))
		utils.CheckError(err)

		err = b.Put([]byte("l"), newBlock.Hash)

		utils.CheckError(err)
		bc.tip = newBlock.Hash

		return nil
	})
	utils.CheckError(err)
	return newBlock
}

// FindTransaction ..
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

// SignTransaction ..
func (bc *Blockchain) SignTransaction(tx *transaction.Transaction, privKey ecdsa.PrivateKey) {
	prevTXs := make(map[string]transaction.Transaction)

	for _, vin := range tx.Vin {
		prevTX, err := bc.FindTransaction(vin.Txid)
		utils.CheckError(err)
		prevTXs[hex.EncodeToString(prevTX.ID)] = prevTX
	}
	tx.Sign(privKey, prevTXs)
}

// VerifyTransaction ..
func (bc *Blockchain) VerifyTransaction(tx *transaction.Transaction) bool {
	if tx.IsCoinBase() {
		return true
	}

	prevTXs := make(map[string]transaction.Transaction)

	for _, vin := range tx.Vin {
		prevTX, err := bc.FindTransaction(vin.Txid)
		utils.CheckError(err)
		prevTXs[hex.EncodeToString(prevTX.ID)] = prevTX
	}
	return tx.Verify(prevTXs)
}

// NewUTXOTransaction -
func NewUTXOTransaction(from, to string, amount int, u *UTXOSet) *transaction.Transaction {
	var inputs []transaction.TXInput
	var outputs []transaction.TXOutput

	wallets, err := wallet.NewWallets()
	utils.CheckError(err)

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
	tx.ID = tx.Hash()
	u.Bc.SignTransaction(&tx, fromWallet.PrivateKey)

	return &tx
}

// AddBlock saves the block into the blockchain
func (bc *Blockchain) AddBlock(block *Block) {
	err := bc.db.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(blocksBuket))
		blockInDb := b.Get(block.Hash)

		if blockInDb != nil {
			return nil
		}

		blockData := utils.GobEncode(block)
		err := b.Put(block.Hash, blockData)
		utils.CheckError(err)

		lastHash := b.Get([]byte("l"))
		lastBlockData := b.Get(lastHash)
		lastBlock := DeserializeBlock(lastBlockData)

		if block.Height > lastBlock.Height {
			err = b.Put([]byte("l"), block.Hash)
			utils.CheckError(err)
			bc.tip = block.Hash
		}
		return nil
	})
	utils.CheckError(err)
}

// View methods

// GetBestHeight ..
func (bc Blockchain) GetBestHeight() int {
	var lastBlock Block

	err := bc.db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(blocksBuket))
		lastHash := b.Get([]byte("l"))
		blockData := b.Get(lastHash)
		lastBlock = *DeserializeBlock(blockData)

		return nil
	})
	if err != nil {
		log.Panic(err)
	}

	return lastBlock.Height
}

// GetBlockHashes ..
func (bc *Blockchain) GetBlockHashes() [][]byte {
	var blocks [][]byte
	bci := bc.Iterator()

	for {
		block := bci.Next()
		blocks = append(blocks, block.Hash)
		if len(block.PrevBlockHash) == 0 {
			break
		}
	}
	return blocks
}

// GetBlock ..
func (bc *Blockchain) GetBlock(blockHash []byte) (Block, error) {
	var block Block

	err := bc.db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(blocksBuket))
		blockData := b.Get(blockHash)
		if blockData == nil {
			return errors.New("Block is not found")
		}

		block = *DeserializeBlock(blockData)
		return nil
	})
	if err != nil {
		return block, err
	}

	return block, nil
}
