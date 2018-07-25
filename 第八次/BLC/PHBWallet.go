package BLC

import (
	"bytes"
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/sha256"
	"fmt"
	"golang.org/x/crypto/ripemd160"
	"log"
)

const version = byte(0x00)
const addressChecksumLen = 4

type PHBWallet struct {
	PHBPrivateKey ecdsa.PrivateKey //私钥
	PHBPublicKey  []byte           //公钥
}

func PHBIsValidForAdress(adress []byte) bool {
	version_public_checksumBytes := PHBBase58Decode(adress)
	fmt.Println(version_public_checksumBytes)
	checkSumBytes := version_public_checksumBytes[len(version_public_checksumBytes)-addressChecksumLen:]
	version_ripemd160 := version_public_checksumBytes[:len(version_public_checksumBytes)-addressChecksumLen]
	checkBytes := PHBCheckSum(version_ripemd160)
	if bytes.Compare(checkSumBytes, checkBytes) == 0 {
		return true
	}
	return false
}

func (w *PHBWallet) PHBGetAddress() []byte {
	ripemd160Hash := PHBRipemd160Hash(w.PHBPublicKey)
	version_ripemd160Hash := append([]byte{version}, ripemd160Hash...)
	checkSumBytes := PHBCheckSum(version_ripemd160Hash)
	bytes := append(version_ripemd160Hash, checkSumBytes...)
	return PHBBase58Encode(bytes)
}

func PHBCheckSum(payload []byte) []byte {
	hash1 := sha256.Sum256(payload)
	hash2 := sha256.Sum256(hash1[:])
	return hash2[:addressChecksumLen]
}

func PHBRipemd160Hash(publicKey []byte) []byte {
	hash256 := sha256.New()
	hash256.Write(publicKey)
	hash := hash256.Sum(nil)
	ripemd160 := ripemd160.New()
	ripemd160.Write(hash)
	return ripemd160.Sum(nil)
}

//创建钱包
func PHBNewWallet() *PHBWallet {
	privateKey, publicKey := phbnewKeyPair()
	return &PHBWallet{privateKey, publicKey}
}

func phbnewKeyPair() (ecdsa.PrivateKey, []byte) {
	curve := elliptic.P256()
	private, err := ecdsa.GenerateKey(curve, rand.Reader)
	if err != nil {
		log.Panic(err)
	}
	pubKey := append(private.PublicKey.X.Bytes(), private.PublicKey.Y.Bytes()...)
	return *private, pubKey
}
