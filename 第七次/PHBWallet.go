package main

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"log"
	"crypto/sha256"
	"golang.org/x/crypto/ripemd160"
	"bytes"
)
const version = byte(0x00)
const addressChecksumLen = 4

type PHBWallet struct {
	PHBPrivateKey ecdsa.PrivateKey
	PHBPublicKey []byte
}

func PHBNewWallet() *PHBWallet {
	private, public := phbNewPair()
	wallet := PHBWallet{private, public}
	return &wallet
}

func (w *PHBWallet) PHBGetAddress() []byte {
	pubKeyHash := PHBHashPubKey(w.PHBPublicKey)
	versionPayload := append([]byte{version}, pubKeyHash...)
	checksum := phbChecksum(versionPayload)
	fullPayload := append(versionPayload, checksum...)
	address := PHBBase58Encode(fullPayload)
	return address
}

func PHBValidateAddress(address string) bool {
	pubKeyHash := PHBBase58Decode([]byte(address))
	actualChecksum := pubKeyHash[len(pubKeyHash) - addressChecksumLen:]
	version := pubKeyHash[0]
	pubKeyHash = pubKeyHash[1:len(pubKeyHash)-addressChecksumLen]
	targetChecksum := phbChecksum(append([]byte{version}, pubKeyHash...))
	return bytes.Compare(actualChecksum, targetChecksum)==0

}

func PHBHashPubKey(pubKey []byte) []byte {
	publickSHA256 := sha256.Sum256(pubKey)
	RIPEMD160Hasher := ripemd160.New()
	_, err := RIPEMD160Hasher.Write(publickSHA256[:])
	if err != nil {
		log.Panic(err)
	}
	publicRIPEMD160 := RIPEMD160Hasher.Sum(nil)
	return  publicRIPEMD160
}

func phbChecksum(payload []byte) []byte {
	sha1 := sha256.Sum256(payload)
	sha2 := sha256.Sum256(sha1[:])
	return sha2[:addressChecksumLen]
}

func phbNewPair() (ecdsa.PrivateKey, []byte) {
	curve := elliptic.P256()
	private, err := ecdsa.GenerateKey(curve, rand.Reader)
	if err != nil {
		log.Panic(err)
	}
	pubKey := append(private.PublicKey.X.Bytes(), private.PublicKey.Y.Bytes()...)
	return *private, pubKey
}

