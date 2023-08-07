package network

import (
	"bytes"
	"encoding/gob"
	"fmt"

	"blockchain_from_scratch/blockchain"
	"blockchain_from_scratch/utils"
)

type block struct {
	AddrFrom string
	Block    []byte
}

//TODO: Thay vì tin tưởng block gửi tới vô điều kiện, chúng ta phải xác minh block mới nhận được có đúng hay không trước khi cho vào blockchain.
//TODO: Thay vì chạy UTXOSet.Reindex(), hãy chạy UTXOSet.Update(block) đối với mỗi block nhận được để giảm thiểu việc quét qua cả blockchain lãng phí.

// handleBlock ..
func handleBlock(request []byte, bc *blockchain.Blockchain) {
	var buff bytes.Buffer
	var payload block

	buff.Write(request[commandLength:])
	dec := gob.NewDecoder(&buff)
	err := dec.Decode(&payload)
	utils.CheckError(err)

	blockData := payload.Block
	block := blockchain.DeserializeBlock(blockData)

	fmt.Println("Received a new block!")
	bc.AddBlock(block)

	fmt.Printf("Added block %x\n", block.Hash)
	if len(blocksInTransit) > 0 {
		blockHash := blocksInTransit[0]
		sendGetData(payload.AddrFrom, "block", blockHash)
		blocksInTransit = blocksInTransit[1:]
	} else {
		UTXOSet := blockchain.UTXOSet{Bc: bc}
		UTXOSet.Reindex()
	}
}

// sendBlock ..
func sendBlock(addr string, b *blockchain.Block) {
	data := block{addr, utils.GobEncode(b)}
	payload := utils.GobEncode(data)
	req := append(commandToBytes("block"), payload...)

	sendData(addr, req)
}
