package main

import "fmt"

func (cli *PHBCLI) phbcreateWallet(nodeID string) {
	wallets, _ := PHBNewWallets(nodeID)
	address := wallets.PHBCreateWallet()
	wallets.PHBSaveToFile(nodeID)

	fmt.Printf("Your new address: %s\n", address)
}