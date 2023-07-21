package main

import (
	"fmt"
)

func main() {
	fmt.Println("Creating a new blockchain ...")

	bc := NewBlockchain()

	bc.AddBlock("Send 1 BTC to Alan")
	bc.AddBlock("Send 2 BTC to Ivan")

	bci := bc.Iterator()

	var b *Block
	b = bci.Next()
	for {
		fmt.Printf("Prev hash of block: %x\n", b.PrevBlockHash)
		fmt.Printf("Data:  %s\n", b.Data)
		fmt.Printf("Block Hash: %x\n", b.Hash)
		b = bci.Next()
		if b == nil {
			break
		}
	}
	// 	while (b != nil) {
	//
	// }
	// for index, b := range bc.blocks {
	// 	// Validate my chain
	// 	pow := NewProofOfWork(b)
	// 	fmt.Printf("Pow validation: %s\n", strconv.FormatBool(pow.Validate()))
	// 	fmt.Println()
	// }
}
