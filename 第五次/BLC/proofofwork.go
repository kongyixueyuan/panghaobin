package BLC

import (
	"math"
	"math/big"
	"bytes"
	"crypto/sha256"
	"fmt"
)

const targetBits = 20
const maxNonce = math.MaxInt64

type ProofOfWork struct {
	block *Block
	target *big.Int
}
// target == 1 左移 256 - targetBits 位
func NewProofWork(b *Block) *ProofOfWork  {
	target := big.NewInt(1)
	target.Lsh(target, uint(256-targetBits))
	pow := &ProofOfWork{b, target}
	return pow
}

//工作了证明需要的数据 PrevBlockHash Data Timestamp targetBits nonce
func (pow *ProofOfWork) prepareData(nonce int) []byte {
	data := bytes.Join(
		[][]byte{
			pow.block.PrevBlockHash,
			pow.block.HashTransactions(),
			IntToHex(pow.block.Timestamp),
			IntToHex(int64(targetBits)),
			IntToHex(int64(nonce)),
			IntToHex(int64(pow.block.Height)),
		},
		[]byte{},
	)
	return data
}

//寻找有效哈希
func (pow *ProofOfWork) Run() (int, []byte)  {
	var hashInt big.Int
	var hash [32]byte
	nonce := 0
	fmt.Printf(" block data: \"%s\"\n", pow.block.Txs)
	fmt.Printf(" block height: \"%d\"\n", pow.block.Height)

	for nonce < maxNonce {
		data := pow.prepareData(nonce)
		hash = sha256.Sum256(data)
		hashInt.SetBytes(hash[:])
		if hashInt.Cmp(pow.target) == -1 {
			fmt.Printf("\r %x", hash)
			break
		} else {
			fmt.Printf("\r %x", hash)
			nonce++
		}

	}
	fmt.Printf("\n block nonce: \"%d\"", nonce)
	fmt.Print("\n\n")
	return nonce, hash[:]
}

//验证工作量
func (pow *ProofOfWork) Validate() bool  {
	var hashInt big.Int
	data := pow.prepareData(pow.block.Nonce)
	hash := sha256.Sum256(data)
	hashInt.SetBytes(hash[:])

	isValid := hashInt.Cmp(pow.target) == -1

	return isValid
}

