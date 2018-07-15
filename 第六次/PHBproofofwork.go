package main

import (
	"math"
	"math/big"
	"bytes"
	"fmt"
	"crypto/sha256"
)

var (
	maxNonce = math.MaxInt64
)

const targetBits = 16

// ProofOfWork represents a proof-of-work
type PHBProofOfWork struct {
	phbblock  *PHBBlock
	phbtarget *big.Int
}

// NewProofOfWork builds and returns a ProofOfWork
func PHBNewProofOfWork(b *PHBBlock) *PHBProofOfWork {
	target := big.NewInt(1)
	target.Lsh(target, uint(256-targetBits))

	pow := &PHBProofOfWork{b, target}

	return pow
}

func (pow *PHBProofOfWork) PHBprepareData(nonce int) []byte {
	data := bytes.Join(
		[][]byte{
			pow.phbblock.PHBPrevBlockHash,
			pow.phbblock.PHBHashTransactions(),
			PHBIntToHex(pow.phbblock.PHBTimestamp),
			PHBIntToHex(int64(targetBits)),
			PHBIntToHex(int64(nonce)),
		},
		[]byte{},
	)

	return data
}

// Run performs a proof-of-work
func (pow *PHBProofOfWork) PHBRun() (int, []byte) {
	var hashInt big.Int
	var hash [32]byte
	nonce := 0

	fmt.Printf("Mining a new block")
	for nonce < maxNonce {
		data := pow.PHBprepareData(nonce)

		hash = sha256.Sum256(data)
		fmt.Printf("\r%x", hash)
		hashInt.SetBytes(hash[:])

		if hashInt.Cmp(pow.phbtarget) == -1 {
			break
		} else {
			nonce++
		}
	}
	fmt.Print("\n\n")

	return nonce, hash[:]
}

// Validate validates block's PoW
func (pow *PHBProofOfWork) PHBValidate() bool {
	var hashInt big.Int
	data := pow.PHBprepareData(pow.phbblock.PHBNonce)
	hash := sha256.Sum256(data)
	hashInt.SetBytes(hash[:])
	isValid := hashInt.Cmp(pow.phbtarget) == -1

	return isValid
}