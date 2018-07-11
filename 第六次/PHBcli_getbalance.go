package main

import (
	"log"
	"fmt"
)

func (cli *PHBCLI) phbgetBalance(address, nodeID string) {
	if !PHBValidateAddress(address) {
		log.Panic("ERROR: Address is not valid")
	}
	bc := PHBNewBlockchain(nodeID)
	UTXOSet := PHBUTXOSet{bc}
	defer bc.phbdb.Close()

	balance := 0
	pubKeyHash := PHBBase58Decode([]byte(address))
	pubKeyHash = pubKeyHash[1 : len(pubKeyHash)-4]
	UTXOs := UTXOSet.PHBFindUTXO(pubKeyHash)

	for _, out := range UTXOs {
		balance += out.PHBValue
	}

	fmt.Printf("Balance of '%s': %d\n", address, balance)
}