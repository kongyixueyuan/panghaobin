package main

import (
	"crypto/ecdsa"
	"encoding/hex"
	"log"
	"fmt"
	"crypto/rand"
	"strings"
	"bytes"
	"encoding/gob"
	"crypto/sha256"
	"crypto/elliptic"
	"math/big"
)
const subsidy = 10
type PHBTransaction struct {
	PHBID []byte
	PHBVin []PHBTXInput
	PHBVout []PHBTXOutput
}

//是否是coinbase交易 coinbase仅有一个TXI，该TXI的Txid为空，Vout设置为-1
func (tx *PHBTransaction)PHBIsCoinbase() bool  {
	return len(tx.PHBVin)==1 && len(tx.PHBVin[0].PHBTxid) == 0 && tx.PHBVin[0].PHBVout == -1
}

func (tx *PHBTransaction) PHBSerialize() []byte {
	var encoded bytes.Buffer
	enc := gob.NewEncoder(&encoded)
	err := enc.Encode(tx)
	if err != nil {
		log.Panic(err)
	}
	return encoded.Bytes()
}

func (tx *PHBTransaction) PHBHash() []byte {
	var hash [32]byte
	txCopy := *tx
	txCopy.PHBID = []byte{}
	hash = sha256.Sum256(txCopy.PHBSerialize())
	return hash[:]
}

func (tx *PHBTransaction) PHBSign(privateKey ecdsa.PrivateKey, prevTxs map[string]PHBTransaction) {
	if tx.PHBIsCoinbase() {
		return
	}

	for _, vin := range tx.PHBVin {
		if prevTxs[hex.EncodeToString(vin.PHBTxid)].PHBID == nil {
			log.Panic("Error")
		}
	}
	txCopy := tx.PHBTrimmedCopy()

	for inID, vin := range txCopy.PHBVin {
		prevTx := prevTxs[hex.EncodeToString(vin.PHBTxid)]
		txCopy.PHBVin[inID].PHBSignature = nil
		txCopy.PHBVin[inID].PHBPubKey = prevTx.PHBVout[vin.PHBVout].PHBPubKeyHash
		dataToSign := fmt.Sprintf("%x\n", txCopy)

		r, s, err := ecdsa.Sign(rand.Reader, &privateKey, []byte(dataToSign))
		if err != nil {
			log.Panic(err)
		}
		signature := append(r.Bytes(), s.Bytes()...)

		tx.PHBVin[inID].PHBSignature = signature
		txCopy.PHBVin[inID].PHBPubKey = nil
	}
}

func (tx *PHBTransaction) PHBString() string {
	var lines []string
	lines = append(lines, fmt.Sprintf("--- Transaction %x:", tx.PHBID))
	for i, input := range tx.PHBVin {
		lines = append(lines, fmt.Sprintf("     Input %d:", i))
		lines = append(lines, fmt.Sprintf("       TXID:      %x", input.PHBTxid))
		lines = append(lines, fmt.Sprintf("       Out:       %d", input.PHBVout))
		lines = append(lines, fmt.Sprintf("       Signature: %x", input.PHBSignature))
		lines = append(lines, fmt.Sprintf("       PubKey:    %x", input.PHBPubKey))
	}
	for i, output := range tx.PHBVout {
		lines = append(lines, fmt.Sprintf("     Output %d:", i))
		lines = append(lines, fmt.Sprintf("       Value:  %d", output.PHBValue))
		lines = append(lines, fmt.Sprintf("       Script: %x", output.PHBPubKeyHash))
	}
	return strings.Join(lines, "\n")
}

func (tx *PHBTransaction) PHBTrimmedCopy() PHBTransaction {
	var inputs []PHBTXInput
	var outputs []PHBTXOutput
	for _, in := range tx.PHBVin {
		inputs = append(inputs, PHBTXInput{in.PHBTxid, in.PHBVout, nil,nil})
	}

	for _, out := range tx.PHBVout {
		outputs = append(outputs, PHBTXOutput{out.PHBValue, out.PHBPubKeyHash})
	}
	txCopy := PHBTransaction{tx.PHBID, inputs, outputs}
	return txCopy
}

