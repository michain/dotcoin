package proof

import (
	"crypto/sha256"
	"math"
	"fmt"
	"math/big"
	"bytes"
	"github.com/michain/dotcoin/util"
)

var (
	maxNonce = math.MaxInt64
)

const targetBits = 16

// ProofOfWork represents a proof-of-work
type ProofOfWork struct {
	target *big.Int
}

// NewProofOfWork builds and returns a ProofOfWork
func NewProofOfWork() *ProofOfWork {
	target := big.NewInt(1)
	target.Lsh(target, uint(256-targetBits))

	pow := &ProofOfWork{target}

	return pow
}

func (pow *ProofOfWork) calculateHash(prevBlockHash, TXsHash []byte, nonce int) [32]byte {
	data := bytes.Join(
		[][]byte{
			prevBlockHash,
			TXsHash,
			util.IntToHex(int64(targetBits)),
			util.IntToHex(int64(nonce)),
		},
		[]byte{},
	)
	return sha256.Sum256(data)
}

func (pow *ProofOfWork) RunAtOnce(prevBlockHash, TXsHash []byte) (int, []byte){
	var hashInt big.Int
	var hash [32]byte
	nonce := 0
	//immediately return for test
	hash = pow.calculateHash(prevBlockHash, TXsHash, nonce)
	hashInt.SetBytes(hash[:])
	return nonce, hash[:]
}

// Run performs a proof-of-work
func (pow *ProofOfWork) Run(prevBlockHash, TXsHash []byte) (int, []byte) {
	var hashInt big.Int
	var hash [32]byte
	nonce := 0

	fmt.Printf("Mining a new block")
	for nonce < maxNonce {
		hash = pow.calculateHash(prevBlockHash, TXsHash, nonce)

		if math.Remainder(float64(nonce), 10000) == 0 {
			fmt.Printf("\r%x", hash)
		}
		hashInt.SetBytes(hash[:])

		if hashInt.Cmp(pow.target) == -1 {
			break
		} else {
			nonce++
		}
	}
	fmt.Print("\n\n")

	return nonce, hash[:]
}

// Validate validates block's PoW
func (pow *ProofOfWork) Validate(prevBlockHash, TXsHash []byte, nonce int) bool {
	var hashInt big.Int

	hash := pow.calculateHash(prevBlockHash, TXsHash, nonce)
	hashInt.SetBytes(hash[:])

	isValid := hashInt.Cmp(pow.target) == -1

	return isValid
}
