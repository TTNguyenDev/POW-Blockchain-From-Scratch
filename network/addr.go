package network

import (
	"bytes"
	"encoding/gob"
	"fmt"

	"blockchain_from_scratch/utils"
)

type addr struct {
	AddrList []string
}

func handleAddr(req []byte) {
	var buff bytes.Buffer
	var payload addr

	buff.Write(req[commandLength:])
	dec := gob.NewDecoder(&buff)
	err := dec.Decode(&payload)
	utils.CheckError(err)

	knownNodes = append(knownNodes, payload.AddrList...)
	fmt.Printf("There are %d known nodes now!\n", len(knownNodes))
	requestBlocks()
}
