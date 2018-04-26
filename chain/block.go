package chain

import (
	"time"
	"bytes"
	"encoding/gob"
	"log"
	"github.com/michain/dotcoin/merkle"
	"github.com/michain/dotcoin/proof"
	"fmt"
	"strings"
)

const genesisReward = 100

// Block represents a block in the blockchain
type Block struct {
	Timestamp     int64
	PrevBlockHash []byte
	MerkleRoot	  []byte
	Hash          []byte
	Difficult	  int
	Nonce         int
	Height        int
	Transactions  []*Transaction
}

func (b *Block) String() string{
	var lines []string

	lines = append(lines, fmt.Sprintf("--- Block %x:", b.Hash))

	lines = append(lines, fmt.Sprintf("    PrevBlockHash:   %x", b.PrevBlockHash))
	lines = append(lines, fmt.Sprintf("    Hash:            %x", b.Hash))
	lines = append(lines, fmt.Sprintf("    MerkleRoot:      %x", b.MerkleRoot))
	lines = append(lines, fmt.Sprintf("    Timestamp:       %d", b.Timestamp))
	lines = append(lines, fmt.Sprintf("    Difficult:       %d", b.Difficult))
	lines = append(lines, fmt.Sprintf("    Nonce:           %d", b.Nonce))
	lines = append(lines, fmt.Sprintf("    Height:          %d", b.Height))

	lines = append(lines, fmt.Sprintf("    Transactions     %d:", len(b.Transactions)))
	for _, tx := range b.Transactions{
		lines = append(lines, fmt.Sprintf(tx.String()))
	}


	return strings.Join(lines, "\n")
}


// SerializeBlock serializes the block
func SerializeBlock(b *Block) []byte{
	var result bytes.Buffer
	encoder := gob.NewEncoder(&result)

	err := encoder.Encode(b)
	if err != nil {
		log.Panic(err)
	}

	return result.Bytes()
}

// DeserializeBlock deserializes a block
func DeserializeBlock(d []byte) *Block {
	var block Block

	decoder := gob.NewDecoder(bytes.NewReader(d))
	err := decoder.Decode(&block)
	if err != nil {
		log.Panic(err)
	}

	return &block
}

// NewBlock creates and returns Block
func NewBlock(transactions []*Transaction, prevBlockHash []byte, height int) *Block {
	block := &Block{}
	block.Timestamp = time.Now().Unix()
	block.Transactions = transactions
	block.PrevBlockHash = prevBlockHash
	block.Difficult = getCorrectDifficult()
	block.Height = height

	merkleRoot := block.HashTransactions()
	block.MerkleRoot = merkleRoot

	pow := proof.NewProofOfWork()
	nonce, hash := pow.Run(prevBlockHash, merkleRoot)
	block.Nonce = nonce
	block.Hash = hash
	return block
}

// NewGenesisBlock creates and returns genesis Block
func NewGenesisBlock(address string) *Block {
	coinbaseTX := NewCoinbaseTX(address, genesisCoinbaseData, genesisReward)
	b:= NewBlock([]*Transaction{coinbaseTX}, []byte{}, 0)
	fmt.Println(b.String())
	return b
}

// HashTransactions returns a hash of the transactions in the block
func (b *Block) HashTransactions() []byte {
	var transactions [][]byte
	for _, tx := range b.Transactions {
		transactions = append(transactions, SerializeTransaction(tx))
	}
	mTree := merkle.NewMerkleTree(transactions)
	return mTree.RootNode.Data
}


func getCorrectDifficult() int{
	return blockDefaultDifficult
}