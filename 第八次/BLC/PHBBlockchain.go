package BLC

import (
	"bytes"
	"crypto/ecdsa"
	"encoding/hex"
	"fmt"
	"github.com/boltdb/bolt"
	"log"
	"math/big"
	"os"
	"strconv"
	"time"
)

const dbName = "blockchain_%s.db"
const blockTableName = "blocks"

type PHBBlockchain struct {
	PHBTip []byte //最新的区块的Hash
	PHBDB  *bolt.DB
}

// 迭代器
func (blockchain *PHBBlockchain) PHBIterator() *PHBBlockchainIterator {

	return &PHBBlockchainIterator{blockchain.PHBTip, blockchain.PHBDB}
}

func PHBDBExists(dbName string) bool {
	if _, err := os.Stat(dbName); os.IsNotExist(err) {
		return false
	}
	return true
}

// 遍历输出所有区块的信息
func (blc *PHBBlockchain) PHBPrintchain() {
	fmt.Println("输出所有区块的信息....")
	blockchainIterator := blc.PHBIterator()
	for {
		fmt.Println("第一次进入for循环.....")
		block := blockchainIterator.PHBNext()

		fmt.Printf("Height：%d\n", block.PHBHeight)
		fmt.Printf("PrevBlockHash：%x\n", block.PHBPrevBlockHash)
		fmt.Printf("Timestamp：%s\n", time.Unix(block.PHBTimestamp, 0).Format("2006-01-02 03:04:05 PM"))
		fmt.Printf("Hash：%x\n", block.PHBHash)
		fmt.Printf("Nonce：%d\n", block.PHBNonce)
		fmt.Println("Txs:")
		for _, tx := range block.PHBTxs {
			fmt.Printf("%x\n", tx.PHBTxHash)
			fmt.Println("Vins:")
			for _, in := range tx.PHBVins {
				fmt.Printf("%x\n", in.PHBTxHash)
				fmt.Printf("%d\n", in.PHBVout)
				fmt.Printf("%x\n", in.PHBPublicKey)
			}

			fmt.Println("Vouts:")
			for _, out := range tx.PHBVouts {
				//fmt.Println(out.Value)
				fmt.Printf("%d\n", out.PHBValue)
				//fmt.Println(out.Ripemd160Hash)
				fmt.Printf("%x\n", out.PHBRipemd160Hash)
			}
		}
		fmt.Println("------------------------------")
		var hashInt big.Int
		hashInt.SetBytes(block.PHBPrevBlockHash)
		if big.NewInt(0).Cmp(&hashInt) == 0 {
			break
		}
	}

}


func (blc *PHBBlockchain) AddBlockToBlockchain(txs []*PHBTransaction) {
	err := blc.PHBDB.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(blockTableName))
		if b != nil {
			blockBytes := b.Get(blc.PHBTip)
			block := PHBDeserializeBlock(blockBytes)
			newBlock := PHBNewBlock(txs, block.PHBHeight+1, block.PHBHash)
			err := b.Put(newBlock.PHBHash, newBlock.PHBSerialize())
			if err != nil {
				log.Panic(err)
			}
			err = b.Put([]byte("l"), newBlock.PHBHash)
			if err != nil {
				log.Panic(err)
			}
			blc.PHBTip = newBlock.PHBHash
		}
		return nil
	})

	if err != nil {
		log.Panic(err)
	}
}

func PHBCreateBlockchainWithGenesisBlock(address string, nodeID string) *PHBBlockchain {
	dbName := fmt.Sprintf(dbName, nodeID)
	if PHBDBExists(dbName) {
		fmt.Println("创世区块已经存在.......")
		os.Exit(1)
	}
	fmt.Println("正在创建创世区块.......")
	db, err := bolt.Open(dbName, 0600, nil)
	if err != nil {
		log.Fatal(err)
	}

	var genesisHash []byte
	err = db.Update(func(tx *bolt.Tx) error {
		b, err := tx.CreateBucket([]byte(blockTableName))
		if err != nil {
			log.Panic(err)
		}
		if b != nil {
			txCoinbase := PHBNewCoinbaseTransaction(address)
			genesisBlock := PHBCreateGenesisBlock([]*PHBTransaction{txCoinbase})
			err := b.Put(genesisBlock.PHBHash, genesisBlock.PHBSerialize())
			if err != nil {
				log.Panic(err)
			}
			err = b.Put([]byte("l"), genesisBlock.PHBHash)
			if err != nil {
				log.Panic(err)
			}
			genesisHash = genesisBlock.PHBHash
		}

		return nil
	})

	return &PHBBlockchain{genesisHash, db}

}

