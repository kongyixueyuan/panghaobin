package BLC

import (
	"fmt"
	"flag"
	"os"
	"log"
	"strings"
)

type CLI struct {
	//BC *BlockChain
}

func PrintfUsage()  {
	fmt.Println("Usage:")
	fmt.Println("  creatblockchain -- 创建包含创世区块的区块链 ")
	fmt.Println("  addblock -data -- 添加区块到区块链")
	fmt.Println("  visitblockchain -- 遍历输出区块链信息")

}

func isValidArgs()  {
	if len(os.Args) < 2 {
		PrintfUsage()
		os.Exit(1)
	}
}

func (cli *CLI) createGenesisBlockchain()  {
	CreatBlockChain()
}

func (cli *CLI) addBlock(data string)  {
	if DBISExist() == false {
		fmt.Println("数据库不存在!")
		os.Exit(1)
	}
	bc := BlockchainObject()
	defer bc.db.Close()
	bc.AddBlock(data)
}

func (cli *CLI) visitBlockChain()  {
	if DBISExist() == false {
		fmt.Println("数据库不存在!")
		os.Exit(1)
	}
	bc := BlockchainObject()
	defer bc.db.Close()
	bc.VisitBlockChain()
}

func (cli *CLI) Run()  {
	isValidArgs()
	creatBlockChainCmd := flag.NewFlagSet("creatblockchain", flag.ExitOnError)
	addBlockCmd := flag.NewFlagSet("addblock", flag.ExitOnError)
	visitBlockChainCmd := flag.NewFlagSet("visitblockchain", flag.ExitOnError)

	addBlockWithData := addBlockCmd.String("data", strings.Join(os.Args[2:], ""), "区块的数据")
	switch os.Args[1] {
	case "creatblockchain":
		err := creatBlockChainCmd.Parse(os.Args[2:])
		if err != nil {
			log.Panic(err)
		}
	case "addblock":
		err := addBlockCmd.Parse(os.Args[2:])
		if err != nil {
			log.Panic(err)
		}
	case "visitblockchain":
		err := visitBlockChainCmd.Parse(os.Args[2:])
		if err != nil {
			log.Panic(err)
		}
	default:
		PrintfUsage()
		os.Exit(1)
	}
	if creatBlockChainCmd.Parsed() {
		cli.createGenesisBlockchain()
	}

	if addBlockCmd.Parsed() {
		if *addBlockWithData == "" {
			addBlockCmd.Usage()
			os.Exit(1)
		}
		cli.addBlock(*addBlockWithData)
	}

	if visitBlockChainCmd.Parsed() {
		fmt.Println("输出区块所有信息：")
		cli.visitBlockChain()
	}
}
