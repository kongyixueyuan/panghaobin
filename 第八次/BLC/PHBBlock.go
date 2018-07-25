package BLC

import (
	"bytes"
	"encoding/gob"
	"fmt"
	"log"
	"time"
)

type PHBBlock struct {
	PHBHeight int64
	PHBPrevBlockHash []byte
	PHBTxs []*PHBTransaction
	PHBTimestamp int64
	PHBHash []byte
	PHBNonce int64
}

func (block *PHBBlock) PHBHashTransactions() []byte {
	var transactions [][]byte
	for _, tx := range block.PHBTxs {
		transactions = append(transactions, tx.PHBSerialize())
	}
	mTree := PHBNewMerkleTree(transactions)
	return mTree.PHBRootNode.PHBData
}

func (block *PHBBlock) PHBSerialize() []byte {
	var result bytes.Buffer
	encoder := gob.NewEncoder(&result)
	err := encoder.Encode(block)
	if err != nil {
		log.Panic(err)
	}
	return result.Bytes()
}

func PHBDeserializeBlock(blockBytes []byte) *PHBBlock {
	var block PHBBlock
	decoder := gob.NewDecoder(bytes.NewReader(blockBytes))
	err := decoder.Decode(&block)
	if err != nil {
		log.Panic(err)
	}
	return &block
}

func PHBNewBlock(txs []*PHBTransaction, height int64, prevBlockHash []byte) *PHBBlock {
	block := &PHBBlock{height, prevBlockHash, txs, time.Now().Unix(), nil, 0}
	pow := NewProofOfWork(block)
	hash, nonce := pow.PHBRun()
	block.PHBHash = hash[:]
	block.PHBNonce = nonce
	fmt.Println()
	return block
}

func PHBCreateGenesisBlock(txs []*PHBTransaction) *PHBBlock {
	return PHBNewBlock(txs, 1, []byte{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0})
}
