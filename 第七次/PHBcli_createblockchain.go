package main

import (
	"log"
	"fmt"
)

func (cli *PHBCLI) phbcreateBlockchain(address, nodeID string) {
	if !PHBValidateAddress(address) {
		log.Panic("ERROR: Address is not valid")
	}
	bc := PHBCreateBlockchain(address, nodeID)
	defer bc.phbdb.Close()

	UTXOSet := PHBUTXOSet{bc}
	UTXOSet.PHBReindex()

	fmt.Println("Done!")
}