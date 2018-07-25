package BLC

import (
	"bytes"
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/sha256"
	"encoding/gob"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"log"
	"math/big"
	"time"
)

type PHBTransaction struct {
	//1. 交易hash
	PHBTxHash []byte
	//2. 输入
	PHBVins []*PHBTXInput
	//3. 输出
	PHBVouts []*PHBTXOutput
}

func (tx *PHBTransaction) PHBPrintTX() {
	fmt.Printf("txHash : %s\n", hex.EncodeToString(tx.PHBTxHash))
	fmt.Println("Vins====")
	for _, vin := range tx.PHBVins {
		vin.PHBPrintInfo()
	}
	fmt.Println("Vouts====")
	for _, vout := range tx.PHBVouts {
		vout.PHBPrintInfo()
	}
	fmt.Println("============================")
}

// 判断当前的交易是否是Coinbase交易
func (tx *PHBTransaction) PHBIsCoinbaseTransaction() bool {
	return len(tx.PHBVins[0].PHBTxHash) == 0 && tx.PHBVins[0].PHBVout == -1
}

func PHBNewCoinbaseTransaction(address string) *PHBTransaction {
	txInput := &PHBTXInput{[]byte{}, -1, nil, []byte{}}
	txOutput := PHBNewTXOutput(10, address)
	txCoinbase := &PHBTransaction{[]byte{}, []*PHBTXInput{txInput}, []*PHBTXOutput{txOutput}}
	//设置hash值
	txCoinbase.PHBHashTransaction()
	return txCoinbase
}

func (tx *PHBTransaction) PHBHashTransaction() {
	var result bytes.Buffer
	encoder := gob.NewEncoder(&result)
	err := encoder.Encode(tx)
	if err != nil {
		log.Panic(err)
	}
	resultBytes := bytes.Join([][]byte{PHBIntToHex(time.Now().Unix()), result.Bytes()}, []byte{})
	hash := sha256.Sum256(resultBytes)
	tx.PHBTxHash = hash[:]
}

//2. 转账时产生的Transaction
func PHBNewSimpleTransaction(from string, to string, amount int64, utxoSet *PHBUTXOSet, txs []*PHBTransaction, nodeID string) *PHBTransaction {
	wallets, _ := PHBNewWallets(nodeID)
	wallet := wallets.PHBWalletsMap[from]
	money, spendableUTXODic := utxoSet.PHBFindSpendableUTXOS(from, amount, txs)
	var txIntputs []*PHBTXInput
	var txOutputs []*PHBTXOutput
	for txHash, indexArray := range spendableUTXODic {
		txHashBytes, _ := hex.DecodeString(txHash)
		for _, index := range indexArray {
			txInput := &PHBTXInput{txHashBytes, index, nil, wallet.PHBPublicKey}
			txIntputs = append(txIntputs, txInput)
		}
	}
	//转账
	txOutput := PHBNewTXOutput(int64(amount), to)
	txOutputs = append(txOutputs, txOutput)
	fmt.Println("找零:", int64(money), "-", int64(amount), ":", int64(money)-int64(amount))
	//找零
	txOutput = PHBNewTXOutput(int64(money)-int64(amount), from)
	txOutputs = append(txOutputs, txOutput)
	tx := &PHBTransaction{[]byte{}, txIntputs, txOutputs}
	//设置hash值
	tx.PHBHashTransaction()
	//进行签名
	utxoSet.PHBBlockchain.PHBSignTransaction(tx, wallet.PHBPrivateKey, txs)
	tx.PHBPrintTX()
	return tx
}

func (tx *PHBTransaction) PHBHash() []byte {
	txCopy := tx
	txCopy.PHBTxHash = []byte{}
	hash := sha256.Sum256(txCopy.PHBSerialize())
	return hash[:]
}

func (tx *PHBTransaction) PHBSerialize() []byte {
	jsonByte, err := json.Marshal(tx)
	if err != nil {
		log.Panic(err)
	}
	return jsonByte
}

func (tx *PHBTransaction) PHBSign(privKey ecdsa.PrivateKey, prevTXs map[string]PHBTransaction) {
	if tx.PHBIsCoinbaseTransaction() {
		return
	}
	for _, vin := range tx.PHBVins {
		if prevTXs[hex.EncodeToString(vin.PHBTxHash)].PHBTxHash == nil {
			log.Panic("ERROR: Previous transaction is not correct")
		}
	}
	txCopy := tx.PHBTrimmedCopy()
	fmt.Println("签名前")
	for inID, vin := range txCopy.PHBVins {
		fmt.Println("签名中")
		prevTx := prevTXs[hex.EncodeToString(vin.PHBTxHash)]
		txCopy.PHBVins[inID].PHBSignature = nil
		txCopy.PHBVins[inID].PHBPublicKey = prevTx.PHBVouts[vin.PHBVout].PHBRipemd160Hash
		txCopy.PHBTxHash = txCopy.PHBHash()
		txCopy.PHBVins[inID].PHBPublicKey = nil
		r, s, err := ecdsa.Sign(rand.Reader, &privKey, txCopy.PHBTxHash)
		if err != nil {
			log.Panic(err)
		}
		signature := append(r.Bytes(), s.Bytes()...)
		tx.PHBVins[inID].PHBSignature = signature
	}
}

func (tx *PHBTransaction) PHBTrimmedCopy() PHBTransaction {
	var inputs []*PHBTXInput
	var outputs []*PHBTXOutput

	for _, vin := range tx.PHBVins {
		inputs = append(inputs, &PHBTXInput{vin.PHBTxHash, vin.PHBVout, nil, nil})
	}

	for _, vout := range tx.PHBVouts {
		outputs = append(outputs, &PHBTXOutput{vout.PHBValue, vout.PHBRipemd160Hash})
	}
	txCopy := PHBTransaction{tx.PHBTxHash, inputs, outputs}
	return txCopy
}

func (tx *PHBTransaction) PHBVerify(prevTXs map[string]PHBTransaction) bool {
	if tx.PHBIsCoinbaseTransaction() {
		return true
	}

	for _, vin := range tx.PHBVins {
		if prevTXs[hex.EncodeToString(vin.PHBTxHash)].PHBTxHash == nil {
			log.Panic("ERROR: Previous transaction is not correct")
		}
	}
	txCopy := tx.PHBTrimmedCopy()
	curve := elliptic.P256()
	for inID, vin := range tx.PHBVins {
		prevTx := prevTXs[hex.EncodeToString(vin.PHBTxHash)]
		txCopy.PHBVins[inID].PHBSignature = nil
		txCopy.PHBVins[inID].PHBPublicKey = prevTx.PHBVouts[vin.PHBVout].PHBRipemd160Hash
		txCopy.PHBTxHash = txCopy.PHBHash()
		txCopy.PHBVins[inID].PHBPublicKey = nil
		// 私钥 ID
		r := big.Int{}
		s := big.Int{}
		sigLen := len(vin.PHBSignature)
		r.SetBytes(vin.PHBSignature[:(sigLen / 2)])
		s.SetBytes(vin.PHBSignature[(sigLen / 2):])
		x := big.Int{}
		y := big.Int{}
		keyLen := len(vin.PHBPublicKey)
		x.SetBytes(vin.PHBPublicKey[:(keyLen / 2)])
		y.SetBytes(vin.PHBPublicKey[(keyLen / 2):])
		rawPubKey := ecdsa.PublicKey{curve, &x, &y}
		if ecdsa.Verify(&rawPubKey, txCopy.PHBTxHash, &r, &s) == false {
			return false
		}
	}
	return true
}
