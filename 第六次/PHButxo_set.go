package main

import (
	"log"
	"github.com/boltdb/bolt"
	"encoding/hex"
)
const utxoBucket = "chainstate"
type PHBUTXOSet struct {
	PHBBlockchain *PHBBlockchain
}

// FindSpendableOutputs finds and returns unspent outputs to reference in inputs
func (u *PHBUTXOSet) FindSpendableOutputs(pubkeyHash []byte, amount int) (int, map[string][]int) {
	unspentOutputs := make(map[string][]int)
	accumulated := 0
	db := u.PHBBlockchain.phbdb

	err := db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(utxoBucket))
		c := b.Cursor()

		for k, v := c.First(); k != nil; k, v = c.Next() {
			txID := hex.EncodeToString(k)
			outs := PHBDeserializeOutputs(v)

			for outIdx, out := range outs.PHBOutPuts {
				if out.PHBIsLockedWithKey(pubkeyHash) && accumulated < amount {
					accumulated += out.PHBValue
					unspentOutputs[txID] = append(unspentOutputs[txID], outIdx)
				}
			}
		}

		return nil
	})
	if err != nil {
		log.Panic(err)
	}

	return accumulated, unspentOutputs
}

// FindUTXO finds UTXO for a public key hash
func (u *PHBUTXOSet) PHBFindUTXO(pubKeyHash []byte) []PHBTXOutput {
	var UTXOs []PHBTXOutput
	db := u.PHBBlockchain.phbdb

	err := db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(utxoBucket))
		c := b.Cursor()

		for k, v := c.First(); k != nil; k, v = c.Next() {
			outs := PHBDeserializeOutputs(v)

			for _, out := range outs.PHBOutPuts {
				if out.PHBIsLockedWithKey(pubKeyHash) {
					UTXOs = append(UTXOs, out)
				}
			}
		}

		return nil
	})
	if err != nil {
		log.Panic(err)
	}

	return UTXOs
}

// CountTransactions returns the number of transactions in the UTXO set
func (u *PHBUTXOSet) PHBCountTransactions() int {
	db := u.PHBBlockchain.phbdb
	counter := 0

	err := db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(utxoBucket))
		c := b.Cursor()

		for k, _ := c.First(); k != nil; k, _ = c.Next() {
			counter++
		}

		return nil
	})
	if err != nil {
		log.Panic(err)
	}

	return counter
}

// Reindex rebuilds the UTXO set
func (u *PHBUTXOSet) PHBReindex() {
	db := u.PHBBlockchain.phbdb
	bucketName := []byte(utxoBucket)

	err := db.Update(func(tx *bolt.Tx) error {
		err := tx.DeleteBucket(bucketName)
		if err != nil && err != bolt.ErrBucketNotFound {
			log.Panic(err)
		}

		_, err = tx.CreateBucket(bucketName)
		if err != nil {
			log.Panic(err)
		}

		return nil
	})
	if err != nil {
		log.Panic(err)
	}

	UTXO := u.PHBBlockchain.PHBFindUTXO()

	err = db.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket(bucketName)

		for txID, outs := range UTXO {
			key, err := hex.DecodeString(txID)
			if err != nil {
				log.Panic(err)
			}

			err = b.Put(key, outs.PHBSerialize())
			if err != nil {
				log.Panic(err)
			}
		}

		return nil
	})
}

// Update updates the UTXO set with transactions from the Block
// The Block is considered to be the tip of a blockchain
func (u *PHBUTXOSet) PHBUpdate(block *PHBBlock) {
	db := u.PHBBlockchain.phbdb

	err := db.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(utxoBucket))

		for _, tx := range block.PHBTransactions {
			if tx.PHBIsCoinbase() == false {
				for _, vin := range tx.PHBVin {
					updatedOuts := PHBTXOutputs{}
					outsBytes := b.Get(vin.PHBTxid)
					outs := PHBDeserializeOutputs(outsBytes)

					for outIdx, out := range outs.PHBOutPuts {
						if outIdx != vin.PHBVout {
							updatedOuts.PHBOutPuts = append(updatedOuts.PHBOutPuts, out)
						}
					}

					if len(updatedOuts.PHBOutPuts) == 0 {
						err := b.Delete(vin.PHBTxid)
						if err != nil {
							log.Panic(err)
						}
					} else {
						err := b.Put(vin.PHBTxid, updatedOuts.PHBSerialize())
						if err != nil {
							log.Panic(err)
						}
					}

				}
			}

			newOutputs := PHBTXOutputs{}
			for _, out := range tx.PHBVout {
				newOutputs.PHBOutPuts = append(newOutputs.PHBOutPuts, out)
			}

			err := b.Put(tx.PHBID, newOutputs.PHBSerialize())
			if err != nil {
				log.Panic(err)
			}
		}

		return nil
	})
	if err != nil {
		log.Panic(err)
	}
}