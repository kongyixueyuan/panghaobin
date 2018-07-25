package BLC

type PHBInv struct {
	PHBAddrFrom string   //自己的地址
	PHBType     string   //类型 block tx
	PHBItems    [][]byte //hash二维数组
}
