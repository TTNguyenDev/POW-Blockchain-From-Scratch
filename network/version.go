package network

import (
	"bytes"
	"encoding/gob"

	"blockchain_from_scratch/blockchain"
	"blockchain_from_scratch/utils"
)

type version struct {
	Version    int
	BestHeight int
	AddrFrom   string
}

func sendVersion(addr string, bc *blockchain.Blockchain) {
	bestHeight := bc.GetBestHeight()
	payload := utils.GobEncode(version{nodeVersion, bestHeight, addr})

	request := append(commandToBytes("version"), payload...)
	sendData(addr, request)
}

func handleVersion(request []byte, bc *blockchain.Blockchain) {
	var buff bytes.Buffer
	var payload version

	buff.Write(request[commandLength:])
	dec := gob.NewDecoder(&buff)
	err := dec.Decode(&payload)
	utils.CheckError(err)

	bestHeight := bc.GetBestHeight()
	foreignerBestHeight := payload.BestHeight

	if bestHeight < foreignerBestHeight {
		sendGetBlocks(payload.AddrFrom)
	} else {
		sendVersion(payload.AddrFrom, bc)
	}

	if !nodeIsKnown(payload.AddrFrom) {
		knownNodes = append(knownNodes, payload.AddrFrom)
	}
}