// 返回Blockchain对象
func PHBBlockchainObject(nodeID string) *PHBBlockchain {
	dbName := fmt.Sprintf(dbName, nodeID)
	if PHBDBExists(dbName) == false {
		fmt.Println("数据库不存在....")
		os.Exit(1)
	}
	db, err := bolt.Open(dbName, 0600, nil)
	if err != nil {
		log.Fatal(err)
	}
	var tip []byte
	err = db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(blockTableName))
		if b != nil {
			// 读取最新区块的Hash
			tip = b.Get([]byte("l"))
		}
		return nil
	})
	return &PHBBlockchain{tip, db}
}

func (blockchain *PHBBlockchain) PHBUnUTXOs(address string, txs []*PHBTransaction) []*PHBUTXO {
	var unUTXOs []*PHBUTXO
	spentTXOutputs := make(map[string][]int)
	for _, tx := range txs {
		if tx.PHBIsCoinbaseTransaction() == false {
			for _, in := range tx.PHBVins {
				publicKeyHash := PHBBase58Decode([]byte(address))
				ripemd160Hash := publicKeyHash[1 : len(publicKeyHash)-4]
				if in.PHBUnLockRipemd160Hash(ripemd160Hash) {
					key := hex.EncodeToString(in.PHBTxHash)
					spentTXOutputs[key] = append(spentTXOutputs[key], in.PHBVout)
				}

			}
		}
	}

	for _, tx := range txs {

	Work1:
		for index, out := range tx.PHBVouts {

			if out.PHBUnLockScriptPubKeyWithAddress(address) {
				fmt.Println("看看是否是俊诚...")
				fmt.Println(address)

				fmt.Println(spentTXOutputs)

				if len(spentTXOutputs) == 0 {
					utxo := &PHBUTXO{tx.PHBTxHash, index, out}
					unUTXOs = append(unUTXOs, utxo)
				} else {
					for hash, indexArray := range spentTXOutputs {
						txHashStr := hex.EncodeToString(tx.PHBTxHash)
						if hash == txHashStr {

							var isUnSpentUTXO bool
							for _, outIndex := range indexArray {
								if index == outIndex {
									isUnSpentUTXO = true
									continue Work1
								}
								if isUnSpentUTXO == false {
									utxo := &PHBUTXO{tx.PHBTxHash, index, out}
									unUTXOs = append(unUTXOs, utxo)
								}
							}
						} else {
							utxo := &PHBUTXO{tx.PHBTxHash, index, out}
							unUTXOs = append(unUTXOs, utxo)
						}
					}
				}

			}

		}

	}

	blockIterator := blockchain.PHBIterator()

	for {
		block := blockIterator.PHBNext()
		fmt.Println(block)
		fmt.Println()
		for i := len(block.PHBTxs) - 1; i >= 0; i-- {
			tx := block.PHBTxs[i]
			if tx.PHBIsCoinbaseTransaction() == false {
				for _, in := range tx.PHBVins {
					//是否能够解锁
					publicKeyHash := PHBBase58Decode([]byte(address))
					ripemd160Hash := publicKeyHash[1 : len(publicKeyHash)-4]

					if in.PHBUnLockRipemd160Hash(ripemd160Hash) {
						key := hex.EncodeToString(in.PHBTxHash)
						spentTXOutputs[key] = append(spentTXOutputs[key], in.PHBVout)
					}
				}
			}
		work:
			for index, out := range tx.PHBVouts {

				if out.PHBUnLockScriptPubKeyWithAddress(address) {
					fmt.Println(out)
					fmt.Println(spentTXOutputs)
					if spentTXOutputs != nil {
						if len(spentTXOutputs) != 0 {
							var isSpentUTXO bool
							for txHash, indexArray := range spentTXOutputs {
								for _, i := range indexArray {
									if index == i && txHash == hex.EncodeToString(tx.PHBTxHash) {
										isSpentUTXO = true
										continue work
									}
								}
							}
							if isSpentUTXO == false {
								utxo := &PHBUTXO{tx.PHBTxHash, index, out}
								unUTXOs = append(unUTXOs, utxo)
							}
						} else {
							utxo := &PHBUTXO{tx.PHBTxHash, index, out}
							unUTXOs = append(unUTXOs, utxo)
						}

					}
				}

			}
		}
		fmt.Println(spentTXOutputs)
		var hashInt big.Int
		hashInt.SetBytes(block.PHBPrevBlockHash)
		if hashInt.Cmp(big.NewInt(0)) == 0 {
			break
		}
	}
	return unUTXOs
}

