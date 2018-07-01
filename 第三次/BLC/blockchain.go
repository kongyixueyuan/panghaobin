package BLC

import (
	"github.com/boltdb/bolt"
	"log"
	"os"
)

const dbPath = "blockchain.db"
const blockBucket = "blocks"
type BlockChain struct {
	tip []byte
	db *bolt.DB
}

func CreatBlockChain() *BlockChain  {
	db, err := bolt.Open(dbPath, 0600, nil)
	if err != nil {
		log.Fatal(err)//exit
	}
	var tip []byte
	err = db.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(blockBucket))//获取表
		if err != nil {
			log.Panic(err)
		}
		if b == nil {
			b, err = tx.CreateBucket([]byte(blockBucket))
			if err != nil {
				log.Panic(err)
			}
		}
		if b != nil { //表存在
			genesisBlock := NewGenesisBlock()
			//将创世区块存储到表
			err := b.Put(genesisBlock.Hash, genesisBlock.Serialize())
			if err != nil {
				log.Panic(err)
			}

			//存储最新区块的hash
			err = b.Put([]byte("l"), genesisBlock.Hash)
			if err != nil {
				log.Panic(err)
			}
			tip = genesisBlock.Hash
		}

		return nil
	})
	//返回区块链
	return  &BlockChain{tip, db}
}

//添加区块到区块链
func (bc *BlockChain)AddBlock(data string)  {
	err := bc.db.Update(func(tx *bolt.Tx) error {
		//获取表
		b := tx.Bucket([]byte(blockBucket))
		if b != nil {
			blockInDB := b.Get(bc.tip)
			block := DeserializeBlock(blockInDB)//反序列化
			newBlock := NewBlock(data, block.Height+1, block.Hash)//将block存储到db
			blockSerialize := newBlock.Serialize()
			err := b.Put(newBlock.Hash, blockSerialize)//存储区块
			if err != nil {
				log.Panic(err)
			}

			//更新db和区块链的tip
			err = b.Put([]byte("l"), newBlock.Hash)
			if err != nil {
				 log.Panic(err)
			}
			bc.tip = newBlock.Hash

		}
		return nil
	})
	if err != nil {
		log.Panic(err)
	}
}

// 返回Blockchain对象
func BlockchainObject() *BlockChain {
	db, err := bolt.Open(dbPath, 0600, nil)
	if err != nil {
		log.Fatal(err)
	}
	var tip []byte
	err = db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(blockBucket))
		if b != nil {
			// 读取最新区块的Hash
			tip = b.Get([]byte("l"))
		}

		return nil
	})
	return &BlockChain{tip,db}
}

func DBISExist() bool {
	if _, err := os.Stat(dbPath); os.IsNotExist(err) {
		return false
	}
	return true
}

