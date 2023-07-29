package cli

import (
	"flag"
	"fmt"
	"log"
	"os"
)

// CLI struct
type CLI struct{}

func (cli *CLI) printUsage() {
	fmt.Println("Usage:")
	fmt.Println(" createwallet - create new wallet")
	fmt.Println(" newblockchain -address ADDRESS - new blockchain with benefician's address")
	fmt.Println(" sendcoins -from ADDRESS -to ADDRESS -amount AMOUNT - send coin function")
	fmt.Println(" getbalance -address ADDRESS - sum of UTXOs of the given address")
	fmt.Println(" printchain - print all the blocks of the blockchain")
	fmt.Println(" listaddresses - print all the addressses in wallet.dat file")
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

	createWalletCmd := flag.NewFlagSet("createwallet", flag.ExitOnError)
	newBlockchainCmd := flag.NewFlagSet("newblockchain", flag.ExitOnError)
	sendCoinsCmd := flag.NewFlagSet("sendcoins", flag.ExitOnError)
	getBalanceCmd := flag.NewFlagSet("getbalance", flag.ExitOnError)
	printChainCmd := flag.NewFlagSet("printchain", flag.ExitOnError)
	listAddressesCmd := flag.NewFlagSet("listaddrs", flag.ExitOnError)

	newBlockchainData := newBlockchainCmd.String("address", "", "Address")
	getBalanceData := getBalanceCmd.String("address", "", "Address")
	fromAddressData := sendCoinsCmd.String("from", "", "From address")
	toAddressData := sendCoinsCmd.String("to", "", "To address")
	amountData := sendCoinsCmd.Int("amount", 0, "Amount")

	switch os.Args[1] {
	case "createwallet":
		err := createWalletCmd.Parse(os.Args[2:])
		if err != nil {
			log.Panic(err)
		}
		cli.createWallet()
	case "newblockchain":
		err := newBlockchainCmd.Parse(os.Args[2:])
		if err != nil {
			log.Panic(err)
		}

		if *newBlockchainData == "" {
			newBlockchainCmd.Usage()
			os.Exit(1)
		}
		cli.createBlockchain(*newBlockchainData)
	case "sendcoins":
		err := sendCoinsCmd.Parse(os.Args[2:])
		if err != nil {
			log.Panic(err)
		}

		if *fromAddressData == "" || *toAddressData == "" || *amountData == 0 {
			sendCoinsCmd.Usage()
			os.Exit(1)
		}
		cli.sendCoins(*fromAddressData, *toAddressData, *amountData)
	case "getbalance":
		err := getBalanceCmd.Parse(os.Args[2:])
		if err != nil {
			log.Panic(err)
		}

		if *getBalanceData == "" {
			getBalanceCmd.Usage()
			os.Exit(1)
		}
		cli.getBalance(*getBalanceData)
	case "printchain":
		err := printChainCmd.Parse(os.Args[2:])
		if err != nil {
			log.Panic(err)
		}
		cli.printChain()
	case "listaddrs":
		err := listAddressesCmd.Parse(os.Args[2:])
		if err != nil {
			log.Panic(err)
		}
		cli.listaddresses()
	default:
		cli.printUsage()
		os.Exit(1)
	}
}
