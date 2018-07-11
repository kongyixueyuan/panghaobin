package main

import (
	"fmt"
	"log"
)

func (cli *PHBCLI) phbstartNode(nodeID, minerAddress string) {
	fmt.Printf("Starting node %s\n", nodeID)
	if len(minerAddress) > 0 {
		if PHBValidateAddress(minerAddress) {
			fmt.Println("Mining is on. Address to receive rewards: ", minerAddress)
		} else {
			log.Panic("Wrong miner address!")
		}
	}
	PHBStartServer(nodeID, minerAddress)
}