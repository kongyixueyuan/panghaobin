package BLC

import (
	"bytes"
	"encoding/hex"
	"fmt"
)

type PHBTXOutput struct {
	PHBValue         int64
	PHBRipemd160Hash []byte //用户名
}

func (txOutput *PHBTXOutput) PHBPrintInfo() {
	fmt.Printf("Value:%s\n", txOutput.PHBValue)
	fmt.Printf("Ripemd160Hash:%s\n", hex.EncodeToString(txOutput.PHBRipemd160Hash))

}

func (txOutput *PHBTXOutput) PHBLock(address string) {
	publicKeyHash := PHBBase58Decode([]byte(address))
	txOutput.PHBRipemd160Hash = publicKeyHash[1 : len(publicKeyHash)-4]
}

func PHBNewTXOutput(value int64, address string) *PHBTXOutput {
	txOutput := &PHBTXOutput{value, nil}
	txOutput.PHBLock(address)
	return txOutput
}

//解锁
func (txOutput *PHBTXOutput) PHBUnLockScriptPubKeyWithAddress(address string) bool {
	publicKeyHash := PHBBase58Decode([]byte(address))
	hash160 := publicKeyHash[1 : len(publicKeyHash)-4]
	return bytes.Compare(txOutput.PHBRipemd160Hash, hash160) == 0
}
