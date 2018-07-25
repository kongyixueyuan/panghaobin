package BLC

import "fmt"

func (cli *PHBCLI) phbaddressLists(nodeID string) {

	fmt.Println("所有的钱包地址:")
	wallets, _ := PHBNewWallets(nodeID)
	for address, _ := range wallets.PHBWalletsMap {
		fmt.Println(address)
	}
}
