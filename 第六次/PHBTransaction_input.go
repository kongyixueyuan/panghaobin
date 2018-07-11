package main

import "bytes"

type PHBTXInput struct {
	PHBTxid []byte
	PHBVout int 
	PHBSignature []byte 
	PHBPubKey []byte 
}

func (in *PHBTXInput) PHBUseKey(pubKeyHash []byte) bool {
	lockHash := PHBHashPubKey(in.PHBPubKey)
	return bytes.Compare(lockHash, pubKeyHash) == 0
}