// 转账时查找可用的UTXO
func (blockchain *PHBBlockchain) PHBFindSpendableUTXOS(from string, amount int, txs []*PHBTransaction) (int64, map[string][]int) {
	utxos := blockchain.PHBUnUTXOs(from, txs)
	spendableUTXO := make(map[string][]int)
	var value int64
	for _, utxo := range utxos {
		value = value + utxo.PHBOutput.PHBValue
		hash := hex.EncodeToString(utxo.PHBTxHash)
		spendableUTXO[hash] = append(spendableUTXO[hash], utxo.PHBIndex)
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

// 挖掘新的区块
func (blockchain *PHBBlockchain) PHBMineNewBlock(from []string, to []string, amount []string, nodeID string) {
	utxoSet := &PHBUTXOSet{blockchain}
	var txs []*PHBTransaction
	for index, address := range from {
		value, _ := strconv.Atoi(amount[index])
		tx := PHBNewSimpleTransaction(address, to[index], int64(value), utxoSet, txs, nodeID)
		txs = append(txs, tx)
	}
	tx := PHBNewCoinbaseTransaction(from[0])
	txs = append(txs, tx)
	var block *PHBBlock
	blockchain.PHBDB.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(blockTableName))
		if b != nil {
			hash := b.Get([]byte("l"))
			blockBytes := b.Get(hash)
			block = PHBDeserializeBlock(blockBytes)
		}
		return nil
	})
	_txs := []*PHBTransaction{}
	for _, tx := range txs {
		if blockchain.PHBVerifyTransaction(tx, _txs) != true {
			log.Panic("ERROR: Invalid transaction")
		}
		_txs = append(_txs, tx)
	}
	//建立新的区块
	block = PHBNewBlock(txs, block.PHBHeight+1, block.PHBHash)
	//将新区块存储到数据库
	blockchain.PHBDB.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(blockTableName))
		if b != nil {
			b.Put(block.PHBHash, block.PHBSerialize())
			b.Put([]byte("l"), block.PHBHash)
			blockchain.PHBTip = block.PHBHash
		}
		return nil
	})

}

// 查询余额
func (blockchain *PHBBlockchain) PHBGetBalance(address string) int64 {
	utxos := blockchain.PHBUnUTXOs(address, []*PHBTransaction{})
	var amount int64
	for _, utxo := range utxos {
		amount = amount + utxo.PHBOutput.PHBValue
	}
	return amount
}

func (bclockchain *PHBBlockchain) PHBSignTransaction(tx *PHBTransaction, privKey ecdsa.PrivateKey, txs []*PHBTransaction) {
	if tx.PHBIsCoinbaseTransaction() {
		return
	}
	prevTXs := make(map[string]PHBTransaction)
	for _, vin := range tx.PHBVins {
		prevTX, err := bclockchain.PHBFindTransaction(vin.PHBTxHash, txs)
		if err != nil {
			log.Panic(err)
		}
		prevTXs[hex.EncodeToString(prevTX.PHBTxHash)] = prevTX
	}
	tx.PHBSign(privKey, prevTXs)
}

func (bc *PHBBlockchain) PHBFindTransaction(ID []byte, txs []*PHBTransaction) (PHBTransaction, error) {
	for _, tx := range txs {
		if bytes.Compare(tx.PHBTxHash, ID) == 0 {
			return *tx, nil
		}
	}
	bci := bc.PHBIterator()
	for {
		block := bci.PHBNext()
		for _, tx := range block.PHBTxs {
			if bytes.Compare(tx.PHBTxHash, ID) == 0 {
				return *tx, nil
			}
		}
		var hashInt big.Int
		hashInt.SetBytes(block.PHBPrevBlockHash)
		if big.NewInt(0).Cmp(&hashInt) == 0 {
			break
		}
	}
	return PHBTransaction{}, nil
}

