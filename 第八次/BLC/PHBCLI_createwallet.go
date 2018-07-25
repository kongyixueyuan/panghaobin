package BLC

import "fmt"

func (cli *PHBCLI) phbcreateWallet(nodeID string) {
	wallets, _ := PHBNewWallets(nodeID)
	wallets.PHBCreateNewWallet(nodeID)
	fmt.Println(len(wallets.PHBWalletsMap))
}
