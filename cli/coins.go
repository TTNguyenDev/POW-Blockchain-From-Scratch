package cli

import (
	"fmt"

	"blockchain_from_scratch/blockchain"
	"blockchain_from_scratch/blockchain/transaction"
	"blockchain_from_scratch/utils"
)

// getBalance
func (cli *CLI) getBalance(address string) {
	bc := blockchain.BCInstance()
	u := blockchain.UTXOSet{Bc: bc}
	defer bc.DB().Close()

	balance := 0
	pubKeyHash := utils.Base58Decode([]byte(address))
	pubKeyHash = pubKeyHash[1 : len(pubKeyHash)-4]
	UTXOs := u.FindUTXO(pubKeyHash)

	for _, out := range UTXOs {
		balance += out.Value
	}

	fmt.Printf("%d", balance)
}

// sendCoins ...
func (cli *CLI) sendCoins(from, to string, amount int) {
	bc := blockchain.BCInstance()
	u := blockchain.UTXOSet{Bc: bc}
	defer bc.DB().Close()

	//Build Input for this transaction
	tx := blockchain.NewUTXOTransaction(from, to, amount, &u)
	cbTx := transaction.NewCoinbaseTX(from, "")
	bc.MineBlock([]*transaction.Transaction{cbTx, tx})

	fmt.Println("Success!")
}
