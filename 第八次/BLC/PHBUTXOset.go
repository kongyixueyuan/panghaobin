package BLC

import (
	"encoding/hex"
	"log"
	"github.com/boltdb/bolt"
	"fmt"
	"bytes"
)

const utxoTableName = "utxoTableName"

type PHBUTXOSet struct {
	PHBBlockchain *PHBBlockchain
}

func (utxoSet *PHBUTXOSet) PHBResetUTXOSet() {
	err := utxoSet.PHBBlockchain.PHBDB.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(utxoTableName))
		if b != nil {
			err := tx.DeleteBucket([]byte(utxoTableName))
			if err != nil {
				log.Panic(err)
			}
		}
		b, _ = tx.CreateBucket([]byte(utxoTableName))
		if b != nil {
			txOutputsMap := utxoSet.PHBBlockchain.PHBFindUTXOMap()
			for keyHash, outs := range txOutputsMap {
				txHash, _ := hex.DecodeString(keyHash)
				b.Put(txHash, outs.PHBSerialize())
			}
		}
		return nil
	})
	if err != nil {
		fmt.Println("重置失败....")
		log.Panic(err)
	}

}

func (utxoSet *PHBUTXOSet) phbfindUTXOForAddress(address string) []*PHBUTXO {
	var utxos []*PHBUTXO
	utxoSet.PHBBlockchain.PHBDB.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(utxoTableName))
		c := b.Cursor()
		for k, v := c.First(); k != nil; k, v = c.Next() {
			txOutputs := PHBDeserializeTXOutputs(v)
			for _, utxo := range txOutputs.PHBUTXOS {
				fmt.Println("$$$txHash:", hex.EncodeToString(utxo.PHBTxHash))
				utxo.PHBOutput.PHBPrintInfo()
				if utxo.PHBOutput.PHBUnLockScriptPubKeyWithAddress(address) {
					utxos = append(utxos, utxo)
				}
			}
		}
		return nil
	})
	return utxos
}

func (utxoSet *PHBUTXOSet) PHBGetBalance(address string) int64 {
	UTXOS := utxoSet.phbfindUTXOForAddress(address)
	var amount int64
	for _, utxo := range UTXOS {
		amount += utxo.PHBOutput.PHBValue
	}
	return amount
}

