package BLC

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"log"
	"crypto/sha256"
	"golang.org/x/crypto/ripemd160"
	"fmt"
	"bytes"
)

const version = byte(0x00)
const addressChechsumLen = 4

type Wallet struct {
	PrivateKey ecdsa.PrivateKey //私钥
	PublickKey []byte  //公钥
}

func NewWallet() *Wallet {
	private, public := newKeyPair()
	wallet := Wallet{private, public}
	return &wallet
}

func newKeyPair() (ecdsa.PrivateKey, []byte) {
	curve := elliptic.P256()
	private, err := ecdsa.GenerateKey(curve, rand.Reader)
	if err != nil {
		log.Panic(err)
	}
	pubKey := append(private.PublicKey.X.Bytes(), private.PublicKey.X.Bytes()...)
	return *private, pubKey
}

func (w Wallet) GetAddress() []byte {
	ripend160Hash := Ripemd160Hash(w.PublickKey)//hash160 20bytes

	version_ripend160hash := append([]byte{version}, ripend160Hash...)//21字节
	checksum := Checksum(version_ripend160hash)//两次hash256

	bytes := append(version_ripend160hash, checksum...)
	return Base58Encode(bytes)

}

func Ripemd160Hash(pubKey []byte) []byte {
	hash256 := sha256.Sum256(pubKey)
	Ripemd160Hasher := ripemd160.New()
	_, err := Ripemd160Hasher.Write(hash256[:])
	if err != nil {
		log.Panic(err)
	}
	pubRipemd160 := Ripemd160Hasher.Sum(nil)
	return pubRipemd160
}

func Checksum(hash []byte) []byte {
	hash1 := sha256.Sum256(hash)
	hash2 := sha256.Sum256(hash1[:])
	return hash2[:addressChechsumLen]

}

func ValidateAddress(address []byte) bool {
	version_pub_checksumBytes := Base58Decode(address)
	fmt.Println(version_pub_checksumBytes)

	checksum := version_pub_checksumBytes[len(version_pub_checksumBytes) - addressChechsumLen:]
	version_ripemd160 := version_pub_checksumBytes[:len(version_pub_checksumBytes) - addressChechsumLen]

	checksumBytes := Checksum(version_ripemd160)

	if bytes.Compare(checksum, checksumBytes) == 0 {
		return true
	}
	return false
}