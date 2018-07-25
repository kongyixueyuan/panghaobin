package BLC

func (cli *PHBCLI) phbcreateGenesisBlockChain(address string, nodeID string) {
	blockchain := PHBCreateBlockchainWithGenesisBlock(address, nodeID)
	defer blockchain.PHBDB.Close()

	utxoSet := &PHBUTXOSet{blockchain}

	utxoSet.PHBResetUTXOSet()
}
