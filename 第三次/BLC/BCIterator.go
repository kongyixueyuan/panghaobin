package BLC

import (
	"github.com/boltdb/bolt"
	"log"
	"fmt"
	"time"
	"math/big"
)

type BlockChainIterator struct {
	currentHash []byte
	db *bolt.DB
}


func (blockChain *BlockChain)Iterator() *BlockChainIterator  {
	return &BlockChainIterator{blockChain.tip, blockChain.db}
}

func (iter *BlockChainIterator) Next() *Block {
	var block *Block
	err := iter.db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(blockBucket))
		currentBlockData := b.Get(iter.currentHash)
		block = DeserializeBlock(currentBlockData)//反序列化获取区块
		iter.currentHash = block.PrevBlockHash //更新迭代器的上一个区块hash
		return nil
	})
	if err != nil {
		log.Panic(err)
	}
	return block

}

func (bc *BlockChain) VisitBlockChain()  {
	iter := bc.Iterator()
	for {
		block := iter.Next()

		fmt.Printf("区块高度：%d\n", block.Height)
		fmt.Printf("上一个区块的hash：%x\n", block.PrevBlockHash)
		fmt.Printf("区块数据：%s\n", block.Data)
		fmt.Printf("区块时间戳：%s\n", time.Unix(block.Timestamp, 0).Format("2006-01-02 15:04:05 AM") )//go诞生之日
		fmt.Printf("区块hash：%x\n", block.Hash)
		fmt.Printf("Nonce: %d\n\n", block.Nonce)
		var hashInt big.Int
		hashInt.SetBytes(block.PrevBlockHash)
		if big.NewInt(0).Cmp(&hashInt) == 0 {
			break
		}

	}
}