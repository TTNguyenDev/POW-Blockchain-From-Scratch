package cli

import (
	"fmt"
	"log"

	"blockchain_from_scratch/network"
	"blockchain_from_scratch/wallet"
)

func (cli *CLI) startNode(nodeID, minerAddress string) {
	fmt.Printf("Staring node %s\n", nodeID)
	if len(minerAddress) > 0 {
		if wallet.ValidateAddress(minerAddress) {
			fmt.Println("Mining is on. Address to receive rewards: ", minerAddress)
		} else {
			log.Panic("Wrong miner address")
		}
	}
	network.StartServer(nodeID, minerAddress)
}
