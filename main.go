package main

import (
	"fmt"
	"strconv"
)

func main() {
	fmt.Println("Creating a new blockchain ...")

	bc := NewBlockchain()

	bc.AddBlock("Send 1 BTC to Alan")
	bc.AddBlock("Send 2 BTC to Ivan")

	for index, b := range bc.blocks {
		fmt.Printf("Prev hash of block %x: %x\n", index, b.PrevBlockHash)
		fmt.Printf("Data:  %s\n", b.Data)
		fmt.Printf("Block Hash: %x\n", b.Hash)

		// Validate my chain
		pow := NewProofOfWork(b)
		fmt.Printf("Pow validation: %s\n", strconv.FormatBool(pow.Validate()))
		fmt.Println()
	}
}
