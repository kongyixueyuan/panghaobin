package BLC

import (
	"github.com/boltdb/bolt"
	"log"
)

type PHBBlockchainIterator struct {
	PHBCurrentHash []byte
	PHBDB          *bolt.DB
}

func (blockchainIterator *PHBBlockchainIterator) PHBNext() *PHBBlock {
	var block *PHBBlock
	err := blockchainIterator.PHBDB.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(blockTableName))
		if b != nil {
			currentBloclBytes := b.Get(blockchainIterator.PHBCurrentHash)
			//获取当前迭代器里的currentHash对应的区块
			block = PHBDeserializeBlock(currentBloclBytes)
			//更新迭代器里面CurrentHash
			blockchainIterator.PHBCurrentHash = block.PHBPrevBlockHash
		}
		return nil
	})
	if err != nil {
		log.Panic(err)
	}
	return block

}
