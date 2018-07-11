package main

import (
	"github.com/boltdb/bolt"
	"fmt"
	"os"
	"log"
	"bytes"
	"errors"
	"encoding/hex"
	"crypto/ecdsa"
)

const dbFile = "blockchain_%s.db"
const blocksBucket = "blocks"
const genesisCoinbaseData = "创世区块"
type PHBBlockchain struct {
	phbtip []byte
	phbdb *bolt.DB
}

func PHBCreateBlockchain(address, nodeID string) *PHBBlockchain {
	dbFile := fmt.Sprintf(dbFile, nodeID)
	if phbDBExist(dbFile) {
		fmt.Println("Blockchain already exists.")
		os.Exit(1)
	}

	var tip []byte

	cbtx := PHBNewCoinbaseTX(address, genesisCoinbaseData)
	genesis := PHBNewGenesisBlock(cbtx)

	db, err := bolt.Open(dbFile, 0600, nil)
	if err != nil {
		log.Panic(err)
	}

	err = db.Update(func(tx *bolt.Tx) error {
		b, err := tx.CreateBucket([]byte(blocksBucket))
		if err != nil {
			log.Panic(err)
		}

		err = b.Put(genesis.PHBHash, genesis.PHBSerialize())
		if err != nil {
			log.Panic(err)
		}

		err = b.Put([]byte("l"), genesis.PHBHash)
		if err != nil {
			log.Panic(err)
		}
		tip = genesis.PHBHash

		return nil
	})
	if err != nil {
		log.Panic(err)
	}

	bc := PHBBlockchain{tip, db}

	return &bc
}

// NewBlockchain creates a new Blockchain with genesis Block
func PHBNewBlockchain(nodeID string) *PHBBlockchain {
	dbFile := fmt.Sprintf(dbFile, nodeID)
	if phbDBExist(dbFile) == false {
		fmt.Println("No existing blockchain found. Create one first.")
		os.Exit(1)
	}

	var tip []byte
	db, err := bolt.Open(dbFile, 0600, nil)
	if err != nil {
		log.Panic(err)
	}

	err = db.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(blocksBucket))
		tip = b.Get([]byte("l"))

		return nil
	})
	if err != nil {
		log.Panic(err)
	}

	bc := PHBBlockchain{tip, db}

	return &bc
}

// AddBlock saves the block into the blockchain
func (bc *PHBBlockchain) PHBAddBlock(block *PHBBlock) {
	err := bc.phbdb.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(blocksBucket))
		blockInDb := b.Get(block.PHBHash)

		if blockInDb != nil {
			return nil
		}

		blockData := block.PHBSerialize()
		err := b.Put(block.PHBHash, blockData)
		if err != nil {
			log.Panic(err)
		}

		lastHash := b.Get([]byte("l"))
		lastBlockData := b.Get(lastHash)
		lastBlock := PHBDeserializeBlock(lastBlockData)

		if block.PHBHeight > lastBlock.PHBHeight {
			err = b.Put([]byte("l"), block.PHBHash)
			if err != nil {
				log.Panic(err)
			}
			bc.phbtip = block.PHBHash
		}

		return nil
	})
	if err != nil {
		log.Panic(err)
	}
}

// FindTransaction finds a transaction by its ID
func (bc *PHBBlockchain) PHBFindTransaction(ID []byte) (PHBTransaction, error) {
	bci := bc.PHBIterator()

	for {
		block := bci.PHBNext()

		for _, tx := range block.PHBTransactions {
			if bytes.Compare(tx.PHBID, ID) == 0 {
				return *tx, nil
			}
		}

		if len(block.PHBPrevBlockHash) == 0 {
			break
		}
	}

	return PHBTransaction{}, errors.New("Transaction is not found")
}

// FindUTXO finds all unspent transaction outputs and returns transactions with spent outputs removed
func (bc *PHBBlockchain) PHBFindUTXO() map[string]PHBTXOutputs {
	UTXO := make(map[string]PHBTXOutputs)
	spentTXOs := make(map[string][]int)
	bci := bc.PHBIterator()

	for {
		block := bci.PHBNext()

		for _, tx := range block.PHBTransactions {
			txID := hex.EncodeToString(tx.PHBID)

		Outputs:
			for outIdx, out := range tx.PHBVout {
				// Was the output spent?
				if spentTXOs[txID] != nil {
					for _, spentOutIdx := range spentTXOs[txID] {
						if spentOutIdx == outIdx {
							continue Outputs
						}
					}
				}

				outs := UTXO[txID]
				outs.PHBOutPuts = append(outs.PHBOutPuts, out)
				UTXO[txID] = outs
			}

			if tx.PHBIsCoinbase() == false {
				for _, in := range tx.PHBVin {
					inTxID := hex.EncodeToString(in.PHBTxid)
					spentTXOs[inTxID] = append(spentTXOs[inTxID], in.PHBVout)
				}
			}
		}

		if len(block.PHBPrevBlockHash) == 0 {
			break
		}
	}

	return UTXO
}

