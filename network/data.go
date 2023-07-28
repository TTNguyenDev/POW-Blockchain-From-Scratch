package network

import (
	"bytes"
	"encoding/gob"
	"encoding/hex"
	"log"

	"blockchain_from_scratch/blockchain"
	"blockchain_from_scratch/utils"
)

type getdata struct {
	AddrFrom string
	Type     string
	ID       []byte
}

func handleGetData(request []byte, bc *blockchain.Blockchain) {
	var buff bytes.Buffer
	var payload getdata

	buff.Write(request[commandLength:])
	dec := gob.NewDecoder(&buff)
	err := dec.Decode(&payload)
	log.Panic(err)

	if payload.Type == "block" {
		block, err := bc.GetBlock([]byte(payload.ID))
		log.Panic(err)
		sendBlock(payload.AddrFrom, &block)
	}

	if payload.Type == "tx" {
		txID := hex.EncodeToString(payload.ID)
		tx := mempool[txID]

		sendTx(payload.AddrFrom, &tx)
	}
}

func sendGetData(address, kind string, id []byte) {
	payload := utils.GobEncode(getdata{address, kind, id})
	req := append(commandToBytes("getdata"), payload...)

	sendData(address, req)
}
