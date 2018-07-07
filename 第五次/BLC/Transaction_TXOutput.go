package BLC

import (
	"fmt"
)

type TxOutput struct {
	Value int64
	PubKey string
}

func (txOutput *TxOutput) UnLockPubKeyWithAddress(address string) bool {
	return txOutput.PubKey == address
}

func (txOutput *TxOutput) printfTxOutput()  {
	fmt.Printf("========================================\n")
	fmt.Printf("交易输出：%d\n", txOutput.Value)
	fmt.Printf("交易输出PubKey：%x\n", txOutput.PubKey)
	fmt.Printf("========================================\n\n")
}