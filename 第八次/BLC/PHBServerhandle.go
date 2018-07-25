package BLC

import (
	"bytes"
	"encoding/gob"
	"encoding/hex"
	"fmt"
	"github.com/boltdb/bolt"
	"log"
	"os"
)

func phbhandleVersion(request []byte, bc *PHBBlockchain) {

	var buff bytes.Buffer
	var payload PHBVersion
	dataBytes := request[COMMANDLENGTH:]
	// 反序列化
	buff.Write(dataBytes)
	dec := gob.NewDecoder(&buff)
	err := dec.Decode(&payload)
	if err != nil {
		log.Panic(err)
	}
	bestHeight := bc.PHBGetBestHeight()
	foreignerBestHeight := payload.PHBBestHeight
	if bestHeight > foreignerBestHeight {
		phbsendVersion(payload.PHBAddrFrom, bc)
	} else if bestHeight < foreignerBestHeight {
		// 去向主节点要信息
		phbsendGetBlocks(payload.PHBAddrFrom)
	}
	if !phbnodeIsKnown(payload.PHBAddrFrom) {
		knowNodes = append(knowNodes, payload.PHBAddrFrom)
	}

}

func phbhandleGetblocks(request []byte, bc *PHBBlockchain) {
	var buff bytes.Buffer
	var payload PHBGetBlocks
	dataBytes := request[COMMANDLENGTH:]
	// 反序列化
	buff.Write(dataBytes)
	dec := gob.NewDecoder(&buff)
	err := dec.Decode(&payload)
	if err != nil {
		log.Panic(err)
	}
	blocks := bc.PHBGetBlockHashes()
	phbsendInv(payload.PHBAddrFrom, BLOCK_TYPE, blocks)

}

func phbhandleGetData(request []byte, bc *PHBBlockchain) {
	var buff bytes.Buffer
	var payload PHBGetData
	dataBytes := request[COMMANDLENGTH:]
	buff.Write(dataBytes)
	dec := gob.NewDecoder(&buff)
	err := dec.Decode(&payload)
	if err != nil {
		log.Panic(err)
	}
	if payload.PHBType == BLOCK_TYPE {

		block, err := bc.PHBGetBlock([]byte(payload.PHBHash))
		if err != nil {
			return
		}
		phbsendBlock(payload.PHBAddrFrom, block)
	}
	if payload.PHBType == TX_TYPE {
		tx := memoryTxPool[hex.EncodeToString(payload.PHBHash)]
		phbsendTx(payload.PHBAddrFrom, tx)
	}
}

func phbhandleAddr(request []byte,bc *PHBBlockchain)  {
}

func phbhandleBlock(request []byte, bc *PHBBlockchain) {
	var buff bytes.Buffer
	var payload PHBBlockData
	dataBytes := request[COMMANDLENGTH:]
	// 反序列化
	buff.Write(dataBytes)
	dec := gob.NewDecoder(&buff)
	err := dec.Decode(&payload)
	if err != nil {
		log.Panic(err)
	}
	blockBytes := payload.PHBBlock
	block := PHBDeserializeBlock(blockBytes)
	fmt.Println("Recevied a new block!")
	bc.PHBAddBlock(block)
	fmt.Printf("Added block %x\n", block.PHBHash)
	if len(transactionArray) > 0 {
		blockHash := transactionArray[0]
		phbsendGetData(payload.PHBAddrFrom, "block", blockHash)
		transactionArray = transactionArray[1:]
	} else {
		fmt.Println("数据库重置......")
		UTXOSet := &PHBUTXOSet{bc}
		UTXOSet.PHBResetUTXOSet()
	}
}

func phbhandleTx(request []byte, bc *PHBBlockchain) {

	var buff bytes.Buffer
	var payload PHBTx
	dataBytes := request[COMMANDLENGTH:]
	buff.Write(dataBytes)
	dec := gob.NewDecoder(&buff)
	err := dec.Decode(&payload)
	if err != nil {
		log.Panic("发序列化错误:", err)
	}
	tx := payload.PHBTx
	memoryTxPool[hex.EncodeToString(tx.PHBTxHash)] = tx

	// 说明主节点自己
	if nodeAddress == knowNodes[0] {
		// 给矿工节点发送交易hash
		for _, nodeAddr := range knowNodes {
			if nodeAddr != nodeAddress && nodeAddr != payload.PHBAddrFrom {
				phbsendInv(nodeAddr, TX_TYPE, [][]byte{tx.PHBTxHash})
			}
		}
	}

	if len(minerAddress) > 0 {
		bc.PHBDB.Close()
		blockchain := PHBBlockchainObject(os.Getenv("NODE_ID"))
		defer blockchain.PHBDB.Close()
		utxoSet := &PHBUTXOSet{blockchain}
		var txs []*PHBTransaction
		txs = append(txs, tx)
		//奖励
		coinTX := PHBNewCoinbaseTransaction(minerAddress)
		txs = append(txs, coinTX)
		//1. 通过相关算法建立Transaction数组
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

		// 在建立新区块之前对txs进行签名验证
		_txs := []*PHBTransaction{}
		for _, tx := range txs {
			if blockchain.PHBVerifyTransaction(tx, _txs) != true {
				log.Panic("ERROR: Invalid transaction")
			}

			_txs = append(_txs, tx)
		}
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
		//转账成功以后，需要更新一下
		utxoSet.PHBUpdate()
		phbsendBlock(knowNodes[0], block.PHBSerialize())
	}

}

func phbhandleInv(request []byte, bc *PHBBlockchain) {

	var buff bytes.Buffer
	var payload PHBInv
	dataBytes := request[COMMANDLENGTH:]
	// 反序列化
	buff.Write(dataBytes)
	dec := gob.NewDecoder(&buff)
	err := dec.Decode(&payload)
	if err != nil {
		log.Panic(err)
	}
	if payload.PHBType == BLOCK_TYPE {
		blockHash := payload.PHBItems[0]
		phbsendGetData(payload.PHBAddrFrom, BLOCK_TYPE, blockHash)
		if len(payload.PHBItems) >= 1 {
			transactionArray = payload.PHBItems[1:]
		}
	}
	if payload.PHBType == TX_TYPE {
		txHash := payload.PHBItems[0]
		if memoryTxPool[hex.EncodeToString(txHash)] == nil {
			phbsendGetData(payload.PHBAddrFrom, TX_TYPE, txHash)
		}

	}

}
