package transaction

import (
	"bytes"

	"blockchain_from_scratch/utils"
)

// TXInput - Definition
type TXInput struct {
	Txid      []byte
	Vout      int
	Signature []byte
	Pubkey    []byte
}

func (in *TXInput) UsesKey(pubHash []byte) bool {
	lockingHash := utils.HashPubKey(in.Pubkey)
	return bytes.Equal(lockingHash, pubHash)
}