// Iterator returns a BlockchainIterat
func (bc *PHBBlockchain) PHBIterator() *PHBBlockchainIterator {
	bci := &PHBBlockchainIterator{bc.phbtip, bc.phbdb}

	return bci
}

// GetBestHeight returns the height of the latest block
func (bc *PHBBlockchain) PHBGetBestHeight() int {
	var lastBlock PHBBlock

	err := bc.phbdb.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(blocksBucket))
		lastHash := b.Get([]byte("l"))
		blockData := b.Get(lastHash)
		lastBlock = *PHBDeserializeBlock(blockData)

		return nil
	})
	if err != nil {
		log.Panic(err)
	}

	return lastBlock.PHBHeight
}

// GetBlock finds a block by its hash and returns it
func (bc *PHBBlockchain) PHBGetBlock(blockHash []byte) (PHBBlock, error) {
	var block PHBBlock

	err := bc.phbdb.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(blocksBucket))

		blockData := b.Get(blockHash)

		if blockData == nil {
			return errors.New("Block is not found.")
		}

		block = *PHBDeserializeBlock(blockData)

		return nil
	})
	if err != nil {
		return block, err
	}

	return block, nil
}

// GetBlockHashes returns a list of hashes of all the blocks in the chain
func (bc *PHBBlockchain) PHBGetBlockHashes() [][]byte {
	var blocks [][]byte
	bci := bc.PHBIterator()

	for {
		block := bci.PHBNext()

		blocks = append(blocks, block.PHBHash)

		if len(block.PHBPrevBlockHash) == 0 {
			break
		}
	}

	return blocks
}

// MineBlock mines a new block with the provided transactions
func (bc *PHBBlockchain) PHBMineBlock(transactions []*PHBTransaction) *PHBBlock {
	var lastHash []byte
	var lastHeight int

	for _, tx := range transactions {
		// TODO: ignore transaction if it's not valid
		if bc.PHBVerifyTransaction(tx) != true {
			log.Panic("ERROR: Invalid transaction")
		}
	}

	err := bc.phbdb.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(blocksBucket))
		lastHash = b.Get([]byte("l"))

		blockData := b.Get(lastHash)
		block := PHBDeserializeBlock(blockData)

		lastHeight = block.PHBHeight

		return nil
	})
	if err != nil {
		log.Panic(err)
	}

	newBlock := PHBNewBlock(transactions, lastHash, lastHeight+1)

	err = bc.phbdb.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(blocksBucket))
		err := b.Put(newBlock.PHBHash, newBlock.PHBSerialize())
		if err != nil {
			log.Panic(err)
		}

		err = b.Put([]byte("l"), newBlock.PHBHash)
		if err != nil {
			log.Panic(err)
		}

		bc.phbtip = newBlock.PHBHash

		return nil
	})
	if err != nil {
		log.Panic(err)
	}

	return newBlock
}

// SignTransaction signs inputs of a Transaction
func (bc *PHBBlockchain) PHBSignTransaction(tx *PHBTransaction, privKey ecdsa.PrivateKey) {
	prevTXs := make(map[string]PHBTransaction)

	for _, vin := range tx.PHBVin {
		prevTX, err := bc.PHBFindTransaction(vin.PHBTxid)
		if err != nil {
			log.Panic(err)
		}
		prevTXs[hex.EncodeToString(prevTX.PHBID)] = prevTX
	}

	tx.PHBSign(privKey, prevTXs)
}

// VerifyTransaction verifies transaction input signatures
func (bc *PHBBlockchain) PHBVerifyTransaction(tx *PHBTransaction) bool {
	if tx.PHBIsCoinbase() {
		return true
	}

	prevTXs := make(map[string]PHBTransaction)

	for _, vin := range tx.PHBVin {
		prevTX, err := bc.PHBFindTransaction(vin.PHBTxid)
		if err != nil {
			log.Panic(err)
		}
		prevTXs[hex.EncodeToString(prevTX.PHBID)] = prevTX
	}

	return tx.PHBVerify(prevTXs)
}

func phbDBExist(dbFile string) bool {
	if _, err := os.Stat(dbFile); os.IsNotExist(err) {
		return false
	}
	return true
}
