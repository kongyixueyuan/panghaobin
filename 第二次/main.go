package main

import (
	"./BLC"
)

func main() {
	bc := BLC.NewBlockChain()

	bc.AddBlock("Block1 : 100PHB")
	bc.AddBlock("Block2 : 200PHB")


}

