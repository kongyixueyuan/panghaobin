package BLC

import (
	"github.com/boltdb/bolt"
	"log"
	"os"
	"encoding/hex"
	"fmt"
	"math/big"
	"strconv"
)

const dbPath = "blockchain.db"
const blockBucket = "blocks"
type BlockChain struct {
	tip []byte
	db *bolt.DB
}

func CreatBlockChain(address string, values int64) *BlockChain  {
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
			txCoinbase := CreateCoinbaseTransaction(address, values)
			genesisBlock := NewGenesisBlock([]*Transaction{txCoinbase})
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
func (bc *BlockChain)AddBlock(txs []*Transaction)  {
	err := bc.db.Update(func(tx *bolt.Tx) error {
		//获取表
		b := tx.Bucket([]byte(blockBucket))
		if b != nil {
			blockInDB := b.Get(bc.tip)
			block := DeserializeBlock(blockInDB)//反序列化
			newBlock := NewBlock(txs, block.Height+1, block.Hash)//将block存储到db
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

//如果一个对应的txouput未花费，那么返回这个transaction
func (bc *BlockChain) UnUTXOs(address string, txs []*Transaction) []*UTXO {
	var unUTXOs []*UTXO
	spentTXOutputs := make(map[string][]int)
	for _, tx := range txs {
		if tx.IsCoinbaseTransaction() == false {
			for _, in := range tx.Vin {
				if in.IsUserKey(address) {
					key := hex.EncodeToString(in.TxHash)
					spentTXOutputs[key] = append(spentTXOutputs[key], in.Vout)
				}
			}
		}
	}

	for _,tx := range txs {
	Work1:
		for index,out := range tx.Vout {

			if out.UnLockPubKeyWithAddress(address) {
				fmt.Println(address)

				fmt.Println(spentTXOutputs)
				if len(spentTXOutputs) == 0 {
					utxo := &UTXO{tx.TxHash, index, out}
					unUTXOs = append(unUTXOs, utxo)
				} else {
					for hash,indexArray := range spentTXOutputs {
						txHashStr := hex.EncodeToString(tx.TxHash)
						if hash == txHashStr {
							var isUnSpentUTXO bool
							for _,outIndex := range indexArray {
								if index == outIndex {
									isUnSpentUTXO = true
									continue Work1
								}
								if isUnSpentUTXO == false {
									utxo := &UTXO{tx.TxHash, index, out}
									unUTXOs = append(unUTXOs, utxo)
								}
							}
						} else {
							utxo := &UTXO{tx.TxHash, index, out}
							unUTXOs = append(unUTXOs, utxo)
						}
					}
				}
			}
		}
	}
	blockIterator := bc.Iterator()
	fmt.Println("#################")
	for {
		block := blockIterator.Next()
		block.printfBlock()
		for i := len(block.Txs) - 1; i >= 0 ; i-- {
			tx := block.Txs[i]
			if tx.IsCoinbaseTransaction() == false {
				for _, in := range tx.Vin {
					//是否能够解锁
					if in.IsUserKey(address) {
						key := hex.EncodeToString(in.TxHash)
						spentTXOutputs[key] = append(spentTXOutputs[key], in.Vout)
					}
				}
			}
		work://out
			for index, out := range tx.Vout {
				if out.UnLockPubKeyWithAddress(address) {
					out.printfTxOutput()
					fmt.Println(spentTXOutputs)
					if spentTXOutputs != nil {
						if len(spentTXOutputs) != 0 {
							var isSpentUTXO bool
							for txHash, indexArray := range spentTXOutputs {
								for _, i := range indexArray {
									if index == i && txHash == hex.EncodeToString(tx.TxHash) {
										isSpentUTXO = true
										continue work
									}
								}
							}
							if isSpentUTXO == false {
								utxo := &UTXO{tx.TxHash, index, out}
								unUTXOs = append(unUTXOs, utxo)
							}
						} else {
							utxo := &UTXO{tx.TxHash, index, out}
							unUTXOs = append(unUTXOs, utxo)
						}
					}
				}

			}
		}

		fmt.Println(spentTXOutputs)
		var hashInt big.Int
		hashInt.SetBytes(block.PrevBlockHash)
		if hashInt.Cmp(big.NewInt(0)) == 0 {
			break;
		}

	}

	return unUTXOs
}

//转账时查找可用的UTXO
func (bc *BlockChain) FindSpendableUTXOs(from string, amount int, txs []*Transaction) (int64, map[string][]int) {
	utxos := bc.UnUTXOs(from,txs)
	spendableUTXO := make(map[string][]int)
	//2. 遍历utxos
	var value int64
	for _, utxo := range utxos {
		value = value + utxo.OutPut.Value
		hash := hex.EncodeToString(utxo.TxHash)
		spendableUTXO[hash] = append(spendableUTXO[hash], utxo.Index)
		if value >= int64(amount) {
			break
		}
	}
	if value < int64(amount) {
		fmt.Printf("%s's fund is 不足\n", from)
		os.Exit(1)
	}
	return value, spendableUTXO
}

// 挖矿
func (bc *BlockChain) MineNewBlock(from []string, to []string, amount []string) {

	//1.建立一笔交易
	fmt.Println(from)
	fmt.Println(to)
	fmt.Println(amount)
	var txs []*Transaction
	for index,address := range from {
		value, _ := strconv.Atoi(amount[index])
		tx := CreateSimpleTransaction(address, to[index], value, bc,txs)
		txs = append(txs, tx)
	}
	//1.通过相关算法建立Transaction数组
	var block *Block
	bc.db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(blockBucket))
		if b != nil {
			hash := b.Get([]byte("l"))
			blockBytes := b.Get(hash)
			block = DeserializeBlock(blockBytes)
		}
		return nil
	})
	//2.建立新区块
	block = NewBlock(txs, block.Height+1, block.Hash)
	//存储到数据库
	bc.db.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(blockBucket))
		if b != nil {
			b.Put(block.Hash, block.Serialize())
			b.Put([]byte("l"), block.Hash)
			bc.tip = block.Hash
		}
		return nil
	})
}

// 查询余额
func (bc *BlockChain) GetBalance(address string) int64 {
	utxos := bc.UnUTXOs(address,[]*Transaction{})
	var amount int64
	for _, utxo := range utxos {
		amount = amount + utxo.OutPut.Value
	}
	return amount
}