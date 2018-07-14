package main

import (
	"log"
	"fmt"
)

func (cli *PHBCLI) phbsend(from, to string, amount int, nodeID string, mineNow bool) {
	if !PHBValidateAddress(from) {
		log.Panic("ERROR: Sender address is not valid")
	}
	if !PHBValidateAddress(to) {
		log.Panic("ERROR: Recipient address is not valid")
	}

	bc := PHBNewBlockchain(nodeID)
	UTXOSet := PHBUTXOSet{bc}
	defer bc.phbdb.Close()

	wallets, err := PHBNewWallets(nodeID)
	if err != nil {
		log.Panic(err)
	}
	wallet := wallets.PHBGetWallet(from)

	tx := PHBNewUTXOTransaction(&wallet, to, amount, &UTXOSet)

	if mineNow {
		cbTx := PHBNewCoinbaseTX(from, "")
		txs := []*PHBTransaction{cbTx, tx}

		newBlock := bc.PHBMineBlock(txs)
		UTXOSet.PHBUpdate(newBlock)
	} else {
		phbsendTx(knownNodes[0], tx)
	}

	fmt.Println("Success!")
}