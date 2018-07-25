package BLC

import (
	"fmt"
	"strconv"
)

func (cli *PHBCLI) phbsend(from []string, to []string, amount []string, nodeID string, mineNow bool) {
	blockchain := PHBBlockchainObject(nodeID)
	utxoSet := &PHBUTXOSet{blockchain}
	defer blockchain.PHBDB.Close()
	if mineNow {
		blockchain.PHBMineNewBlock(from, to, amount, nodeID)
		utxoSet.PHBUpdate()
	} else {
		// 把交易发送到矿工节点去进行验证
		fmt.Println("由矿工节点处理......")
		value, _ := strconv.Atoi(amount[0])
		tx := PHBNewSimpleTransaction(from[0], to[0], int64(value), utxoSet, []*PHBTransaction{}, nodeID)
		phbsendTx(knowNodes[0], tx)
	}

}
