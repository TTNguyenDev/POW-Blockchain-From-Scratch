package main

import (
	"encoding/hex"
	"flag"
	"fmt"
	"log"
	"os"
	"strconv"
)

// CLI struct
type CLI struct{}

func (cli *CLI) createBlockchain(benefician string) {
	bc := NewBlockchain(benefician)
	bc.db.Close()
	fmt.Println("Done!")

}

func (cli *CLI) printUsage() {
	fmt.Println("Usage:")
	fmt.Println(" newblockchain -address ADDRESS - new blockchain with benefician's address")
	fmt.Println(" getbalance -address ADDRESS - sum of UTXOs of the given address")
	fmt.Println(" printchain - print all the blocks of the blockchain")
}

func (cli *CLI) validateArgs() {
	if len(os.Args) < 2 {
		cli.printUsage()
		os.Exit(1)
	}
}

// Run ...
func (cli *CLI) Run() {
	cli.validateArgs()

	newBlockchainCmd := flag.NewFlagSet("newblockchain", flag.ExitOnError)
	getBalanceCmd := flag.NewFlagSet("getbalance", flag.ExitOnError)
	printChainCmd := flag.NewFlagSet("printchain", flag.ExitOnError)

	newBlockchainData := newBlockchainCmd.String("address", "", "Address")
	getBalanceData := getBalanceCmd.String("address", "", "Address")

	switch os.Args[1] {
	case "newblockchain":
		err := newBlockchainCmd.Parse(os.Args[2:])
		if err != nil {
			log.Panic(err)
		}
	case "printchain":
		err := printChainCmd.Parse(os.Args[2:])
		if err != nil {
			log.Panic(err)
		}
	case "getbalance":
		err := getBalanceCmd.Parse(os.Args[2:])
		if err != nil {
			log.Panic(err)
		}
	default:
		cli.printUsage()
		os.Exit(1)
	}

	if newBlockchainCmd.Parsed() {
		if *newBlockchainData == "" {
			newBlockchainCmd.Usage()
			os.Exit(1)
		}
		cli.createBlockchain(*newBlockchainData)
	}

	if getBalanceCmd.Parsed() {
		if *getBalanceData == "" {
			getBalanceCmd.Usage()
			os.Exit(1)
		}
		cli.getBalance(*getBalanceData)
	}

	if printChainCmd.Parsed() {
		cli.printChain()
	}
}

// getBalance
func (cli *CLI) getBalance(address string) {
	bc := BCInstance()
	defer bc.db.Close()

	balance := 0
	UTXOs := bc.FindUTXO(address)

	for _, out := range UTXOs {
		balance += out.Value
	}

	fmt.Printf("Balance of %s: %d\n", address, balance)
}

func (cli *CLI) printChain() {
	bc := BCInstance()
	defer bc.db.Close()
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
		pow := NewProofOfWork(b)
		fmt.Printf("IsValid: %s \n\n", strconv.FormatBool(pow.Validate()))
	}
}
