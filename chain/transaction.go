package chain

import (
	"bytes"
	"encoding/gob"
	"log"
	"crypto/ecdsa"
	"fmt"
	"crypto/rand"
	"crypto/elliptic"
	"math/big"
	"github.com/michain/dotcoin/wallet"
	"crypto/sha256"
	"strings"
	"github.com/michain/dotcoin/util"
	"github.com/michain/dotcoin/util/hashx"
)



// Transaction a transaction with ID\input\output
type Transaction struct{
	ID      hashx.Hash
	Inputs  []TXInput
	Outputs []TXOutput
}

func (tx Transaction) StringID() string{
	return tx.ID.String()
}


// IsCoinBase checks whether the transaction is coinbase
func (tx Transaction) IsCoinBase() bool {
	return len(tx.Inputs) == 1 && tx.Inputs[0].PreviousOutPoint.Hash.IsEqual(hashx.ZeroHash()) && tx.Inputs[0].PreviousOutPoint.Index == -1
}


// String returns a human-readable representation of a transaction
func (tx Transaction) String() string {
	var lines []string

	lines = append(lines, fmt.Sprintf("--- Transaction %v:", tx.ID))

	for i, input := range tx.Inputs {

		lines = append(lines, fmt.Sprintf("     Input %d:", i))
		lines = append(lines, fmt.Sprintf("       TXID:         %v", input.PreviousOutPoint.Hash.String()))
		lines = append(lines, fmt.Sprintf("       OutIndex:     %d", input.PreviousOutPoint.Index))
		lines = append(lines, fmt.Sprintf("       Signature:    %x", input.Signature))
		lines = append(lines, fmt.Sprintf("       PubKey:       %x", input.PubKey))
	}

	for i, output := range tx.Outputs {
		lines = append(lines, fmt.Sprintf("     Output %d:", i))
		lines = append(lines, fmt.Sprintf("       Value:        %d", output.Value))
		lines = append(lines, fmt.Sprintf("       PubKeyHash:   %x", output.PubKeyHash))
	}

	return strings.Join(lines, "\n")
}

// Sign signs each input of a Transaction
// must match input's prev TX exists
func (tx *Transaction) Sign(privateKey ecdsa.PrivateKey, prevTXs map[string]Transaction) {
	if tx.IsCoinBase() {
		return
	}

	//check input's prev TX exists
	for _, vin := range tx.Inputs {
		if _, exists:= prevTXs[vin.PreviousOutPoint.StringHash()];!exists{
			log.Panic("ERROR: Previous transaction is not correct")
		}
	}

	//get TX's trimmed copy
	txCopy := tx.TrimmedCopy()

	for inID, input := range txCopy.Inputs {
		prevTx := prevTXs[input.PreviousOutPoint.StringHash()]
		txCopy.Inputs[inID].Signature = nil
		txCopy.Inputs[inID].PubKey = prevTx.Outputs[input.PreviousOutPoint.Index].PubKeyHash //why no use input's raw public key?

		dataToSign := fmt.Sprintf("%x\n", txCopy)

		r, s, err := ecdsa.Sign(rand.Reader, &privateKey, []byte(dataToSign))
		if err != nil {
			log.Panic(err)
		}
		signature := append(r.Bytes(), s.Bytes()...)

		tx.Inputs[inID].Signature = signature
		txCopy.Inputs[inID].PubKey = nil
	}
}

// Verify verifies signatures of Transaction inputs
// use signature & rawPubKey on ecdsa.Verify
func (tx *Transaction) Verify(prevTXs map[string]Transaction) bool {
	if tx.IsCoinBase() {
		return true
	}

	for _, vin := range tx.Inputs {
		if _, exists:= prevTXs[vin.PreviousOutPoint.StringHash()];!exists{
			log.Panic("ERROR: Previous transaction is not correct")
		}
	}

	txCopy := tx.TrimmedCopy()
	curve := elliptic.P256()

	for inID, vin := range tx.Inputs {
		prevTx := prevTXs[vin.PreviousOutPoint.StringHash()]
		txCopy.Inputs[inID].Signature = nil
		txCopy.Inputs[inID].PubKey = prevTx.Outputs[vin.PreviousOutPoint.Index].PubKeyHash

		r := big.Int{}
		s := big.Int{}
		sigLen := len(vin.Signature)
		r.SetBytes(vin.Signature[:(sigLen / 2)])
		s.SetBytes(vin.Signature[(sigLen / 2):])

		x := big.Int{}
		y := big.Int{}
		keyLen := len(vin.PubKey)
		x.SetBytes(vin.PubKey[:(keyLen / 2)])
		y.SetBytes(vin.PubKey[(keyLen / 2):])

		dataToVerify := fmt.Sprintf("%x\n", txCopy)

		rawPubKey := ecdsa.PublicKey{Curve: curve, X: &x, Y: &y}
		if ecdsa.Verify(&rawPubKey, []byte(dataToVerify), &r, &s) == false {
			return false
		}
		txCopy.Inputs[inID].PubKey = nil
	}

	return true
}


