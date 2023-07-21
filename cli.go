package main

import (
	"flag"
	"fmt"
	"os"
	"strconv"
)

type CLI struct{}

func (cli *CLI) createBlockchain(benefician string) {
	bc := NewBlockchain(benefician)
	bc.db.Close()
	fmt.Println("Done!")

}

func (cli *CLI) printUsage() {
	fmt.Println("Usage:")
	fmt.Println(" newblockchain -address ADDRESS - new blockchain with benefician's address")
	fmt.Println(" addblock -data BLOCK_DATA - add a block to the blockchain")
	fmt.Println(" printchain - print all the blocks of the blockchain")
}

func (cli *CLI) validateArgs() {
	if len(os.Args) < 2 {
		cli.printUsage()
		os.Exit(1)
	}
}
func (cli *CLI) Run() {
	cli.validateArgs()

	newBlockchainCmd := flag.NewFlagSet("newblockchain", flag.ExitOnError)
	addBlockCmd := flag.NewFlagSet("addblock", flag.ExitOnError)
	printChainCmd := flag.NewFlagSet("printchain", flag.ExitOnError)

	newBlockchainData := newBlockchainCmd.String("address", "", "Address")
	addBlockData := addBlockCmd.String("data", "", "Block data")

	switch os.Args[1] {
	case "newblockchain":
		newBlockchainCmd.Parse(os.Args[2:])
	case "addblock":
		addBlockCmd.Parse(os.Args[2:])
	case "printchain":
		printChainCmd.Parse(os.Args[2:])
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

	if addBlockCmd.Parsed() {
		if *addBlockData == "" {
			addBlockCmd.Usage()
			os.Exit(1)
		}
		cli.addBlock(*addBlockData)
	}

	if printChainCmd.Parsed() {
		cli.printChain()
	}
}

func (cli *CLI) addBlock(data string) {
	// bc := BCInstance()
	// defer bc.db.Close()
	// bc.AddBlock(data)
	fmt.Println("A new block is Added")
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
		fmt.Printf("Data:  %s\n", b.Data)
		fmt.Printf("Block Hash: %x\n", b.Hash)
		pow := NewProofOfWork(b)
		fmt.Printf("IsValid: %s \n\n", strconv.FormatBool(pow.Validate()))
	}
}
