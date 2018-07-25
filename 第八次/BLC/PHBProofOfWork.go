package BLC

import (
	"bytes"
	"crypto/sha256"
	"fmt"
	"math/big"
)

const targetBit = 20

type PHBProofOfWork struct {
	phbBlock  *PHBBlock
	phbtarget *big.Int
}

func (pow *PHBProofOfWork) phbprepareData(nonce int, txHash []byte) []byte {
	data := bytes.Join(
		[][]byte{
			pow.phbBlock.PHBPrevBlockHash,
			txHash,
			PHBIntToHex(pow.phbBlock.PHBTimestamp),
			PHBIntToHex(int64(targetBit)),
			PHBIntToHex(int64(nonce)),
			PHBIntToHex(int64(pow.phbBlock.PHBHeight)),
		},
		[]byte{},
	)

	return data
}

func (proofOfWork *PHBProofOfWork) PHBRun() ([]byte, int64) {
	nonce := 0
	var hashInt big.Int // 存储我们新生成的hash
	var hash [32]byte
	txHash := proofOfWork.phbBlock.PHBHashTransactions()
	for {
		dataBytes := proofOfWork.phbprepareData(nonce, txHash)
		// 生成hash
		hash = sha256.Sum256(dataBytes)
		fmt.Printf("\r%x", hash)
		hashInt.SetBytes(hash[:])
		if proofOfWork.phbtarget.Cmp(&hashInt) == 1 {
			break
		}
		nonce = nonce + 1
	}
	return hash[:], int64(nonce)
}

func NewProofOfWork(block *PHBBlock) *PHBProofOfWork {
	target := big.NewInt(1)
	target = target.Lsh(target, 256-targetBit)
	return &PHBProofOfWork{block, target}
}
