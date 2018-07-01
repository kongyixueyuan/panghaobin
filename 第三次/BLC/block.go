package BLC

import (
	"time"
	"bytes"
	"encoding/gob"
	"log"
)

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
	return NewBlock("Genesis Block", 0, []byte{})
}

//区块序列化
func (b *Block) Serialize() []byte  {
	var result bytes.Buffer
	encoder := gob.NewEncoder(&result)
	err := encoder.Encode(b)
	if err != nil {
		log.Panic(err)
	}

	return result.Bytes()
}

//反序列化
func DeserializeBlock(d []byte) *Block  {
	var block Block

	decoder := gob.NewDecoder(bytes.NewBuffer(d))
	err := decoder.Decode(&block)
	if err != nil {
		log.Panic(err)
	}
	return &block
}

