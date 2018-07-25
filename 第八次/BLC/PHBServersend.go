package BLC

import (
	"bytes"
	"io"
	"log"
	"net"
)

//COMMAND_VERSION
func phbsendVersion(toAddress string, bc *PHBBlockchain) {
	bestHeight := bc.PHBGetBestHeight()
	payload := phbgobEncode(PHBVersion{NODE_VERSION, bestHeight, nodeAddress})
	request := append(phbcommandToBytes(COMMAND_VERSION), payload...)
	phbsendData(toAddress, request)
}

//COMMAND_GETBLOCKS
func phbsendGetBlocks(toAddress string) {
	payload := phbgobEncode(PHBGetBlocks{nodeAddress})
	request := append(phbcommandToBytes(COMMAND_GETBLOCKS), payload...)
	phbsendData(toAddress, request)
}

func phbsendInv(toAddress string, kind string, hashes [][]byte) {
	payload := phbgobEncode(PHBInv{nodeAddress, kind, hashes})
	request := append(phbcommandToBytes(COMMAND_INV), payload...)
	phbsendData(toAddress, request)
}

func phbsendGetData(toAddress string, kind string, blockHash []byte) {
	payload := phbgobEncode(PHBGetData{nodeAddress, kind, blockHash})
	request := append(phbcommandToBytes(COMMAND_GETDATA), payload...)
	phbsendData(toAddress, request)
}

func phbsendBlock(toAddress string, block []byte) {
	payload := phbgobEncode(PHBBlockData{nodeAddress, block})
	request := append(phbcommandToBytes(COMMAND_BLOCK), payload...)
	phbsendData(toAddress, request)
}

func phbsendTx(toAddress string, tx *PHBTransaction) {
	payload := phbgobEncode(PHBTx{nodeAddress, tx})
	request := append(phbcommandToBytes(COMMAND_TX), payload...)
	phbsendData(toAddress, request)

}

func phbsendData(to string, data []byte) {
	conn, err := net.Dial("tcp", to)
	if err != nil {
		panic("error")
	}
	defer conn.Close()
	_, err = io.Copy(conn, bytes.NewReader(data))
	if err != nil {
		log.Panic(err)
	}
}
