package network

import (
	"fmt"
	"io/ioutil"
	"log"
	"net"

	"blockchain_from_scratch/blockchain"
	"blockchain_from_scratch/blockchain/transaction"
)

const nodeVersion = 1
const commandLength = 12

// TODO: We should replace these lines of code with the p2p find peers function
var protocol = "tcp"
var miningAddress string
var nodeAddress string
var knownNodes = []string{"localhost:3000"}
var blocksInTransit = [][]byte{}
var mempool = make(map[string]transaction.Transaction)

func StartServer(nodeID, minerAddress string) {
	nodeAddress = fmt.Sprintf("localhost:%s", nodeID)
	miningAddress := minerAddress
	listner, err := net.Listen(protocol, nodeAddress)
	if err != nil {
		log.Panic(err)
	}
	defer listner.Close()

	bc := blockchain.NewBlockchain(miningAddress) //TODO: Start with nodeID

	if nodeAddress != knownNodes[0] {
		sendVersion(knownNodes[0], bc)
	}

	for {
		conn, err := listner.Accept()
		if err != nil {
			log.Panic(err)
		}
		go handleConnection(conn, bc)
	}
}

func handleConnection(conn net.Conn, bc *blockchain.Blockchain) {
	request, err := ioutil.ReadAll(conn)
	log.Panic(err)
	command := bytesToCommand(request[:commandLength])
	fmt.Printf("Received %s command\n", command)

	switch command {
	case "version":

	default:
		fmt.Println("Unknown command!")
	}
	conn.Close()
}

// View methods
