package cli

import (
	"fmt"
	"log"

	"blockchain_from_scratch/wallet"
)

// createWallet ..
func (cli *CLI) createWallet() {
	wallets, _ := wallet.NewWallets()
	address := wallets.CreateWallet()

	fmt.Printf("Your new address: %s\n", address)
}

// listaddresses ..
func (cli *CLI) listaddresses() {
	wallets, err := wallet.NewWallets()
	if err != nil {
		log.Panic(err)
	}
	addresses := wallets.GetAddresses()

	for _, address := range addresses {
		fmt.Println(address)
	}
}
