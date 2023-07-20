package main

import (
	"bytes"
	"crypto/sha256"
	"strconv"
	"time"
)

type Block struct {
	TimeStamp     int64
	Data          []byte
	PrevBlockHash []byte
	Hash          []byte
}

func (b *Block) SetHash() {
	timeStamp := []byte(strconv.FormatInt(b.TimeStamp, 10))
	header := bytes.Join([][]byte{b.PrevBlockHash, b.Data, timeStamp}, []byte{})
	hash := sha256.Sum256(header)

	b.Hash = hash[:]
}

func NewBlock(data string, prevBlockHash []byte) *Block {
	block := &Block{time.Now().Unix(), []byte(data), prevBlockHash, []byte{}}
	block.SetHash()
	return block
}
