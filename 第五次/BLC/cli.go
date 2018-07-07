package BLC

import (
	"fmt"
	"flag"
	"os"
	"log"
	"strconv"
)

type CLI struct {
	//BC *BlockChain
}

func PrintfUsage()  {
	fmt.Println("Usage:")
	fmt.Println("  creatblockchain -address ADDRESS -tokens VALUE-- 创建包含创世区块的区块链并可以选择设定初始token ./main creatblockchain -address pang -tokens 101")
	//fmt.Println("  addblock -data -- 添加区块到区块链")
	fmt.Println("  send -from From -to To -amount AMOUNT -- 交易明细./main send -from '[\"pang\"]' -to '[\"lisi\"]' -amount '[\"2.5\"]")
	fmt.Println("  getbalance -address -- 获取区块信息")
	fmt.Println("  visitblockchain -- 遍历输出区块链信息")
	fmt.Println("  createwallet -- 创建钱包")
	fmt.Println("  printfwalletaddresslist -- 输出所有钱包地址")

}

func isValidArgs()  {
	if len(os.Args) < 2 {
		PrintfUsage()
		os.Exit(1)
	}
}

func (cli *CLI) createGenesisBlockchain(address string, values string)  {
	val, err := strconv.Atoi(values)
	if err != nil {
		CreatBlockChain(address, 100)
	} else {
		if val < 1 {
			val = 100
		}
		CreatBlockChain(address, int64(val) )
	}

}

// 转账
func (cli *CLI) send(from []string,to []string,amount []string)  {
	if !DBISExist() {
		fmt.Println("数据库不存在")
		os.Exit(1)
	}
	bc := BlockchainObject()
	defer bc.db.Close()
	bc.MineNewBlock(from,to,amount)
}

func (cli *CLI) getBalance(address string)  {
	fmt.Println("地址bc：" + address)
	bc := BlockchainObject()
	defer bc.db.Close()
	amount := bc.GetBalance(address)
	fmt.Printf("%s一共有%d个Token\n",address,amount)
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

func (cli *CLI) createwallet()  {
	wallets, _ := NewWallets()
	wallets.CreateWallet()
	fmt.Println(len(wallets.WalletsMap))
}

func (cli *CLI) printfwalletaddresslist()  {
	fmt.Println("所有钱包地址：")
	wallets, _ := NewWallets()
	for address, _ := range wallets.WalletsMap {
		fmt.Println(address)
	}
}

func (cli *CLI) Run()  {
	isValidArgs()
	creatBlockChainCmd := flag.NewFlagSet("creatblockchain", flag.ExitOnError)
	sendBlockCmd := flag.NewFlagSet("send", flag.ExitOnError)
	visitBlockChainCmd := flag.NewFlagSet("visitblockchain", flag.ExitOnError)
	getbalanceCmd := flag.NewFlagSet("getBalance", flag.ExitOnError)
	createwalletCmd := flag.NewFlagSet("createwallet", flag.ExitOnError)
	printfwalletaddresslistCmd := flag.NewFlagSet("printfwalletaddresslist", flag.ExitOnError)

	creatBlockChainAddress := creatBlockChainCmd.String("address",  "", "创建创世区块地址")
	creatBlockTokens := creatBlockChainCmd.String("tokens",  "", "设定tokens")

	//	$ ./main send -from '["pang"]' -to '["lisi"]' -amount '["2.5"]'
	flagFrom := sendBlockCmd.String("from","","转账源地址......")
	flagTo := sendBlockCmd.String("to","","转账目的地地址......")
	flagAmount := sendBlockCmd.String("amount","","转账金额......")

	getbalanceAdress := getbalanceCmd.String("address", "", "要查询某一个账号的余额")

	//addBlockWithData := addBlockCmd.String("data", strings.Join(os.Args[2:], ""), "区块的数据")

	
	

	switch os.Args[1] {
	case "creatblockchain":
		err := creatBlockChainCmd.Parse(os.Args[2:])
		if err != nil {
			log.Panic(err)
		}
	case "send":
		err := sendBlockCmd.Parse(os.Args[2:])
		if err != nil {
			log.Panic(err)
		}
	case "getbalance":
		err := getbalanceCmd.Parse(os.Args[2:])
		if err != nil {
			log.Panic(err)
		}
	case "visitblockchain":
		err := visitBlockChainCmd.Parse(os.Args[2:])
		if err != nil {
			log.Panic(err)
		}

	case "createwallet":
		err := createwalletCmd.Parse(os.Args[2:])
		if err != nil {
			log.Panic(err)
		}

	case "printfwalletaddresslist":
		err := printfwalletaddresslistCmd.Parse(os.Args[2:])
		if err != nil {
			log.Panic(err)
		}

	default:
		PrintfUsage()
		os.Exit(1)
	}
	if creatBlockChainCmd.Parsed() {
		if *creatBlockChainAddress == "" || *creatBlockTokens == "" {
			PrintfUsage()
			os.Exit(1)
		}
		if ValidateAddress([]byte(*creatBlockChainAddress)) == false {
			fmt.Println("地址无效")
			PrintfUsage()
			os.Exit(1)
		}
		cli.createGenesisBlockchain(*creatBlockChainAddress, *creatBlockTokens)
	}

	if sendBlockCmd.Parsed() {
		if *flagFrom == "" || *flagTo == "" || *flagAmount == ""{
			PrintfUsage()
			os.Exit(1)
		}
		from := JSONToArray(*flagFrom)
		to := JSONToArray(*flagTo)
		amount := JSONToArray(*flagAmount)
		cli.send(from,to,amount)
	}
	if getbalanceCmd.Parsed() {
		if *getbalanceAdress == "" {
			fmt.Println("地址不能为空")
			PrintfUsage()
			os.Exit(1)
		}
		cli.getBalance(*getbalanceAdress)
	}
	if visitBlockChainCmd.Parsed() {
		fmt.Println("输出区块所有信息：")
		cli.visitBlockChain()
	}
	if createwalletCmd.Parsed() {
		cli.createwallet()
	}
	if printfwalletaddresslistCmd.Parsed() {
		cli.printfwalletaddresslist()
	}
}
