package main

import (
	"log"
	"fmt"
)

func (cli *PHBCLI) phblistAddresses(nodeID string) {
	wallets, err := PHBNewWallets(nodeID)
	if err != nil {
		log.Panic(err)
	}
	addresses := wallets.PHBGetAddresses()
	for _, address := range addresses {
		fmt.Println(address)
	}
}