// 验证数字签名
func (bc *PHBBlockchain) PHBVerifyTransaction(tx *PHBTransaction, txs []*PHBTransaction) bool {
	prevTXs := make(map[string]PHBTransaction)
	for _, vin := range tx.PHBVins {
		prevTX, err := bc.PHBFindTransaction(vin.PHBTxHash, txs)
		if err != nil {
			log.Panic(err)
		}
		prevTXs[hex.EncodeToString(prevTX.PHBTxHash)] = prevTX
	}
	return tx.PHBVerify(prevTXs)
}

// [string]*TXOutputs
func (blc *PHBBlockchain) PHBFindUTXOMap() map[string]*PHBTXOutputs {
	blcIterator := blc.PHBIterator()
	spentableUTXOsMap := make(map[string][]*PHBTXInput)
	utxoMaps := make(map[string]*PHBTXOutputs)
	for {
		block := blcIterator.PHBNext()
		for i := len(block.PHBTxs) - 1; i >= 0; i-- {
			txOutputs := &PHBTXOutputs{[]*PHBUTXO{}}
			tx := block.PHBTxs[i]
			if tx.PHBIsCoinbaseTransaction() == false {
				for _, txInput := range tx.PHBVins {
					txHash := hex.EncodeToString(txInput.PHBTxHash)
					spentableUTXOsMap[txHash] = append(spentableUTXOsMap[txHash], txInput)
				}
			}
			txHash := hex.EncodeToString(tx.PHBTxHash)
			txInputs := spentableUTXOsMap[txHash]
			if len(txInputs) > 0 {
			WorkOutLoop:
				for index, out := range tx.PHBVouts {
					for _, in := range txInputs {
						outPublicKey := out.PHBRipemd160Hash
						inPublicKey := in.PHBPublicKey
						if bytes.Compare(outPublicKey, PHBRipemd160Hash(inPublicKey)) == 0 {
							if index == in.PHBVout {
								continue WorkOutLoop
							} else {
								utxo := &PHBUTXO{tx.PHBTxHash, index, out}
								txOutputs.PHBUTXOS = append(txOutputs.PHBUTXOS, utxo)
							}
						}
					}
				}
			} else {

				for index, out := range tx.PHBVouts {
					utxo := &PHBUTXO{tx.PHBTxHash, index, out}
					txOutputs.PHBUTXOS = append(txOutputs.PHBUTXOS, utxo)
				}
			}
			utxoMaps[txHash] = txOutputs

		}
		// 找到创世区块时退出
		var hashInt big.Int
		hashInt.SetBytes(block.PHBPrevBlockHash)
		if hashInt.Cmp(big.NewInt(0)) == 0 {
			break
		}
	}

	return utxoMaps
}

func (bc *PHBBlockchain) PHBGetBestHeight() int64 {
	block := bc.PHBIterator().PHBNext()
	return block.PHBHeight
}

func (bc *PHBBlockchain) PHBGetBlockHashes() [][]byte {
	blockIterator := bc.PHBIterator()
	var blockHashs [][]byte
	for {
		block := blockIterator.PHBNext()
		blockHashs = append(blockHashs, block.PHBHash)
		var hashInt big.Int
		hashInt.SetBytes(block.PHBPrevBlockHash)
		if hashInt.Cmp(big.NewInt(0)) == 0 {
			break
		}
	}
	return blockHashs
}

func (bc *PHBBlockchain) PHBGetBlock(blockHash []byte) ([]byte, error) {
	var blockBytes []byte
	err := bc.PHBDB.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(blockTableName))
		if b != nil {
			blockBytes = b.Get(blockHash)
		}
		return nil
	})
	return blockBytes, err
}

func (bc *PHBBlockchain) PHBAddBlock(block *PHBBlock) {
	err := bc.PHBDB.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(blockTableName))
		if b != nil {
			blockExist := b.Get(block.PHBHash)
			if blockExist != nil {
				// 如果存在，不需要做任何过多的处理
				return nil
			}
			err := b.Put(block.PHBHash, block.PHBSerialize())

			if err != nil {
				log.Panic(err)
			}
			blockHash := b.Get([]byte("l"))
			blockBytes := b.Get(blockHash)
			blockInDB := PHBDeserializeBlock(blockBytes)
			if blockInDB.PHBHeight < block.PHBHeight {

				b.Put([]byte("l"), block.PHBHash)
				bc.PHBTip = block.PHBHash
			}
		}
		return nil
	})
	if err != nil {
		log.Panic(err)
	}
}
