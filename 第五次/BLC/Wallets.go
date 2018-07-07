package BLC

import (
	"os"
	"io/ioutil"
	"log"
	"encoding/gob"
	"crypto/elliptic"
	"bytes"
	"fmt"
)

const walletFile = "wallets.dat"

type Wallets struct {
	WalletsMap map[string]*Wallet
}

func NewWallets() (*Wallets, error) {
	if _, err := os.Stat(walletFile); os.IsNotExist(err) {
		wallets := &Wallets{}
		wallets.WalletsMap = make(map[string]*Wallet)
		return wallets, err
	}
	fileContent, err := ioutil.ReadFile(walletFile)
	if err != nil {
		log.Panic(err)
	}
	var wallets Wallets
	gob.Register(elliptic.P256())
	decoder := gob.NewDecoder(bytes.NewReader(fileContent))
	err = decoder.Decode(&wallets)
	if err != nil {
		log.Panic(err)
	}
	return &wallets, nil

}

//写入钱包信息
func (w *Wallets) SaveWallets()  {
	var content bytes.Buffer
	gob.Register(elliptic.P256())

	encoder := gob.NewEncoder(&content)
	err := encoder.Encode(&w)
	if err != nil {
		log.Panic(err)
	}

	err = ioutil.WriteFile(walletFile, content.Bytes(), 0644)//写入文件 覆盖原文件
	if err != nil {
		log.Panic(err)
	}
}

//创建一个钱包
func (w *Wallets) CreateWallet()  {
	wallet := NewWallet()
	fmt.Printf("钱包地址： %s \n", wallet.GetAddress())
	w.WalletsMap[string(wallet.GetAddress())] = wallet
	w.SaveWallets()
}