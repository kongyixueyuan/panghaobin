package BLC

var knowNodes = []string{"localhost:3000"}
var nodeAddress string
var transactionArray [][]byte
var minerAddress string
var memoryTxPool = make(map[string]*PHBTransaction)
