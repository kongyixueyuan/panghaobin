package BLC

type TxInput struct {
	TxHash []byte //交易hash
	Vout int  //txoutput 索引
	Signature string //用户名
}

func (txInput *TxInput) IsUserKey(address string) bool  {
	return txInput.Signature == address
}
