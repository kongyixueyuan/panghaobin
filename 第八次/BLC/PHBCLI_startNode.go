package BLC

import (
	"fmt"
	"os"
)

func (cli *PHBCLI) phbstartNode(nodeID string, minerAdd string) {
	if minerAdd == "" || PHBIsValidForAdress([]byte(minerAdd)) {
		//  启动服务器
		fmt.Printf("启动服务器:localhost:%s\n", nodeID)
		phbstartServer(nodeID, minerAdd)
	} else {
		fmt.Println("指定的地址无效....")
		os.Exit(0)
	}

}