// 返回要凑多少钱，对应TXOutput的TX的Hash和index
func (utxoSet *PHBUTXOSet) FindUnPackageSpendableUTXOS(from string, txs []*PHBTransaction) []*PHBUTXO {
	var unUTXOs []*PHBUTXO
	spentTXOutputs := make(map[string][]int)
	for _, tx := range txs {
		if tx.PHBIsCoinbaseTransaction() == false {
			for _, in := range tx.PHBVins {
				//是否能够解锁
				publicKeyHash := PHBBase58Decode([]byte(from))
				ripemd160Hash := publicKeyHash[1:len(publicKeyHash)-4]
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
			if out.PHBUnLockScriptPubKeyWithAddress(from) {
				fmt.Println(from)
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

	return unUTXOs

}

func (utxoSet *PHBUTXOSet) PHBFindSpendableUTXOS(from string, amount int64, txs []*PHBTransaction) (int64, map[string][]int) {
	unPackageUTXOS := utxoSet.FindUnPackageSpendableUTXOS(from, txs)
	spentableUTXO := make(map[string][]int)
	var money int64 = 0
	for _, UTXO := range unPackageUTXOS {
		money += UTXO.PHBOutput.PHBValue
		txHash := hex.EncodeToString(UTXO.PHBTxHash)
		spentableUTXO[txHash] = append(spentableUTXO[txHash], UTXO.PHBIndex)
		if money >= amount {
			fmt.Println("$$$$unPackageUTXOS$$$$")
			return money, spentableUTXO
		}
	}
	// 钱还不够
	utxoSet.PHBBlockchain.PHBDB.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(utxoTableName))
		if b != nil {
			c := b.Cursor()
		UTXOBREAK:
			for k, v := c.First(); k != nil; k, v = c.Next() {
				txOutputs := PHBDeserializeTXOutputs(v)
				for _, utxo := range txOutputs.PHBUTXOS {
					if utxo.PHBOutput.PHBUnLockScriptPubKeyWithAddress(from) {
						fmt.Println("拿到的钱:", utxo.PHBOutput.PHBValue, "hash:", hex.EncodeToString(utxo.PHBTxHash))
						money += utxo.PHBOutput.PHBValue
						txHash := hex.EncodeToString(utxo.PHBTxHash)
						spentableUTXO[txHash] = append(spentableUTXO[txHash], utxo.PHBIndex)
						if money >= amount {
							break UTXOBREAK
						}
					}
				}
			}
		}
		return nil
	})
	if money < amount {
		log.Panic("余额不足......")
	}
	return money, spentableUTXO
}


func (utxoSet *PHBUTXOSet) PHBUpdate() {
	block := utxoSet.PHBBlockchain.PHBIterator().PHBNext()
	ins := []*PHBTXInput{}
	outsMap := make(map[string]*PHBTXOutputs)

	for _, tx := range block.PHBTxs {
		for _, in := range tx.PHBVins {
			ins = append(ins, in)
		}
	}
	//拿出新区快中所有的outputs
	for _, tx := range block.PHBTxs {
		utxos := []*PHBUTXO{}
		for index, out := range tx.PHBVouts {
			isSpent := false
			for _, in := range ins {
				if in.PHBVout == index && bytes.Compare(tx.PHBTxHash, in.PHBTxHash) == 0 && bytes.Compare(out.PHBRipemd160Hash, PHBRipemd160Hash(in.PHBPublicKey)) == 0 {
					fmt.Println("有已花费交易")
					isSpent = true
					continue
				}
			}
			if isSpent == false {
				utxo := &PHBUTXO{tx.PHBTxHash, index, out}
				utxos = append(utxos, utxo)
			}
		}
		if len(utxos) > 0 {
			txHash := hex.EncodeToString(tx.PHBTxHash)
			fmt.Println("outsMap:", txHash)
			outsMap[txHash] = &PHBTXOutputs{utxos}
		}
	}
	err := utxoSet.PHBBlockchain.PHBDB.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(utxoTableName))
		if b != nil {
			// 删除
			for _, in := range ins {
				txOutputsBytes := b.Get(in.PHBTxHash)
				if len(txOutputsBytes) == 0 {
					continue
				}
				fmt.Println("DeserializeTXOutputs")
				txOutputs := PHBDeserializeTXOutputs(txOutputsBytes)
				UTXOS := []*PHBUTXO{}
				// 判断是否需要
				isNeedDelete := false
				for _, utxo := range txOutputs.PHBUTXOS {
					if in.PHBVout == utxo.PHBIndex && bytes.Compare(utxo.PHBOutput.PHBRipemd160Hash, PHBRipemd160Hash(in.PHBPublicKey)) == 0 {
						isNeedDelete = true
					} else {
						UTXOS = append(UTXOS, utxo)
					}
				}
				if isNeedDelete {
					b.Delete(in.PHBTxHash)
					if len(UTXOS) > 0 {
						preTXOutputs := outsMap[hex.EncodeToString(in.PHBTxHash)]
						if preTXOutputs == nil {
							preTXOutputs = new(PHBTXOutputs)
						}
						preTXOutputs.PHBUTXOS = append(preTXOutputs.PHBUTXOS, UTXOS...)
						outsMap[hex.EncodeToString(in.PHBTxHash)] = preTXOutputs

					}
				}
			}
			for keyHash, outPuts := range outsMap {
				keyHashBytes, _ := hex.DecodeString(keyHash)
				b.Put(keyHashBytes, outPuts.PHBSerialize())
			}
		}
		return nil
	})
	if err != nil {
		log.Panic(err)
	}

}