package main

type Blockchain struct {
	blocks []*Block
}

func NewBlockchain() *Blockchain {
	return &Blockchain{
		[]*Block{
			NewGenesisBlock(),
		},
	}
}

func NewGenesisBlock() *Block {
	return NewBlock("Genesis Block", []byte{})
}

func (bc *Blockchain) AddBlock(data string) {
	len := len(bc.blocks) - 1
	prevBlock := bc.blocks[len]
	newBlock := NewBlock(data, prevBlock.Hash)
	bc.blocks = append(bc.blocks, newBlock)
}
