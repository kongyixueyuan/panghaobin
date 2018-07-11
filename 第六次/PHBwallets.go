package main

import (
	"io/ioutil"
	"log"
	"encoding/gob"
	"crypto/elliptic"
	"bytes"
	"fmt"
	"os"
)

const walletFile = "wallet_%s.dat"

// Wallets stores a collection of wallets
type PHBWallets struct {
	PHBWallets map[string]*PHBWallet
}

// NewWallets creates Wallets and fills it from a file if it exists
func PHBNewWallets(nodeID string) (*PHBWallets, error) {
	wallets := PHBWallets{}
	wallets.PHBWallets = make(map[string]*PHBWallet)

	err := wallets.PHBLoadFromFile(nodeID)

	return &wallets, err
}

// CreateWallet adds a Wallet to Wallets
func (ws *PHBWallets) PHBCreateWallet() string {
	wallet := PHBNewWallet()
	address := fmt.Sprintf("%s", wallet.PHBGetAddress())

	ws.PHBWallets[address] = wallet

	return address
}

// GetAddresses returns an array of addresses stored in the wallet file
func (ws *PHBWallets) PHBGetAddresses() []string {
	var addresses []string

	for address := range ws.PHBWallets {
		addresses = append(addresses, address)
	}

	return addresses
}

// GetWallet returns a Wallet by its address
func (ws *PHBWallets) PHBGetWallet(address string) PHBWallet {
	return *ws.PHBWallets[address]
}

// LoadFromFile loads wallets from the file
func (ws *PHBWallets) PHBLoadFromFile(nodeID string) error {
	walletFile := fmt.Sprintf(walletFile, nodeID)
	if _, err := os.Stat(walletFile); os.IsNotExist(err) {
		return err
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

	ws.PHBWallets = wallets.PHBWallets

	return nil
}

// SaveToFile saves wallets to a file
func (ws *PHBWallets) PHBSaveToFile(nodeID string) {
	var content bytes.Buffer
	walletFile := fmt.Sprintf(walletFile, nodeID)

	gob.Register(elliptic.P256())

	encoder := gob.NewEncoder(&content)
	err := encoder.Encode(ws)
	if err != nil {
		log.Panic(err)
	}

	err = ioutil.WriteFile(walletFile, content.Bytes(), 0644)
	if err != nil {
		log.Panic(err)
	}
}