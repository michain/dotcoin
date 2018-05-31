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
	Nonce int64
	Hash [32]byte
}

// NewProofOfWork builds and returns a ProofOfWork
func NewProofOfWorkT(targetBits int) *ProofOfWork {
	target := big.NewInt(1)
	target.Lsh(target, uint(targetBits))

	pow := &ProofOfWork{target:target}

	return pow
}

// NewProofOfWork builds and returns a ProofOfWork
func NewProofOfWork() *ProofOfWork {
	target := big.NewInt(1)
	target.Lsh(target, uint(256-targetBits))

	pow := &ProofOfWork{target:target}

	return pow
}

// calculateHash calc hash with bestBlockHash and Txs hashes
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

// solveHash solve right hash which less than the target difficulty
// it will be stop when received quit signal
func (pow *ProofOfWork) solveHash(prevBlockHash, TXsHash []byte, quit chan struct{}) bool{
	var hashInt big.Int
	var hash [32]byte
	nonce := 0
	for nonce < maxNonce {
		select {
		case <-quit:
			return false
		default:
			hash = pow.calculateHash(prevBlockHash, TXsHash, nonce)

			if math.Remainder(float64(nonce), 10000) == 0 {
				fmt.Printf("\r%x", hash)
			}

			hashInt.SetBytes(hash[:])
			if hashInt.Cmp(pow.target) == -1 {
				pow.Nonce = int64(nonce)
				pow.Hash = hash
				return true
			} else {
				nonce++
			}
		}
	}
	return false
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

// SolveHash loop calc hash to solve target
func (pow *ProofOfWork) SolveHash(prevBlockHash, TXsHash []byte, quit chan struct{}) bool {
	fmt.Println("Mining a new block begin")
	isSolve := pow.solveHash(prevBlockHash, TXsHash, quit)
	fmt.Println("Mining a new block end, IsSolve:", isSolve)
	return isSolve
}

// Validate validates block's PoW
func (pow *ProofOfWork) Validate(prevBlockHash, TXsHash []byte, nonce int) bool {
	var hashInt big.Int

	hash := pow.calculateHash(prevBlockHash, TXsHash, nonce)
	hashInt.SetBytes(hash[:])

	isValid := hashInt.Cmp(pow.target) == -1

	return isValid
}
