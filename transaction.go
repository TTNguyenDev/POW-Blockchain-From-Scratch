package main

import (
	"bytes"
	"crypto/sha256"
	"encoding/gob"
	"fmt"
	"log"
)

const subsidy = 10000

type Transaction struct {
	ID   []byte
	Vin  []TXInput
	Vout []TXOutput
}
type TXOutput struct {
	Value        int
	ScriptPubKey string
}

type TXInput struct {
	Txid      []byte
	Vout      int
	ScriptSig string
}

func (tx *Transaction) SetID() {
	var encoded bytes.Buffer
	var hash [32]byte

	enc := gob.NewEncoder(&encoded)
	err := enc.Encode(tx)

	if err != nil {
		log.Panic(err)
	}

	hash = sha256.Sum256(encoded.Bytes())
	tx.ID = hash[:]
}
func NewCoinbaseTX(to, data string) *Transaction {
	if len(data) == 0 {
		data = fmt.Sprintf("Reward to %s", to)
	}

	txin := TXInput{
		[]byte{},
		-1, /*index of output for coinbase transaction*/
		data,
	}

	txout := TXOutput{
		subsidy,
		to,
	}

	tx := Transaction{nil, []TXInput{txin}, []TXOutput{txout}}
	tx.SetID()

	return &tx
}
