package main

import "fmt"

func (cli *PHBCLI) phbreindexUTXO(nodeID string) {
	bc := PHBNewBlockchain(nodeID)
	UTXOSet := PHBUTXOSet{bc}
	UTXOSet.PHBReindex()

	count := UTXOSet.PHBCountTransactions()
	fmt.Printf("Done! There are %d transactions in the UTXO set.\n", count)
}