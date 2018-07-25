package BLC

func (cli *PHBCLI) phbresetUTXOSet(nodeID string) {
	blockchain := PHBBlockchainObject(nodeID)
	defer blockchain.PHBDB.Close()
	utxoSet := &PHBUTXOSet{blockchain}
	utxoSet.PHBResetUTXOSet()

}
