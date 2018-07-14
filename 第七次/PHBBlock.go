package main

import (
	"time"
	"bytes"
	"encoding/gob"
	"log"
)

type PHBBlock struct {
	PHBTimestamp int64
	PHBTransactions []*PHBTransaction
	PHBPrevBlockHash []byte
	PHBHash []byte
	PHBNonce int
	PHBHeight int
}

// NewBlock creates and returns Block
func PHBNewBlock(transactions []*PHBTransaction, prevBlockHash []byte, height int) *PHBBlock {
	block := &PHBBlock{time.Now().Unix(), transactions, prevBlockHash, []byte{}, 0, height}
	pow := PHBNewProofOfWork(block)
	nonce, hash := pow.PHBRun()

	block.PHBHash = hash[:]
	block.PHBNonce = nonce

	return block
}


func PHBNewGenesisBlock(coinbase *PHBTransaction) *PHBBlock {
	return PHBNewBlock([]*PHBTransaction{coinbase}, []byte{}, 0)
}


func (b *PHBBlock) PHBHashTransactions() []byte {
	var transactions [][]byte

	for _, tx := range b.PHBTransactions {
		transactions = append(transactions, tx.PHBSerialize())
	}
	mTree := PHBNewMerkleTree(transactions)

	return mTree.PHBRootNode.PHBData
}


func (b *PHBBlock) PHBSerialize() []byte {
	var result bytes.Buffer
	encoder := gob.NewEncoder(&result)

	err := encoder.Encode(b)
	if err != nil {
		log.Panic(err)
	}

	return result.Bytes()
}


func PHBDeserializeBlock(d []byte) *PHBBlock {
	var block PHBBlock

	decoder := gob.NewDecoder(bytes.NewReader(d))
	err := decoder.Decode(&block)
	if err != nil {
		log.Panic(err)
	}

	return &block
}
