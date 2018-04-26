package chain

import (
	"github.com/michain/dotcoin/wallet"
	"bytes"
	"encoding/gob"
	"log"
)

// TXOutput represents a transaction output
type TXOutput struct {
	Value      int
	PubKeyHash []byte
}



// Lock set PublicKeyHash to signs the output
// input must check this value to use
func (out *TXOutput) Lock(address []byte) {
	out.PubKeyHash = wallet.GetPubKeyHashFromAddress(address)
}

// IsLockedWithKey checks if the output can be used by the owner of the pubkey
func (out *TXOutput) IsLockedWithKey(pubKeyHash []byte) bool {
	return bytes.Compare(out.PubKeyHash, pubKeyHash) == 0
}

// NewTXOutput create a new TXOutput
func NewTXOutput(value int, address string) *TXOutput {
	out := &TXOutput{value, nil}
	out.Lock([]byte(address))
	return out
}


// Serialize serializes []TXOutput
func SerializeOutputs(outs []TXOutput) []byte {
	var buff bytes.Buffer

	enc := gob.NewEncoder(&buff)
	err := enc.Encode(outs)
	if err != nil {
		log.Panic(err)
	}

	return buff.Bytes()
}

// DeserializeOutputs deserializes TXOutputs
func DeserializeOutputs(data []byte) []TXOutput {
	var outputs []TXOutput

	dec := gob.NewDecoder(bytes.NewReader(data))
	err := dec.Decode(&outputs)
	if err != nil {
		log.Panic(err)
	}

	return outputs
}


