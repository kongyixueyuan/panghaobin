package BLC

import "time"

type Block struct {
	Timestamp int64 //时间戳

	Height int64 //区块高度
	PrevBlockHash []byte //上一个区块HASH
	Hash []byte //本区块的hash
	Data []byte //交易数据
	Nonce int //计算PoW
}

func NewBlock(data string, height int64, prevBlockHash []byte) *Block  {
	block := &Block{
		Timestamp: time.Now().Unix(),
		Height: height,
		PrevBlockHash: prevBlockHash,
		Hash: []byte{},
		Data: []byte(data),
		Nonce: 0}
	pow := NewProofWork(block)
	nonce, hash := pow.Run()
	block.Hash = hash[:]
	block.Nonce = nonce

	return block
}

func NewGenesisBlock() *Block {
	return NewBlock("Genesis Block", 1, []byte{})
}


