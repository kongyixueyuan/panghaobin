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
		fmt.Printf("========================================\n")
		fmt.Printf("区块高度：%d\n", block.Height)
		fmt.Printf("上一个区块的hash：%x\n", block.PrevBlockHash)
		fmt.Printf("区块交易数据\n")
		for _, tx := range block.Txs {
			fmt.Printf("交易hash：%x\n", tx.TxHash)
			fmt.Println("  交易输入:")
			for _, in := range tx.Vin {
				fmt.Printf("%x\n", in.TxHash)
				fmt.Printf("%d\n", in.Vout)
				fmt.Printf("用户签名：%s\n", in.Signature)
			}

			fmt.Println("  交易输出:")
			for _, out := range tx.Vout {
				fmt.Println(out.Value)
				fmt.Println(out.PubKey)
			}
		}
		fmt.Println("------------------------------\n")
		fmt.Printf("区块时间戳：%s\n", time.Unix(block.Timestamp, 0).Format("2006-01-02 15:04:05 AM") )//go诞生之日
		fmt.Printf("区块hash：%x\n", block.Hash)
		fmt.Printf("Nonce: %d\n", block.Nonce)
		fmt.Printf("========================================\n\n")
		var hashInt big.Int
		hashInt.SetBytes(block.PrevBlockHash)
		if big.NewInt(0).Cmp(&hashInt) == 0 {
			break
		}

	}
}