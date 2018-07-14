package main

import (
	"fmt"
	"strconv"
)

func (cli *PHBCLI) phbprintChain(nodeID string) {
	bc := PHBNewBlockchain(nodeID)
	defer bc.phbdb.Close()

	bci := bc.PHBIterator()

	for {
		block := bci.PHBNext()

		fmt.Printf("============ Block %x ============\n", block.PHBHash)
		fmt.Printf("Height: %d\n", block.PHBHeight)
		fmt.Printf("Prev. block: %x\n", block.PHBPrevBlockHash)
		pow := PHBNewProofOfWork(block)
		fmt.Printf("PoW: %s\n\n", strconv.FormatBool(pow.PHBValidate()))
		for _, tx := range block.PHBTransactions {
			fmt.Println(tx)
		}
		fmt.Printf("\n\n")

		if len(block.PHBPrevBlockHash) == 0 {
			break
		}
	}
}