func (tx *PHBTransaction) PHBVerify(prevTXs map[string]PHBTransaction) bool {
	if tx.PHBIsCoinbase() {
		return true
	}

	for _, vin := range tx.PHBVin {
		if prevTXs[hex.EncodeToString(vin.PHBTxid)].PHBID == nil {
			log.Panic("ERROR: Previous transaction is not correct")
		}
	}

	txCopy := tx.PHBTrimmedCopy()
	curve := elliptic.P256()

	for inID, vin := range tx.PHBVin {
		prevTx := prevTXs[hex.EncodeToString(vin.PHBTxid)]
		txCopy.PHBVin[inID].PHBSignature = nil
		txCopy.PHBVin[inID].PHBPubKey = prevTx.PHBVout[vin.PHBVout].PHBPubKeyHash

		r := big.Int{}
		s := big.Int{}
		sigLen := len(vin.PHBSignature)
		r.SetBytes(vin.PHBSignature[:(sigLen / 2)])
		s.SetBytes(vin.PHBSignature[(sigLen / 2):])

		x := big.Int{}
		y := big.Int{}
		keyLen := len(vin.PHBPubKey)
		x.SetBytes(vin.PHBPubKey[:(keyLen / 2)])
		y.SetBytes(vin.PHBPubKey[(keyLen / 2):])

		dataToVerify := fmt.Sprintf("%x\n", txCopy)

		rawPubKey := ecdsa.PublicKey{Curve: curve, X: &x, Y: &y}
		if ecdsa.Verify(&rawPubKey, []byte(dataToVerify), &r, &s) == false {
			return false
		}
		txCopy.PHBVin[inID].PHBPubKey = nil
	}

	return true
}

// NewCoinbaseTX creates a new coinbase transaction
func PHBNewCoinbaseTX(to, data string) *PHBTransaction {
	if data == "" {
		randData := make([]byte, 20)
		_, err := rand.Read(randData)
		if err != nil {
			log.Panic(err)
		}

		data = fmt.Sprintf("%x", randData)
	}

	txin := PHBTXInput{[]byte{}, -1, nil, []byte(data)}
	txout := PHBNewTXOutPut(subsidy, to)
	tx := PHBTransaction{nil, []PHBTXInput{txin}, []PHBTXOutput{*txout}}
	tx.PHBID = tx.PHBHash()

	return &tx
}

// NewUTXOTransaction creates a new transaction
func  PHBNewUTXOTransaction(wallet *PHBWallet, to string, amount int, UTXOSet *PHBUTXOSet) *PHBTransaction {
	var inputs []PHBTXInput
	var outputs []PHBTXOutput

	pubKeyHash := PHBHashPubKey(wallet.PHBPublicKey)
	acc, validOutputs := UTXOSet.FindSpendableOutputs(pubKeyHash, amount)

	if acc < amount {
		log.Panic("ERROR: Not enough funds")
	}

	// Build a list of inputs
	for txid, outs := range validOutputs {
		txID, err := hex.DecodeString(txid)
		if err != nil {
			log.Panic(err)
		}

		for _, out := range outs {
			input := PHBTXInput{txID, out, nil, wallet.PHBPublicKey}
			inputs = append(inputs, input)
		}
	}

	// Build a list of outputs
	from := fmt.Sprintf("%s", wallet.PHBGetAddress())
	outputs = append(outputs, *PHBNewTXOutPut(amount, to))
	if acc > amount {
		outputs = append(outputs, *PHBNewTXOutPut(acc-amount, from)) // a change
	}

	tx := PHBTransaction{nil, inputs, outputs}
	tx.PHBID = tx.PHBHash()
	UTXOSet.PHBBlockchain.PHBSignTransaction(&tx, wallet.PHBPrivateKey)

	return &tx
}

// DeserializeTransaction deserializes a transaction
func PHBDeserializeTransaction(data []byte) PHBTransaction {
	var transaction PHBTransaction

	decoder := gob.NewDecoder(bytes.NewReader(data))
	err := decoder.Decode(&transaction)
	if err != nil {
		log.Panic(err)
	}

	return transaction
}