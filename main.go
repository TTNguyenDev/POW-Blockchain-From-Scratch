package main

import (
	"fmt"
)

func main() {
	fmt.Println("Creating a new blockchain ...")

	bc := NewBlockchain()

	bc.AddBlock("Send 1 BTC to Alan")
	bc.AddBlock("Send 2 BTC to Ivan")

	for index, b := range bc.blocks {
		fmt.Printf("Prev hash of block %x: %x\n", index, b.PrevBlockHash)
		fmt.Printf("Data:  %s\n", b.Data)
		fmt.Printf("Block Hash: %x\n\n", b.Hash)
	}
}
