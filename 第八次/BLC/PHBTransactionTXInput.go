package BLC

import (
	"bytes"
	"encoding/hex"
	"fmt"
)

type PHBTXInput struct {
	PHBTxHash    []byte //交易的Hash
	PHBVout      int    //存储TXOutput在Vout里面的索引
	PHBSignature []byte //数字签名
	PHBPublicKey []byte //公钥，钱包里面
}

func (txInput *PHBTXInput) PHBPrintInfo() {
	fmt.Printf("txHash:%s\n", hex.EncodeToString(txInput.PHBTxHash))
	fmt.Printf("Vout:%d\n", txInput.PHBVout)
	fmt.Printf("Signature:%x\n", txInput.PHBSignature)
	fmt.Printf("PublicKey:%x\n", txInput.PHBPublicKey)
}

//判断当前的消费是谁的钱
func (txInput *PHBTXInput) PHBUnLockRipemd160Hash(ripemd160Hash []byte) bool {
	publicKey := PHBRipemd160Hash(txInput.PHBPublicKey)
	return bytes.Compare(publicKey, ripemd160Hash) == 0
}
