package main

import (
	"./BLC"
)

func main() {
	//bc := BLC.CreatBlockChain()

	//bc.AddBlock("Block1 : 100PHB")
	//bc.AddBlock("Block2 : 200PHB")

	//bc.VisitBlockChain()
	cli := BLC.CLI{}
	cli.Run()
}

