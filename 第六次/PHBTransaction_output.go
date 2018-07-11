package main

import (
	"bytes"
	"encoding/gob"
	"log"
)

type PHBTXOutput struct {
	PHBValue int
	PHBPubKeyHash []byte
}

type PHBTXOutputs struct {
	PHBOutPuts []PHBTXOutput
}

func (out *PHBTXOutput) PHBLock(address []byte)  {
	pubKeyHash := PHBBase58Decode(address)
	pubKeyHash = pubKeyHash[1:len(pubKeyHash)-4]
	out.PHBPubKeyHash = pubKeyHash
}

func (out *PHBTXOutput) PHBIsLockedWithKey(pubKeyHash []byte) bool {
	return bytes.Compare(out.PHBPubKeyHash, pubKeyHash) == 0
}

func PHBNewTXOutPut(value int, address string) *PHBTXOutput {
	txo := &PHBTXOutput{value, nil}
	txo.PHBLock([]byte(address))
	return txo
}

func (outs *PHBTXOutputs) PHBSerialize() []byte {
	var buff bytes.Buffer
	enc := gob.NewEncoder(&buff)
	err := enc.Encode(outs)
	if err != nil {
		log.Panic(err)
	}
	return buff.Bytes()
}

func PHBDeserializeOutputs(data []byte) PHBTXOutputs {
	var outputs PHBTXOutputs
	dec := gob.NewDecoder(bytes.NewReader(data))
	err := dec.Decode(&outputs)
	if err != nil {
		log.Panic(err)
	}
	return outputs
}


