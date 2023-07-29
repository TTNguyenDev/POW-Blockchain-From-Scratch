package transaction

import (
	"bytes"
	"encoding/gob"
	"log"

	"blockchain_from_scratch/utils"
)

// TXOutput - Definition
type TXOutput struct {
	Value      int
	PubKeyHash []byte
}

// Lock method sets the pubHash, which is decoded from the input address, to the out.PubKeyHash.
func (out *TXOutput) Lock(address []byte) {
	pubHash := utils.Base58Decode(address)
	pubHash = pubHash[1 : len(pubHash)-4] // remove version & checksum
	out.PubKeyHash = pubHash
}

// IsLockedWithKey ..
func (out *TXOutput) IsLockedWithKey(pubHash []byte) bool {
	return bytes.Equal(out.PubKeyHash, pubHash)
}

// NewTxOutput ..
func NewTxOutput(value int, address string) *TXOutput {
	txo := &TXOutput{value, nil}
	txo.Lock([]byte(address))
	return txo
}

// DeSerializeTXOuputs ..
func DeSerializeTXOuputs(outputs []byte) []TXOutput {

	var result []TXOutput

	decoder := gob.NewDecoder(bytes.NewReader(outputs))
	err := decoder.Decode(&result)

	if err != nil {
		log.Panic(err)
	}

	return result
}
