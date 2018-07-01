package BLC

type BlockChain struct {
	blocks []*Block
}

func NewBlockChain() *BlockChain {
	return &BlockChain{[]*Block{NewGenesisBlock()}}
}

func (bc *BlockChain) AddBlock(data string)  {
	prevBlock := bc.blocks[len(bc.blocks)-1]
	newBlock := NewBlock(data, bc.blocks[len(bc.blocks)-1].Height+1, prevBlock.Hash)
	bc.blocks = append(bc.blocks, newBlock)
}