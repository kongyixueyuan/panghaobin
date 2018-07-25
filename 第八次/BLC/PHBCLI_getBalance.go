package BLC

import "fmt"

func (cli *PHBCLI) phbgetBalance(address string, nodeID string) {
	fmt.Println("地址：" + address)
	// 获取某一个节点的blockchain对象
	blockchain := PHBBlockchainObject(nodeID)
	defer blockchain.PHBDB.Close()
	utxoSet := &PHBUTXOSet{blockchain}
	amount := utxoSet.PHBGetBalance(address)
	fmt.Printf("%s一共有%d个Token\n", address, amount)

}
