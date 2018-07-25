package BLC

import (
	"fmt"
	"io/ioutil"
	"log"
	"net"
)

func phbstartServer(nodeID string, minerAdd string) {

	nodeAddress = fmt.Sprintf("localhost:%s", nodeID)
	minerAddress = minerAdd
	ln, err := net.Listen(PROTOCOL, nodeAddress)
	if err != nil {
		log.Panic(err)
	}
	defer ln.Close()
	bc := PHBBlockchainObject(nodeID)
	defer bc.PHBDB.Close()
	if nodeAddress != knowNodes[0] {
		phbsendVersion(knowNodes[0], bc)
	}
	for {
		conn, err := ln.Accept()
		if err != nil {
			log.Panic(err)
		}
		go phbhandleConnection(conn, bc)
	}

}

func phbhandleConnection(conn net.Conn, bc *PHBBlockchain) {

	// 读取客户端发送过来的所有的数据
	request, err := ioutil.ReadAll(conn)
	if err != nil {
		log.Panic(err)
	}

	fmt.Printf("Receive a Message:%s\n", request[:COMMANDLENGTH])
	command := phbbytesToCommand(request[:COMMANDLENGTH])

	switch command {
	case COMMAND_VERSION:
		phbhandleVersion(request, bc)
	case COMMAND_ADDR:
		phbhandleAddr(request, bc)
	case COMMAND_BLOCK:
		phbhandleBlock(request, bc)
	case COMMAND_GETBLOCKS:
		phbhandleGetblocks(request, bc)
	case COMMAND_GETDATA:
		phbhandleGetData(request, bc)
	case COMMAND_INV:
		phbhandleInv(request, bc)
	case COMMAND_TX:
		phbhandleTx(request, bc)
	default:
		fmt.Println("Unknown command!")
	}
	conn.Close()
}

func phbnodeIsKnown(addr string) bool {
	for _, node := range knowNodes {
		if node == addr {
			return true
		}
	}
	return false
}