// TrimmedCopy creates a trimmed copy of Transaction to be used in signing
// set sign & pubkey nil
func (tx *Transaction) TrimmedCopy() Transaction {
	var inputs []TXInput
	var outputs []TXOutput

	for _, vin := range tx.Inputs {
		inputs = append(inputs, *NewTXInput(&vin.PreviousOutPoint, nil, nil))
	}
	for _, vout := range tx.Outputs {
		outputs = append(outputs, TXOutput{vout.Value, vout.PubKeyHash})
	}

	txCopy := Transaction{tx.ID, inputs, outputs}
	return txCopy
}

// Hash returns the hash of the Transaction
func (tx *Transaction) Hash() hashx.Hash {
	var hash [32]byte

	txCopy := *tx
	txCopy.ID = *hashx.ZeroHash()

	hash = sha256.Sum256(SerializeTransaction(&txCopy))

	return hash
}


func (tx *Transaction) StringHash() string{
	return tx.Hash().String()
}

// SerializeTransaction serializes a transaction for []byte
func SerializeTransaction(tx *Transaction) []byte{
	var encoded bytes.Buffer

	enc := gob.NewEncoder(&encoded)
	err := enc.Encode(tx)
	if err != nil {
		log.Panic(err)
	}
	return encoded.Bytes()
}

// DeserializeTransaction deserializes a transaction
func DeserializeTransaction(data []byte) (*Transaction, error) {
	var transaction Transaction

	decoder := gob.NewDecoder(bytes.NewReader(data))
	err := decoder.Decode(&transaction)
	return &transaction, err
}

// NewCoinbaseTX creates a new coinbase transaction
func NewCoinbaseTX(to, data string, reward int) *Transaction {
	if data == ""{
		data = util.GetRandData()
	}
	fmt.Println("NewCoinbaseTX RandData", data)
	txin := NewTXInput(NewOutPoint(hashx.ZeroHash(), -1), nil, []byte(data))
	txout := NewTXOutput(reward, to)
	fmt.Println(txout)
	tx := Transaction{*hashx.ZeroHash(), []TXInput{*txin}, []TXOutput{*txout}}

	tx.ID = tx.Hash()

	return &tx
}

// NewUTXOTransaction creates a new transaction
func NewUTXOTransaction(fromWallet *wallet.Wallet, to string, amount int, UTXOSet *UTXOSet, txPool TxPool) *Transaction {
	var inputs []TXInput
	var outputs []TXOutput

	pubKeyHash := wallet.HashPublicKey(fromWallet.PublicKey)
	acc, validOutputs := UTXOSet.FindSpendableOutputs(pubKeyHash, amount, txPool)

	if acc < amount {
		log.Panic("ERROR: Not enough funds")
	}

	// Build a list of inputs
	for txid, outs := range validOutputs {
		var txID hashx.Hash
		err := hashx.Decode(&txID, txid)
		fmt.Println("NewUTXOTransaction", txID, txid)
		if err != nil {
			log.Panic(err)
		}

		for _, out := range outs {
			input := TXInput{*NewOutPoint(&txID, out), nil, fromWallet.PublicKey}
			inputs = append(inputs, input)
		}
	}

	// Build a list of outputs
	from := fromWallet.GetStringAddress()
	outputs = append(outputs, *NewTXOutput(amount, to))
	if acc > amount {
		outputs = append(outputs, *NewTXOutput(acc-amount, from)) // a change
	}

	tx := Transaction{*hashx.ZeroHash(), inputs, outputs}
	tx.ID = tx.Hash()
	UTXOSet.Blockchain.SignTransaction(&tx, fromWallet.PrivateKey)

	return &tx
}

