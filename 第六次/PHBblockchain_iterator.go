package main

import (
	"github.com/boltdb/bolt"
	"log"
)

type PHBBlockchainIterator struct {
	phbcurrentHash []byte
	phbdb          *bolt.DB
}

// Next returns next block starting from the tip
func (i *PHBBlockchainIterator) PHBNext() *PHBBlock {
	var block *PHBBlock
	err := i.phbdb.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(blocksBucket))
		encodedBlock := b.Get(i.phbcurrentHash)
		block = PHBDeserializeBlock(encodedBlock)

		return nil
	})
	if err != nil {
		log.Panic(err)
	}
	i.phbcurrentHash = block.PHBPrevBlockHash

	return block
}
