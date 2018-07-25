package BLC

func (cli *PHBCLI) phbprintchain(nodeID string) {
	blockchain := PHBBlockchainObject(nodeID)
	defer blockchain.PHBDB.Close()
	blockchain.PHBPrintchain()

}
