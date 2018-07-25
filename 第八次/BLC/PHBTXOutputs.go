package BLC

import (
	"bytes"
	"encoding/gob"
	"log"
)

type PHBTXOutputs struct {
	PHBUTXOS []*PHBUTXO
}

// 将区块序列化成字节数组
func (txOutputs *PHBTXOutputs) PHBSerialize() []byte {
	var result bytes.Buffer
	encoder := gob.NewEncoder(&result)
	err := encoder.Encode(txOutputs)
	if err != nil {
		log.Panic(err)
	}
	return result.Bytes()
}

// 反序列化
func PHBDeserializeTXOutputs(txOutputsBytes []byte) *PHBTXOutputs {
	var txOutputs PHBTXOutputs
	decoder := gob.NewDecoder(bytes.NewReader(txOutputsBytes))
	err := decoder.Decode(&txOutputs)
	if err != nil {
		log.Panic(err)
	}
	return &txOutputs
}
