package BLC

import (
	"bytes"
	"encoding/gob"
	"log"
	"crypto/sha256"
	"encoding/hex"
)

type Transaction struct {
	TxHash []byte //交易hash
	Vin []*TxInput
	Vout []*TxOutput
}

func (tx *Transaction) TxHashTransaction()  {
	var hashbuf bytes.Buffer
	encoder := gob.NewEncoder(&hashbuf)
	err := encoder.Encode(tx)
	if err != nil {
		log.Panic(err)
	}
	hash := sha256.Sum256(hashbuf.Bytes())
	tx.TxHash = hash[:]
}

func (tx *Transaction) IsCoinbaseTransaction() bool {//是否Coinbase交易
	return len(tx.Vin[0].TxHash) == 0 && tx.Vin[0].Vout == -1
}

//1 创世区块创建时的Transaction
func CreateCoinbaseTransaction(address string, values int64) *Transaction  {
	txin := &TxInput{[]byte{}, -1, "Genesis Data"}
	txout := &TxOutput{values, address}

	txCoinbase := &Transaction{[]byte{}, []*TxInput{txin}, []*TxOutput{txout}}

	txCoinbase.TxHashTransaction()
	return txCoinbase
}

//2 转账时的transaction
func CreateSimpleTransaction(from string, to string, amount int, bc *BlockChain, txs []*Transaction) *Transaction {

	money,spendableUTXODic := bc.FindSpendableUTXOs(from, amount, txs)
	var txIntputs []*TxInput
	var txOutputs []*TxOutput
	for txHash,indexArray := range spendableUTXODic  {
		txHashBytes,_ := hex.DecodeString(txHash)
		for _,index := range indexArray  {
			txInput := &TxInput{txHashBytes,index,from}
			txIntputs = append(txIntputs,txInput)
		}
	}
	// 转账
	txOutput := &TxOutput{int64(amount), to}
	txOutputs = append(txOutputs,txOutput)
	// 找零
	txOutput = &TxOutput{int64(money) - int64(amount), from}
	txOutputs = append(txOutputs,txOutput)
	tx := &Transaction{[]byte{}, txIntputs, txOutputs}
	//设置hash值
	tx.TxHashTransaction()
	return tx
}






