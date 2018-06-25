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
	"github.com/michain/dotcoin/util/hashx"
)

const genesisReward = 100
const genesisBlockHeight = 1


// Block represents a block in the blockchain
type Block struct {
	Timestamp     time.Time
	PrevBlockHash []byte
	MerkleRoot	  []byte
	Hash          []byte
	Difficult	  uint32
	Nonce         int64
	Height        int32
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
func NewBlock(transactions []*Transaction, prevBlockHash []byte, height int32, quit chan struct{}) (*Block, bool) {
	block := &Block{}
	block.Timestamp = time.Now()
	block.Transactions = transactions
	block.PrevBlockHash = prevBlockHash
	block.Difficult = getCorrectDifficult()
	block.Height = height

	merkleRoot := block.HashTransactions()
	block.MerkleRoot = merkleRoot

	isSolve := false
	if height > genesisBlockHeight{
		pow := proof.NewProofOfWork()
		isSolve = pow.SolveHash(prevBlockHash, merkleRoot, quit)
		block.Nonce = pow.Nonce
		block.Hash = pow.Hash[:]
	}
	return block, isSolve
}

// NewGenesisBlock creates and returns genesis Block
func NewGenesisBlock(address string) *Block {
	coinbaseTX := NewCoinbaseTX(address, genesisCoinbaseData, genesisReward)
	b, _:= NewBlock([]*Transaction{coinbaseTX}, []byte{}, genesisBlockHeight, nil)
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

// SetHeight sets the height of the block
func (b *Block) SetHeight(height int32) {
	b.Height = height
}

func (b *Block) GetHash() (*hashx.Hash) {
	hash,_ := hashx.NewHash(b.Hash)
	return hash
}

func (b *Block) GetPrevHash() (*hashx.Hash) {
	hash,_ := hashx.NewHash(b.PrevBlockHash)
	return hash
}



func getCorrectDifficult() uint32{
	return blockDefaultDifficult
}