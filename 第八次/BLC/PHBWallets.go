package BLC

import (
	"bytes"
	"crypto/elliptic"
	"encoding/gob"
	"fmt"
	"io/ioutil"
	"log"
	"os"
)

const walletFile = "Wallets_%s.dat"

type PHBWallets struct {
	PHBWalletsMap map[string]*PHBWallet
}

// 创建钱包集合
func PHBNewWallets(nodeID string) (*PHBWallets, error) {
	walletFile := fmt.Sprintf(walletFile, nodeID)
	if _, err := os.Stat(walletFile); os.IsNotExist(err) {
		wallets := &PHBWallets{}
		wallets.PHBWalletsMap = make(map[string]*PHBWallet)
		return wallets, err
	}
	fileContent, err := ioutil.ReadFile(walletFile)
	if err != nil {
		log.Panic(err)
	}
	var wallets PHBWallets
	gob.Register(elliptic.P256())
	decoder := gob.NewDecoder(bytes.NewReader(fileContent))
	err = decoder.Decode(&wallets)
	if err != nil {
		log.Panic(err)
	}
	return &wallets, nil
}

func (w *PHBWallets) PHBCreateNewWallet(nodeID string) {
	wallet := PHBNewWallet()
	fmt.Printf("Address：%s\n", wallet.PHBGetAddress())
	w.PHBWalletsMap[string(wallet.PHBGetAddress())] = wallet
	w.PHBSaveWallets(nodeID)
}

func (w *PHBWallets) PHBSaveWallets(nodeID string) {
	walletFile := fmt.Sprintf(walletFile, nodeID)
	var content bytes.Buffer
	gob.Register(elliptic.P256())
	encoder := gob.NewEncoder(&content)
	err := encoder.Encode(&w)
	if err != nil {
		log.Panic(err)
	}
	err = ioutil.WriteFile(walletFile, content.Bytes(), 0644)
	if err != nil {
		log.Panic(err)
	}
}
