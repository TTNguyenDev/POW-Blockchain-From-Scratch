package cli

import (
	"encoding/hex"
	"flag"
	"fmt"
	"log"
	"os"
	"strconv"

	"blockchain_from_scratch/blockchain"
	"blockchain_from_scratch/blockchain/transaction"
	"blockchain_from_scratch/utils"
	"blockchain_from_scratch/wallet"
)

// CLI struct
type CLI struct{}

func (cli *CLI) createBlockchain(benefician string) {
	bc := blockchain.NewBlockchain(benefician)
	defer bc.DB().Close()
	u := blockchain.UTXOSet{Bc: bc}
	u.Reindex()
	fmt.Println("Done!")
}

func (cli *CLI) createWallet() {
	wallets, _ := wallet.NewWallets()
	address := wallets.CreateWallet()

	fmt.Printf("Your new address: %s\n", address)
}

func (cli *CLI) printUsage() {
	fmt.Println("Usage:")
	fmt.Println(" createwallet - create new wallet")
	fmt.Println(" newblockchain -address ADDRESS - new blockchain with benefician's address")
	fmt.Println(" sendcoins -from ADDRESS -to ADDRESS -amount AMOUNT - send coin function")
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

	createWalletCmd := flag.NewFlagSet("createwallet", flag.ExitOnError)

	newBlockchainCmd := flag.NewFlagSet("newblockchain", flag.ExitOnError)
	getBalanceCmd := flag.NewFlagSet("getbalance", flag.ExitOnError)
	sendCoinsCmd := flag.NewFlagSet("sendcoins", flag.ExitOnError)
	printChainCmd := flag.NewFlagSet("printchain", flag.ExitOnError)

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
	case "printchain":
		err := printChainCmd.Parse(os.Args[2:])
		if err != nil {
			log.Panic(err)
		}
		cli.printChain()
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
	default:
		cli.printUsage()
		os.Exit(1)
	}
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

	fmt.Printf("Balance of %s: %d\n", address, balance)
}

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
		pow := blockchain.NewProofOfWork(b)
		fmt.Printf("IsValid: %s \n\n", strconv.FormatBool(pow.Validate()))
	}
}
