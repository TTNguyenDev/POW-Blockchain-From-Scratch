package cli

import (
	"encoding/hex"
	"fmt"
	"strconv"

	"blockchain_from_scratch/blockchain"
	"blockchain_from_scratch/blockchain/consensus"
)

// createBlockchain ..
func (cli *CLI) createBlockchain(benefician string) {
	bc := blockchain.NewBlockchain(benefician)
	defer bc.DB().Close()
	u := blockchain.UTXOSet{Bc: bc}
	u.Reindex()
	fmt.Println("Done!")
}

// printChain ..
func (cli *CLI) printChain() {
	bc := blockchain.BCInstance()
	defer bc.DB().Close()
	bci := bc.Iterator()

	fmt.Printf("Querying blockchain data:\n")
	for {
		b := bci.Next()
		if b == nil {
			fmt.Println("Ended")
			break
		}
		fmt.Printf("Prev hash of block: %x\n", b.PrevBlockHash)
		for _, tx := range b.Transactions {
			fmt.Printf("Transaction: %s\n", hex.EncodeToString(tx.ID))
		}
		fmt.Printf("Block Hash: %x\n", b.Hash)
		pow := consensus.NewProofOfWork(b)
		fmt.Printf("IsValid: %s \n\n", strconv.FormatBool(pow.Validate()))
	}
}
