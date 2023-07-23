package transaction

import (
	"bytes"

	"blockchain_from_scratch/utils"
)

// TXOutput - Definition
type TXOutput struct {
	Value      int
	PubKeyHash []byte
}

func (out *TXOutput) Lock(address []byte) {
	pubHash := utils.Base58Decode(address)
	pubHash = pubHash[1 : len(pubHash)-4] // remove version & checksum
	out.PubKeyHash = pubHash
}

func (out *TXOutput) IsLockedWithKey(pubHash []byte) bool {
	return bytes.Equal(out.PubKeyHash, pubHash)
}

func NewTxOutput(value int, address string) *TXOutput {
	txo := &TXOutput{value, nil}
	txo.Lock([]byte(address))
	return txo
}
