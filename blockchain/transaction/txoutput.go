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

func (out TXOutput) Serialize() []byte {
	var result bytes.Buffer
	encoder := gob.NewEncoder(&result)

	err := encoder.Encode(out)
	if err != nil {
		log.Panic(err)
	}

	return result.Bytes()
}

func SerializeTXOutputs(outputs []TXOutput) []byte {
	var result bytes.Buffer
	encoder := gob.NewEncoder(&result)

	err := encoder.Encode(outputs)
	if err != nil {
		log.Panic(err)
	}

	return result.Bytes()
}

func DeSerializeTXOuputs(outputs []byte) []TXOutput {

	var result []TXOutput

	decoder := gob.NewDecoder(bytes.NewReader(outputs))
	err := decoder.Decode(&result)

	if err != nil {
		log.Panic(err)
	}

	return result
}
