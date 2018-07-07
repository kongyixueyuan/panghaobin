package BLC

import (
	"time"
	"bytes"
	"encoding/gob"
	"log"
	"crypto/sha256"
	"fmt"
)

type Block struct {
	Timestamp int64 //时间戳

	Height int64 //区块高度
	PrevBlockHash []byte //上一个区块HASH
	Hash []byte //本区块的hash
	Txs []*Transaction //交易数据
	Nonce int //计算PoW
}

func NewBlock(txs []*Transaction, height int64, prevBlockHash []byte) *Block  {
	block := &Block{
		Timestamp: time.Now().Unix(),
		Height: height,
		PrevBlockHash: prevBlockHash,
		Hash: []byte{},
		Txs: txs,
		Nonce: 0}
	pow := NewProofWork(block)
	nonce, hash := pow.Run()
	block.Hash = hash[:]
	block.Nonce = nonce

	return block
}

func NewGenesisBlock(txs []*Transaction) *Block {
	return NewBlock(txs, 0, []byte{})
}

//区块序列化
func (b *Block) Serialize() []byte  {
	var result bytes.Buffer
	encoder := gob.NewEncoder(&result)
	err := encoder.Encode(b)
	if err != nil {
		log.Panic(err)
	}

	return result.Bytes()
}

//反序列化
func DeserializeBlock(d []byte) *Block  {
	var block Block

	decoder := gob.NewDecoder(bytes.NewBuffer(d))
	err := decoder.Decode(&block)
	if err != nil {
		log.Panic(err)
	}
	return &block
}

//返回Transaction的字节数组
func (b *Block) HashTransactions() []byte {
	var txHashes [][]byte
	var txHash [32]byte
	for _, tx := range b.Txs {
		txHashes = append(txHashes, tx.TxHash)
	}
	txHash = sha256.Sum256(bytes.Join(txHashes, []byte{}))
	return txHash[:]
}

func (block *Block) printfBlock()  {
	fmt.Printf("========================================\n")
	fmt.Printf("区块高度：%d\n", block.Height)
	fmt.Printf("上一个区块的hash：%x\n", block.PrevBlockHash)
	fmt.Printf("区块交易数据\n")
	for _, tx := range block.Txs {
		fmt.Printf("交易hash：%x\n", tx.TxHash)
		fmt.Println("  交易输入:")
		for _, in := range tx.Vin {
			fmt.Printf("%x\n", in.TxHash)
			fmt.Printf("%d\n", in.Vout)
			fmt.Printf("用户签名：%s\n", in.Signature)
		}

		fmt.Println("  交易输出:")
		for _, out := range tx.Vout {
			fmt.Println(out.Value)
			fmt.Println(out.PubKey)
		}
	}
	fmt.Println("------------------------------\n")
	fmt.Printf("区块时间戳：%s\n", time.Unix(block.Timestamp, 0).Format("2006-01-02 15:04:05 AM") )//go诞生之日
	fmt.Printf("区块hash：%x\n", block.Hash)
	fmt.Printf("Nonce: %d\n", block.Nonce)
	fmt.Printf("========================================\n\n")
}
